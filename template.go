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
    "bytes"
    "strings"
    "io/ioutil"
    "path/filepath"
    "template"
    //"github.com/hoisie/mustache.go"
)

var (
    NoTemplateError = os.NewError("Requested template does not exist")
    ParseError      = os.NewError("Couldn't parse template")
)

func TestName(filename string) string {
    var test = filename
    if strings.HasSuffix(test, ".go") {
        test = test[:len(test)-3]
    }
    if strings.HasSuffix(test, "_test") {
        test = test[:len(test)-5]
    }
    return strings.Title(test)
}

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
    return []string{GetGoroot(), "src", "pkg",
        "github.com", "bmatsuo", "gonew", "templates"}
}

//  Get a full template path from a path slice relative to the templates
//  directory.
func GetTemplatePath(relpath []string) string {
    var (
        rootpath = GetTemplateRoot()
        path     = make([]string, len(rootpath)+len(relpath))
    )
    copy(path, rootpath)
    copy(path[len(rootpath):], relpath)
    joined := filepath.Join(path...)
    if stat, err := os.Stat(joined); stat == nil || err != nil {
        return ""
    }
    return joined
}

func GetAltTemplatePath(relpath []string) string {
    if AppConfig.AltRoot == "" {
        if DEBUG {
            log.Print("No alt root found.")
        }
        return ""
    }
    altpath := GetRootedTemplatePath([]string{AppConfig.AltRoot}, relpath)
    if stat, err := os.Stat(altpath); stat == nil || err != nil {
        if DEBUG {
            log.Printf("Error stat'ing %s.", altpath)
        }
        return ""
    }
    return altpath
}

//  Get a full template path from a path slice relative to another path
//  slice.
func GetRootedTemplatePath(rootpath []string, relpath []string) string {
    path := make([]string, len(rootpath)+len(relpath))
    copy(path, rootpath)
    copy(path[len(rootpath):], relpath)
    return filepath.Join(path...)
}

func extraData(filename string) map[string]string {
    return map[string]string{"file": filename, "test": TestName(filename)}
}

func combined(dict, extra map[string]string) map[string]string {
    numEntries := len(dict) + len(extra)
    comb := make(map[string]string, numEntries)
    add := func(d map[string]string) {
        for k, v := range d {
            comb[k] = v
        }
    }
    add(dict)
    add(extra)
    return comb
}

func ParseAltTemplate(filename string, dict map[string]string, relpath []string) (string, os.Error) {
    var templ *template.Template
    if tpath := GetTemplatePath(relpath); tpath == "" {
        return "", NoTemplateError
    } else {
        templ = template.MustParseFile(tpath, nil)
        if DEBUG && DEBUG_LEVEL > 0 {
            log.Printf("scanning: %s", tpath)
            if DEBUG_LEVEL > 1 {
                log.Printf("context:\n%v", dict)
            }
        }
    }

    buff := bytes.NewBuffer(make([]byte, 0, 1<<20))
    errTExec := template.Execute(buff, combined(dict, extraData(filename)))
    return buff.String(), errTExec
}

//  Given a filename and dictionary context, create a context dict+("file"=>filename),
//  and read a template specified by relpath. See GetTemplatePath().
func ParseTemplate(filename string, dict map[string]string, relpath []string) (string, os.Error) {
    var templ *template.Template
    if tpath := GetTemplatePath(relpath); tpath == "" {
        return "", NoTemplateError
    } else {
        templ = template.MustParseFile(tpath, nil)

        if DEBUG && DEBUG_LEVEL > 0 {
            log.Printf("scanning: %s", tpath)
            if DEBUG_LEVEL > 1 {
                log.Printf("context:\n%v", dict)
            }
        }
    }

    buff := bytes.NewBuffer(make([]byte, 0, 1<<20))
    errTExec := template.Execute(buff, combined(dict, extraData(filename)))
    return buff.String(), errTExec
}

//  Given a filename, dictionary context, and the path to a template,
//  write the parsed template to the specified filename. The context of
//  the template will have a rule "file":filename which should override
//  any previous "file" rule in dict.
func WriteTemplate(filename, desc string, dict map[string]string, relpath ...string) os.Error {
    var templ string
    if altt, err := ParseAltTemplate(filename, dict, relpath); err == nil {
        templ = altt
        if DEBUG || VERBOSE {
            fmt.Printf("Using alternate template %s\n", GetAltTemplatePath(relpath))
        }
    } else if stdt, err := ParseTemplate(filename, dict, relpath); err == nil {
        templ = stdt
    } else {
        return err
    }

    if DEBUG || VERBOSE {
        fmt.Printf("Creating %s %s\n", desc, filename)
        if DEBUG && DEBUG_LEVEL > 2 {
            log.Print("\n", template, "\n")
        }
    }

    templout := make([]byte, len(template))
    copy(templout, template)
    return ioutil.WriteFile(filename, templout, FilePermissions)
}
func AppendTemplate(filename, desc string, dict map[string]string, relpath ...string) os.Error {
    var templ string
    if altt, err := ParseAltTemplate(filename, dict, relpath); err == nil {
        templ = altt
        if DEBUG || VERBOSE {
            fmt.Printf("Using alternate template %s\n", GetAltTemplatePath(relpath))
        }
    } else if stdt, err := ParseTemplate(filename, dict, relpath); err == nil {
        templ = stdt
    } else {
        return err
    }

    if DEBUG || VERBOSE {
        fmt.Printf("Appending %s %s\n", desc, filename)
        if DEBUG && DEBUG_LEVEL > 2 {
            log.Print("\n", templ, "\n")
        }
    }

    fout, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND, FilePermissions)
    if err != nil {
        return err
    }
    if _, err := fout.WriteString(templ); err != nil {
        return err
    }
    if err := fout.Close(); err != nil {
        return err
    }
    return nil
}

/* Some functions for tests and debugging. */
func getDebugTemplateRoot() []string { return []string{"templates"} }
func getDebugTemplatePath(relpath ...string) string {
    return GetRootedTemplatePath(getDebugTemplateRoot(), relpath)
}
