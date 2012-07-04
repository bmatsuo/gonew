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
)

type ProjectConfig struct {
	Inherits []string
	Hooks    map[string]string // a hook should be more than a string
	Files    map[string]*ProjectFileConfig
}

func (config *ProjectConfig) Validate() error {
	if config.Hooks == nil {
		config.Hooks = make(map[string]string)
	}
	if config.Files != nil {
		for k, file := range config.Files {
			if err := tryConfigValidate(file); err != nil {
				return fmt.Errorf("Files[%q]%v", k, err)
			}
		}
	}
	return nil
}
