// Copyright 2012, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package config

/*  Filename:    project_config.go
 *  Author:      Bryan Matsuo <bryan.matsuo [at] gmail.com>
 *  Created:     2012-07-03 18:20:32.475203 -0700 PDT
 *  Description: 
 */

import (
	"errors"
	"fmt"
	"github.com/bmatsuo/go-validate"
	"strings"
	"unicode"
)

// A set of uniquely named project types.
type ProjectsConfig map[string]*ProjectConfig

func (config ProjectsConfig) inheritanceGraph() configInheritanceGraph {
	g := make(configInheritanceGraph, len(config))
	for v := range config {
		g[v] = make(map[string]interface{}, len(config[v].Inherits))
		for _, w := range config[v].Inherits {
			g[v][w] = true
		}
	}
	return g
}

// Validates the following for each project
//		- The name must not contain spaces
//		- The environment must be valid (see EnvironmentConfig.Validate)
//		- All inherited environments must exist.
func (config ProjectsConfig) Validate() (err error) {
	for k, project := range config {
		if strings.IndexFunc(k, unicode.IsSpace) > -1 {
			return validate.Invalid("name", k)
		}
		if err = validate.Property(k, project); err != nil {
			return
		}
		err = validate.IndexFunc(k, func() (err error) {
			err = validate.PropertyFunc("Inherits", func() (err error) {
				for i, k2 := range project.Inherits {
					err = validate.IndexFunc(i, func() (err error) {
						if _, ok := config[k2]; !ok {
							err = fmt.Errorf("unknown project: %q", k2)
						}
						return
					})
					if err != nil {
						return
					}
				}
				return
			})
			return
		})
		if err != nil {
			return
		}
	}
	graph := config.inheritanceGraph()
	for k := range config {
		err = validate.IndexFunc(k, func() (err error) {
			if b, _ := graph.HasCycles(k); b {
				err = errors.New("inheritance cycle")
			}
			return
		})
		if err != nil {
			return
		}
	}
	return
}

type ProjectConfig struct {
	Inherits []string
	Hooks    *ProjectHooksConfig
	Files    map[string]*ProjectFileConfig
}

func (config *ProjectConfig) Validate() (err error) {
	if config.Hooks == nil {
		config.Hooks = new(ProjectHooksConfig)
	}
	if config.Files == nil {
		config.Files = make(map[string]*ProjectFileConfig)
	}
	err = validate.Property("Files", func() (err error) {
		for k, file := range config.Files {
			if err = validate.Index(k, file); err != nil {
				return
			}
		}
		return
	})
	return
}

func (config *ProjectConfig) Merge(other *ProjectConfig) {
	if other.Hooks != nil {
		if config.Hooks == nil {
			config.Hooks = new(ProjectHooksConfig)
		}
		config.Hooks.Merge(other.Hooks)
	}
	if other.Files != nil {
		if config.Files == nil {
			config.Files = make(map[string]*ProjectFileConfig, len(other.Files))
		}
		for name, otherFile := range other.Files {
			file, present := config.Files[name]
			if !present {
				file = new(ProjectFileConfig)
				config.Files[name] = file
			}
			file.Merge(otherFile)
		}
	}
}

type ProjectHooksConfig struct {
	Pre, Post []*HookConfig
}

func (config *ProjectHooksConfig) Merge(other *ProjectHooksConfig) {
	if other.Pre != nil {
		config.Pre = append(other.Pre, config.Pre...)
	}
	if other.Post != nil {
		config.Post = append(other.Post, config.Post...)
	}
}

type HookConfig struct {
	Cwd      string
	Commands []string
}

func (config *HookConfig) Merge(other *ProjectHooksConfig) {
}
