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
)

type EnvironmentUserConfig struct {
	Name  string
	Email string
}

type EnvironmentConfig struct {
	Inherits	   []string
	User           *EnvironmentUserConfig
	License        string
	VersionControl *VersionControlConfig
}

func (config EnvironmentConfig) Validate() error {
	if config.User == nil {
		newConfigMissingPropertyError("User")
	}
	if err := tryConfigValidate(config.User); err != nil {
		return fmt.Errorf("User%v", err)
	}
	switch (config.License) {
	case "":
		fallthrough
	case "none":
	case "newbsd":
	default:
		return newConfigInvalidPropertyError("License", config)
	}
	if err := tryConfigValidate(config.VersionControl); err != nil {
		return fmt.Errorf("VersionControl%v", err)
	}
	return nil
}
