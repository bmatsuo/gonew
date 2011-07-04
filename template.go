package main
/*
 *  Filename:    template.go
 *  Package:     main
 *  Author:      Bryan Matsuo <bmatsuo@soe.ucsc.edu>
 *  Created:     Sun Jul  3 17:55:40 PDT 2011
 *  Description: 
 */
import (
    "os"
    "fmt"
    "log"
    "io/ioutil"
    "path/filepath"
    "github.com/hoisie/mustache.go"
)

//  The $GOROOT environment variable.
func GetGoroot() string {
    goroot, err := os.Getenverror("GOROOT")
    if err != nil {
        panic("goroot")
    }
    return goroot
}

//  The template directory of the goinstall'ed gonew package.
func GetTemplateRoot() []string {
    var goroot = GetGoroot()
    return []string{goroot, "src", "pkg",
            "github.com", "bmatsuo", "gonew", "templates"}
}

//  Get a full template path from a path slice relative to the templates
//  directory.
func GetTemplatePath(relpath []string) string {
    var (
        rootpath = GetTemplateRoot()
        path = make([]string, len(rootpath)+len(relpath))
    )
    copy(path, rootpath)
    copy(path[len(rootpath):], relpath)
    return filepath.Join(path...)
}
//  Get a full template path from a path slice relative to another path
//  slice.
func GetRootedTemplatePath(rootpath []string, relpath []string) string {
    var path = make([]string, len(rootpath)+len(relpath))
    copy(path, rootpath)
    copy(path[len(rootpath):], relpath)
    return filepath.Join(path...)
}
//  Given a filename and dictionary context, create a context dict+("file"=>filename),
//  and read a template specified by relpath. See GetTemplatePath().
func ParseTemplate(filename string, dict map[string]string, relpath []string) string {
    var tpath = GetTemplatePath(relpath)
    if DEBUG && DEBUG_LEVEL > 0 {
        log.Printf("scanning: %s", tpath)
        if DEBUG_LEVEL > 1 {
            log.Printf("context:\n%v", dict)
        }
    }
    return mustache.RenderFile(tpath, dict, map[string]string{"file":filename})
}
//  Given a filename, dictionary context, and the path to a template,
//  write the parsed template to the specified filename. The context of
//  the template will have a rule "file":filename which should override
//  any previous "file" rule in dict.
func WriteTemplate(filename, desc string, dict map[string]string, relpath...string) os.Error {
    var template = ParseTemplate(filename, dict, relpath)
	if DEBUG || VERBOSE {
		fmt.Printf("Creating %s %s", desc, filename)
        if DEBUG && DEBUG_LEVEL > 2 {
            log.Print("\n", template, "\n")
        }
    }
    var templout = make([]byte, len(template))
    copy(templout, template)
    var errWrite = ioutil.WriteFile(filename, templout, FilePermissions)
    return errWrite
}
func AppendTemplate(filename, desc string, dict map[string]string, relpath...string) os.Error {
    var template = ParseTemplate(filename, dict, relpath)
	if DEBUG || VERBOSE {
		fmt.Printf("Creating %s %s", desc, filename)
        if DEBUG && DEBUG_LEVEL > 2 {
            log.Print("\n", template, "\n")
        }
    }
    var fout, errOpen = os.Open(filename)
    if errOpen != nil {
        return errOpen
    }
    var _, errAppend = fout.WriteString(template)
    if errAppend != nil {
        return errAppend
    }
    var errClose = fout.Close()
    if errClose != nil {
        return errClose
    }
    return nil
}

/* Some functions for tests and debugging. */
func getDebugTemplateRoot() []string {
    return []string{"templates"}
}
func getDebugTemplatePath(relpath...string) string {
    return GetRootedTemplatePath(getDebugTemplateRoot(), relpath)
}
