// Copyright 2012, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package config

/*  Filename:    environment_spec.go
 *  Author:      Bryan Matsuo <bryan.matsuo [at] gmail.com>
 *  Created:     2012-07-03 18:04:29.869451 -0700 PDT
 *  Description:
 */

import (
	"errors"
	"fmt"
	"github.com/bmatsuo/go-validate"
	"strings"
	"unicode"
)

// A set of uniquely named environments.
type Environments map[string]*Environment

func (config Environments) inheritanceGraph() configInheritanceGraph {
	g := make(configInheritanceGraph, len(config))
	for v := range config {
		g[v] = make(map[string]interface{}, len(config[v].Inherits))
		for _, w := range config[v].Inherits {
			g[v][w] = true
		}
	}
	return g
}

// Validates the following for each environment
//		- The name must not contain spaces
//		- The environment must be valid (see Environment.Validate)
//		- All inherited environments must exist.
func (config Environments) Validate() (err error) {
	for k, env := range config {
		if strings.IndexFunc(k, unicode.IsSpace) > -1 {
			return validate.Invalid("name", k)
		}
		if err = validate.Property(k, env); err != nil {
			return
		}
		err = validate.IndexFunc(k, func() (err error) {
			err = validate.PropertyFunc("Inherits", func() (err error) {
				for i, k2 := range env.Inherits {
					err = validate.IndexFunc(i, func() (err error) {
						if _, ok := config[k2]; !ok {
							err = fmt.Errorf("unknown environment: %q", k2)
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

// User (project author) details. All fields are optional.
type EnvironmentUserConfig struct {
	Name  string // A real name or pseudonym
	Email string // An email address (potentially malformed)
}

func (config *EnvironmentUserConfig) Merge(other *EnvironmentUserConfig) {
	if other.Name != "" {
		config.Name = other.Name
	}
	if other.Email != "" {
		config.Email = other.Email
	}
}

// Specifies the environment for template generation.
type Environment struct {
	BaseImportPath string                 // Base import path for templates
	Inherits       []string               // Environments to inherit configs from
	User           *EnvironmentUserConfig // User info for templates
}

// Merges other into config. Inherits are not merged, as this is used to eliminate inheritence.
func (config *Environment) Merge(other *Environment) {
	if other.BaseImportPath != "" {
		config.BaseImportPath = other.BaseImportPath
	}
	if other.User != nil {
		if config.User == nil {
			config.User = new(EnvironmentUserConfig)
		}
		config.User.Merge(other.User)
	}
}

// Requires a User.
func (config *Environment) Validate() (err error) {
	err = validate.PropertyFunc("User", func() (err error) {
		if config.User == nil {
			err = fmt.Errorf("missing")
		}
		return err
	})
	if err != nil {
		return
	}
	if err = validate.Property("User", config.User); err != nil {
		return
	}
	return nil
}
