// Copyright 2012, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package project

/*  Filename:    project.go
 *  Author:      Bryan Matsuo <bryan.matsuo [at] gmail.com>
 *  Created:     2012-07-06 23:08:16.03525 -0700 PDT
 *  Description: 
 */

import (
	"path"

	"github.com/bmatsuo/gonew/config"
	"github.com/bmatsuo/gonew/extension"
)

var BaseImportPath string

func importPath(pkg string) string {
	if BaseImportPath == "" {
		return pkg
	}
	return path.Join(BaseImportPath, pkg)
}

func Context(filename, filetype string, p Interface) interface{} {
	return map[string]interface{}{
		"File": map[string]interface{}{
			"Name": filename,
			"Type": filetype,
		},
		"Prefix":  p.Prefix(),
		"Package": p.Package(),
		"Project": p,
		"Env":     p.Env(),
		"X":       extension.Extensions,
	}
}

type Interface interface {
	Name() string
	Prefix() string
	Package() string
	Import() string
	Env() *config.EnvironmentConfig
}

func New(name, pkg string, env *config.EnvironmentConfig) Interface {
	return &project{name, pkg, env}
}

type project struct {
	name string
	pkg  string
	env  *config.EnvironmentConfig
}

func (p *project) Name() string                   { return p.name }
func (p *project) Prefix() string                 { return "./" + p.name } // XXX could be smarter
func (p *project) Package() string                { return p.pkg }
func (p *project) Import() string                 { return importPath(p.Package()) }
func (p *project) Env() *config.EnvironmentConfig { return p.env }
