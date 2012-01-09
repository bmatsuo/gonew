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
	"errors"
	"fmt"
	"io"
	"os"
	//"log"
	"bytes"
	//"strings"
	//"io/ioutil"
	"path/filepath"
	"runtime"
	"text/template"
	//"github.com/hoisie/mustache.go"
)

const (
	TemplateFileExt = ".t"
)

var (
	NoTemplateError = errors.New("Requested template does not exist")
	ParseError      = errors.New("Couldn't parse template")
)

type Executor interface {
	Execute(io.Writer, interface{}) error
}

type ExecutorSet interface {
	Execute(io.Writer, string, interface{}) error
}

func collectBytes(f func(w io.Writer) error) (p []byte, err error) {
	b := new(bytes.Buffer)
	err = f(b)
	p = b.Bytes()
	return
}

//  Call t.Execute with the given data, writing to a []byte buffer.
func Executed(t Executor, data interface{}) ([]byte, error) {
	return collectBytes(func(w io.Writer) error { return t.Execute(w, data) })
}

//  Call s.Execute with the given name data, writing to a []byte buffer.
func ExecutedSet(s ExecutorSet, name string, data interface{}) ([]byte, error) {
	return collectBytes(func(w io.Writer) error { return s.Execute(w, name, data) })
}

//  Returns an os-specific pattern <root>/"*"/*<TemplateFileExt>
func TemplateGlobPattern(root string) string {
	return filepath.Join(root, "*", fmt.Sprintf("*%s", TemplateFileExt))
}

type TemplateMultiSet []*template.Template

//  Call function CollectTemplates on each given root and create a TemplateMultiSet.
func MakeTemplateMultiSet(f template.FuncMap, roots ...string) (ms TemplateMultiSet, err error) {
	var s *template.Template
	ms = make(TemplateMultiSet, 0, len(roots))
	for i := range roots {
		if s, err = CollectTemplates(roots[i], f); err != nil {
			fmt.Fprintln(os.Stderr, err)
		} else {
			ms = append(ms, s)
		}
	}
	if len(ms) == 0 {
		err = errors.New("No valid templates found")
	}
	return
}

//  Execute a template in first s in ms for which s.Template(name) is non-nil.
func (ms TemplateMultiSet) Execute(wr io.Writer, name string, data interface{}) error {
	for _, s := range ms {
		if t := s.Lookup(name); t != nil {
			return t.Execute(wr, data)
		}
	}
	return NoTemplateError
}

func emptyTemplate(name string) (*template.Template, error) { return template.New(name).Parse("") }

//  Parse templates <root>/*/*.t, allowing them a given function map.
func CollectTemplates(root string, f template.FuncMap) (s *template.Template, err error) {
	switch s, err = emptyTemplate(root); {
	case err != nil:
		return
	case f != nil:
		s.Funcs(f)
	}
	return s.ParseGlob(TemplateGlobPattern(root))
}

//  The template directory of the goinstall'ed gonew package.
func GetTemplateRoot() []string {
	return []string{runtime.GOROOT(), "src", "pkg",
		"github.com", "bmatsuo", "gonew", "templates"}
}
