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
	"strings"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"text/template"
)

// The file extension used by Gonew template files. Custom templates must all
// use this extension.
const (
	TemplateFileExt = ".t"
)

// Template errors
var (
	ErrNoTemplates = errors.New("No valid templates found")
	ErrNoExist     = errors.New("Requested template does not exist")
	ErrParse       = errors.New("Couldn't parse template")
)

func TemplateType(name string) FileType {
	i := strings.Index(name, ".")
	if i < 0 {
		panic(fmt.Sprintf("bad template %s", name))
	}
	switch t := name[:i]; strings.ToLower(t) {
	case "go":
		return GoFile
	case "readme":
		return ReadmeFile
	case "makefile":
		return MakeFile
	case "license":
		return LicenseFile
	case "other":
		return OtherFile
	}
	panic(fmt.Sprintf("unknown template type: %s", name[:i]))
}

// An abstraction of the template.Set type.
type Executor interface {
	Execute(io.Writer, string, interface{}) error
}

//  Call s.Execute with the given name data, writing to a []byte buffer.
func Executed(s Executor, name string, data interface{}) ([]byte, error) {
	b := new(bytes.Buffer)
	err := s.Execute(b, name, data)
	return b.Bytes(), err
}

//  Returns an os-specific pattern <root>/*/*<TemplateFileExt>
func TemplateGlobPattern(root string) string {
	return filepath.Join(root, "*", fmt.Sprintf("*%s", TemplateFileExt))
}

// A linear hierarchy of template (sets).
type TemplateHierarchy []*template.Template

// Note: this should never ignore an error if the template package is working 'properly'
func emptyTemplate(name string) *template.Template { t, _ := template.New(name).Parse(""); return t }

// Create a template multiset containing one template.Template for each root
// diven. Template precedence decreases as roots go from left to right.
func makeTemplateHierarchy(f template.FuncMap, roots ...string) (ms TemplateHierarchy, err error) {
	var s *template.Template
	ms = make(TemplateHierarchy, 0, len(roots))
	for i := range roots {
		s, err = emptyTemplate(roots[i]).Funcs(f).ParseGlob(TemplateGlobPattern(roots[i]))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		} else {
			ms = append(ms, s)
		}
	}
	if len(ms) == 0 {
		err = ErrNoTemplates
	}
	return
}

// Collect all visible sets of templates visible.
func FindTemplates() (TemplateHierarchy, error) {
	troots := make([]string, 0, 2)
	if alt := AppConfig.AltRoot; alt != "" {
		troots = append(troots, alt)
	}
	troots = append(troots, TemplateRoot)
	return makeTemplateHierarchy(DefaultFuncMap(), troots...)
}

//  Execute the named template from the first set in which such a template exists.
func (ms TemplateHierarchy) Execute(wr io.Writer, name string, data interface{}) error {
	for _, s := range ms {
		if t := s.Lookup(name); t != nil {
			return t.Execute(wr, data)
		}
	}
	return ErrNoTemplates
}
