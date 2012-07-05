// Copyright 2012, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

/*  Filename:    version_control_spec.go
 *  Author:      Bryan Matsuo <bryan.matsuo [at] gmail.com>
 *  Created:     2012-07-03 18:02:20.295289 -0700 PDT
 *  Description: 
 */

import (
	"go-validate"
)

type VersionControlConfig struct {
	Type   string
	Remote map[string]interface{}
}

func (config *VersionControlConfig) Validate() (err error) {
	switch config.Type {
	case "":
	case "git":
	case "hg":
	default:
		err = validate.Invalid(config.Type)
	}
	return
}

func (config *VersionControlConfig) Merge(other *VersionControlConfig) {
	switch {
	case other.Type == "":
		fallthrough
	case other.Type == config.Type:
		if other.Remote == nil {
			return
		}
		if config.Remote == nil {
			config.Remote = make(map[string]interface{}, len(other.Remote))
		}
	default:
		config.Type = other.Type
		if other.Remote == nil {
			config.Remote = nil
			return
		}
		config.Remote = make(map[string]interface{}, len(other.Remote))
	}
	if other.Remote != nil {
		for k := range other.Remote {
			config.Remote[k] = other.Remote[k]
		}
	}
}
