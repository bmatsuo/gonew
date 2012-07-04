// Copyright 2012, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

/*  Filename:    file_spec.go
 *  Author:      Bryan Matsuo <bryan.matsuo [at] gmail.com>
 *  Created:     2012-06-30 22:13:24.768101 -0700 PDT
 *  Description: 
 */

import (
	"bytes"
	"text/template"
)

type ProjectDefinition interface {
	ProjectPrefix() string
	ProjectPackage() string
	ProjectLib() string
}

type FileSpec struct {
	Path      string
	Type      string
	Templates []string
}

func (file *FileSpec) ProjectPath(definition ProjectDefinition) string {
	templ := template.Must(template.New("path").Parse(file.Path))
	buff := new(bytes.Buffer)
	if err := templ.Execute(buff, definition); err != nil {
		panic(err)
	}
	return string(buff.Bytes())
}
