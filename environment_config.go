// Copyright 2012, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

/*  Filename:    environment_spec.go
 *  Author:      Bryan Matsuo <bryan.matsuo [at] gmail.com>
 *  Created:     2012-07-03 18:04:29.869451 -0700 PDT
 *  Description: 
 */

import (
	"fmt"
	"go-validate"
	"strings"
	"unicode"
)

// A set of uniquely named environments.
type EnvironmentsConfig map[string]EnvironmentConfig

// Validates the following for each environment
//		- The name must not contain spaces
//		- The environment must be valid (see EnvironmentConfig.Validate)
//		- All inherited environments must exist.
func (config EnvironmentsConfig) Validate() (err error) {
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
	return
}

// User (project author) details.
type EnvironmentUserConfig struct {
	Name  string
	Email string
}

// Project License configuration.
type LicenseConfig string

// Can be "newbsd", "none", or missing.
func (config LicenseConfig) Validate() (err error) {
	switch config {
	case "newbsd":
	case "none":
		fallthrough
	case "":
	default:
		err = validate.Invalid(config)
	}
	return
}

// Specifies the environment for template generation.
type EnvironmentConfig struct {
	Inherits       []string
	User           *EnvironmentUserConfig
	License        LicenseConfig
	VersionControl *VersionControlConfig
}

// Requires a User. Requires License and VersionControl to be valid
// (see LicenseConfig.Validate, VersionControlConfig.Validate).
func (config EnvironmentConfig) Validate() (err error) {
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
	if err = validate.Property("License", config.License); err != nil {
		return
	}
	if err = validate.Property("VersionControl", config.VersionControl); err != nil {
		return
	}
	return nil
}
