// Copyright 2012, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The templating system for gonew (v2).
package templates

/*  Filename:    templates.go
 *  Author:      Bryan Matsuo <bryan.matsuo [at] gmail.com>
 *  Created:     2012-07-05 23:02:29.719822 -0700 PDT
 *  Description: 
 */

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"text/template"
)

// Output file creation helper. Creates any missing parent directories. Does not
// overwrite existing files.
func FileCreate(path string) (*os.File, error) {
	if _, err := os.Stat(path); err == nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0644); err != nil { // TODO configurable? smarter?
		return nil, err
	}
	return os.Create(path)
}

// An artificial wrapper struct to make template rendering less cumbersome.
type Environment struct{ v interface{} }

// Create an Environment.
func Env(v interface{}) Environment { return Environment{v} }

// Render named templates to `out`. Halts if a rendering error is encountered.
func (env Environment) Render(out io.Writer, ts Interface, names ...string) (err error) {
	for _, name := range names {
		if err = ts.Render(out, name, env); err != nil {
			return
		}
	}
	return
}

type ErrSourceType struct{ v interface{} }
type ErrNoTemplate string

func (err ErrSourceType) Error() string { return fmt.Sprint("not a source: ", err.v) }
func (err ErrNoTemplate) Error() string { return "no template: " + string(err) }

// A set of templates relatively in-line with template.Template.
type Interface interface {
	Render(io.Writer, string, interface{}) error // Render a named template.
	Source(interface{}) error                    // Add a template source.
	Funcs(template.FuncMap) error                // Add a set of functions.
}

// The straight-forward implementation of Interface.
type templates struct {
	t   *template.Template
	ext string
}

// Create a new template set that recognizes ext as a template file extension.
func New(ext string) Interface { return (&templates{ext: ext}).setup() }

func (ts *templates) Render(out io.Writer, name string, environment interface{}) error {
	if ts.t == nil {
		return ErrNoTemplate(name)
	}
	return ts.t.ExecuteTemplate(out, name, environment)
}

func (ts *templates) setup() *templates {
	if ts.t == nil {
		fns := template.FuncMap{"gonew": func() string { return "gonew v2" }}
		ts.t = template.Must(template.New("gonew").Funcs(fns).Parse("{{gonew}}"))
	}
	return ts
}

func (ts *templates) Source(src interface{}) (err error) {
	switch src.(type) {
	case string:
		dir := src.(string)
		if !isDir(dir) {
			return fmt.Errorf("not a directory: %s", dir)
		}
		var paths []string
		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				paths = append(paths, path)
			}
			return nil
		})
		_, err = ts.setup().t.ParseFiles(paths...)
	case template.Template:
		t := src.(template.Template)
		_, err = ts.setup().t.AddParseTree(t.Name(), t.Tree)
	default:
		err = ErrSourceType{src}
	}
	return
}

func (ts *templates) Funcs(fns template.FuncMap) error {
	ts.setup().t.Funcs(fns)
	return nil
}

func isDir(d string) bool { // be careful
	info, err := os.Stat(d)
	return err == nil && info.IsDir()
}
