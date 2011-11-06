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
	"io"
	"fmt"
	//"log"
	"bytes"
	//"strings"
	//"io/ioutil"
	"path/filepath"
	"template"
	"runtime"
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

type TemplateMultiSet []*template.Set

//  Call function CollectTemplates on each given root and create a TemplateMultiSet.
func MakeTemplateMultiSet(f template.FuncMap, roots ...string) (ms TemplateMultiSet, err error) {
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
func (ms TemplateMultiSet) Execute(wr io.Writer, name string, data interface{}) error {
	for _, s := range ms {
		if t := s.Template(name); t != nil {
			return t.Execute(wr, data)
		}
	}
	return NoTemplateError
}

//  Parse templates <root>/*/*.t, allowing them a given function map.
func CollectTemplates(root string, f template.FuncMap) (s *template.Set, err error) {
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

//  The template directory of the goinstall'ed gonew package.
func GetTemplateRoot() []string {
	return []string{runtime.GOROOT(), "src", "pkg",
		"github.com", "bmatsuo", "gonew", "templates"}
}
