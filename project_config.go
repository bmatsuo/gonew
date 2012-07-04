// Copyright 2012, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

/*  Filename:    project_config.go
 *  Author:      Bryan Matsuo <bryan.matsuo [at] gmail.com>
 *  Created:     2012-07-03 18:20:32.475203 -0700 PDT
 *  Description: 
 */

import (
	"fmt"
	"go-validate"
	"strings"
	"unicode"
)

// A set of uniquely named project types.
type ProjectsConfig map[string]*ProjectConfig

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

type ProjectHooksConfig struct {
	Pre  string // a hook should be more than a string
	Post string
}
