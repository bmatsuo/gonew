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
    "io"
    "os"
    "fmt"
    //"log"
    "bytes"
    "strings"
    "io/ioutil"
    "path/filepath"
    "template"
    //"github.com/hoisie/mustache.go"
)

const (
    TemplateFileExt = ".t"
)

var (
    NoTemplateError = os.NewError("Requested template does not exist")
    ParseError      = os.NewError("Couldn't parse template")
)

type Executor interface {
    Execute(io.Writer, interface{}) os.Error
}

type ExecutorSet interface {
    Execute(io.Writer, string, interface{}) os.Error
}

func collectBytes(f func(w io.Writer) os.Error) (p []byte, err os.Error) {
    b := new(bytes.Buffer)
    err = f(b)
    p = b.Bytes()
    return
}

//  Call t.Execute with the given data, writing to a []byte buffer.
func Executed(t Executor, data interface{}) ([]byte, os.Error) {
    return collectBytes(func(w io.Writer)os.Error { return t.Execute(w, data) })
}

//  Call s.Execute with the given name data, writing to a []byte buffer.
func ExecutedSet(s ExecutorSet, name string, data interface{}) ([]byte, os.Error) {
    return collectBytes(func(w io.Writer)os.Error { return s.Execute(w, name, data) })
}

//  Returns an os-specific pattern <root>/"*"/*<TemplateFileExt>
func TemplateGlobPattern(root string) string {
    return filepath.Join(root, "*", fmt.Sprintf("*%s", TemplateFileExt))
}

type TemplateMultiSet []*template.Set

//  Call function CollectTemplates on each given root and create a TemplateMultiSet.
func MakeTemplateMultiSet(f template.FuncMap, roots... string) (ms TemplateMultiSet, err os.Error) {
    var s *template.Set
    ms = make(TemplateMultiSet, len(roots))
    for i := range roots {
        if s, err = CollectTemplates(roots[i], f); err != nil {
            break
        } else {
            ms[i] = s
        }
    }
    return
}

//  Execute a template in first s in ms for which s.Template(name) is non-nil.
func (ms TemplateMultiSet) Execute(wr io.Writer, name string, data interface{}) os.Error {
    for _, s := range ms {
        if t := s.Template(name); t != nil {
            return t.Execute(wr, data)
        }
    }
    return NoTemplateError
}

//  Parse templates <root>/*/*.t, allowing them a given function map.
func CollectTemplates(root string, f template.FuncMap) (s *template.Set, err os.Error) {
    s = new(template.Set)
    if f != nil {
        s.Funcs(f)
    }
    s, err = s.ParseTemplateGlob(TemplateGlobPattern(root))
    if err != nil {
        return
    }
    return
}

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
        Debug(0, "No alt root found.")
        return ""
    }
    altpath := GetRootedTemplatePath([]string{AppConfig.AltRoot}, relpath)
    if stat, err := os.Stat(altpath); stat == nil || err != nil {
        Debug(0, fmt.Sprintf("Error stat'ing %s.", altpath))
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
        templ = template.Must(template.ParseFile(tpath))
        Debug(0, fmt.Sprintf("scanning: %s", tpath))
        Debug(1, fmt.Sprintf("context:\n%v", dict))
    }

    buff := bytes.NewBuffer(make([]byte, 0, 1<<20))
    errTExec := templ.Execute(buff, combined(dict, extraData(filename)))
    return buff.String(), errTExec
}

//  Given a filename and dictionary context, create a context dict+("file"=>filename),
//  and read a template specified by relpath. See GetTemplatePath().
func ParseTemplate(filename string, dict map[string]string, relpath []string) (string, os.Error) {
    var templ *template.Template
    if tpath := GetTemplatePath(relpath); tpath == "" {
        return "", NoTemplateError
    } else {
        templ = template.Must(template.ParseFile(tpath))

        Debug(0, fmt.Sprintf("scanning: %s", tpath))
        Debug(1, fmt.Sprintf("context:\n%v", dict))
    }

    buff := bytes.NewBuffer(make([]byte, 0, 1<<20))
    errTExec := templ.Execute(buff, combined(dict, extraData(filename)))
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
        Verbose(fmt.Sprintf("Using alternate template %s\n", GetAltTemplatePath(relpath)))
    } else if stdt, err := ParseTemplate(filename, dict, relpath); err == nil {
        templ = stdt
    } else {
        return err
    }

    Verbose(fmt.Sprintf("Creating %s %s\n", desc, filename))
    Debug(2, fmt.Sprint("\n", templ, "\n"))

    templout := make([]byte, len(templ))
    copy(templout, templ)
    return ioutil.WriteFile(filename, templout, FilePermissions)
}
func AppendTemplate(filename, desc string, dict map[string]string, relpath ...string) os.Error {
    var templ string
    if altt, err := ParseAltTemplate(filename, dict, relpath); err == nil {
        templ = altt
        Verbose(fmt.Sprintf("Using alternate template %s\n", GetAltTemplatePath(relpath)))
    } else if stdt, err := ParseTemplate(filename, dict, relpath); err == nil {
        templ = stdt
    } else {
        return err
    }

    Verbose(fmt.Sprintf("Appending %s %s\n", desc, filename))
    Debug(2, fmt.Sprint("\n", templ, "\n"))

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
