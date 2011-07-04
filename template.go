// Copyright 2011, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
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

var (
    NoTemplateError = os.NewError("Requested template does not exist")
    ParseError      = os.NewError("Couldn't parse template")
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
    var (
        joined = filepath.Join(path...)
        stat, errStat = os.Stat(joined)
    )
    if stat == nil || errStat != nil {
        return ""
    }
    return joined
}
func GetAltTemplatePath(relpath []string) string {
    if AppConfig.AltRoot == "" {
        return ""
    }
    var (
        altpath = GetRootedTemplatePath([]string{AppConfig.AltRoot}, relpath)
        stat, errStat = os.Stat(altpath)
    )
    if stat == nil || errStat != nil {
        return ""
    }
    return altpath
}
//  Get a full template path from a path slice relative to another path
//  slice.
func GetRootedTemplatePath(rootpath []string, relpath []string) string {
    var path = make([]string, len(rootpath)+len(relpath))
    copy(path, rootpath)
    copy(path[len(rootpath):], relpath)
    return filepath.Join(path...)
}
func ParseAltTemplate(filename string, dict map[string]string, relpath []string) (string, os.Error) {
    var tpath = GetAltTemplatePath(relpath)
    if tpath == "" {
        return "", NoTemplateError
    }
    if DEBUG && DEBUG_LEVEL > 0 {
        log.Printf("scanning: %s", tpath)
        if DEBUG_LEVEL > 1 {
            log.Printf("context:\n%v", dict)
        }
    }
    // TODO Stat the template to make sure it exists.
    return mustache.RenderFile(tpath, dict, map[string]string{"file":filename}), nil
}
//  Given a filename and dictionary context, create a context dict+("file"=>filename),
//  and read a template specified by relpath. See GetTemplatePath().
func ParseTemplate(filename string, dict map[string]string, relpath []string) (string, os.Error) {
    var tpath = GetTemplatePath(relpath)
    if tpath == "" {
        return "", NoTemplateError
    }
    if DEBUG && DEBUG_LEVEL > 0 {
        log.Printf("scanning: %s", tpath)
        if DEBUG_LEVEL > 1 {
            log.Printf("context:\n%v", dict)
        }
    }
    // TODO Stat the template to make sure it exists.
    return mustache.RenderFile(tpath, dict, map[string]string{"file":filename}), nil
}
//  Given a filename, dictionary context, and the path to a template,
//  write the parsed template to the specified filename. The context of
//  the template will have a rule "file":filename which should override
//  any previous "file" rule in dict.
func WriteTemplate(filename, desc string, dict map[string]string, relpath...string) os.Error {
    var template string
    var alttemplate, errParseAlt = ParseAltTemplate(filename, dict, relpath)
    if errParseAlt == nil {
        template = alttemplate
    } else {
        var stdtemplate, errParseStd = ParseTemplate(filename, dict, relpath)
        if errParseStd != nil {
            return errParseStd
        }
        template = stdtemplate
    }
	if DEBUG || VERBOSE {
		fmt.Printf("Creating %s %s\n", desc, filename)
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
    var template string
    var alttemplate, errParseAlt = ParseAltTemplate(filename, dict, relpath)
    if errParseAlt == nil {
        template = alttemplate
    } else {
        var stdtemplate, errParseStd = ParseTemplate(filename, dict, relpath)
        if errParseStd != nil {
            return errParseStd
        }
        template = stdtemplate
    }
	if DEBUG || VERBOSE {
		fmt.Printf("Appending %s %s\n", desc, filename)
        if DEBUG && DEBUG_LEVEL > 2 {
            log.Print("\n", template, "\n")
        }
    }
    var fout, errOpen = os.OpenFile(filename, os.O_WRONLY|os.O_APPEND, FilePermissions)
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
