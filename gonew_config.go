// Copyright 2012, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

/*  Filename:    gonew_config.go
 *  Author:      Bryan Matsuo <bryan.matsuo [at] gmail.com>
 *  Created:     2012-07-03 18:18:19.325777 -0700 PDT
 *  Description: 
 */

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"unicode"
)

type configValidator interface {
	Validate() error
}

func newConfigValidationError(msg string) error { return fmt.Errorf(": %s", msg) }
func newConfigMissingPropertyError(property string) error {
	return newConfigValidationError(fmt.Sprintf("missing property: %q", property))
}
func newConfigInvalidPropertyError(property string, value interface{}) error {
	return newConfigValidationError(fmt.Sprintf("invalid property: %q (%#v)", property, value))
}
func newConfigInvalidNameError(nameType string, name string) error {
	return newConfigValidationError(fmt.Sprintf("invalid %s name: %q", nameType, name))
}

func tryConfigValidate(v interface{}) error {
	switch v.(type) {
	case configValidator:
		return v.(configValidator).Validate()
	}
	return nil
}

type GonewConfig2 struct {
	Environments map[string]*EnvironmentConfig
	Projects     map[string]*ProjectConfig
}

func (config GonewConfig2) Validate() error {
	if config.Environments == nil {
		return fmt.Errorf("$%v", newConfigMissingPropertyError("Environments"))
	}
	for k, env := range config.Environments {
		if strings.IndexFunc(k, unicode.IsSpace) > -1 {
			return fmt.Errorf("$.Environments%v", newConfigInvalidNameError("Environment", k))
		}
		if err := tryConfigValidate(env); err != nil {
			return fmt.Errorf("$.Environments[%q]%v", k, err)
		}
		if env.Inherits != nil {
			for i, k2 := range env.Inherits {
				if _, ok := config.Environments[k2]; !ok {
					return fmt.Errorf("$.Environments[%q].Inherits[%d]", k, i, fmt.Errorf("missing Environment: %q", k2))
				}
			}
		}
	}

	if config.Projects == nil {
		return fmt.Errorf("$%v", newConfigMissingPropertyError("Projects"))
	}
	for k, project := range config.Projects {
		if strings.IndexFunc(k, unicode.IsSpace) > -1 {
			return fmt.Errorf("$.Projects%v", newConfigInvalidNameError("Project", k))
		}
		if err := tryConfigValidate(project); err != nil {
			return fmt.Errorf("$.Projects[%q]%v", k, err)
		}
		if project.Inherits != nil {
			for i, k2 := range project.Inherits {
				if _, ok := config.Projects[k2]; !ok {
					return fmt.Errorf("$.Projects[%q].Inherits[%d]", k, i, fmt.Errorf("missing Project: %q", k2))
				}
			}
		}
	}
	return nil
}

// Implements json.Marshaler
func (config *GonewConfig2) marshalJSON() ([]byte, error) { return json.Marshal(config) }
func (config *GonewConfig2) MarshalFileJSON(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	p, err := config.marshalJSON()
	if err != nil {
		return err
	}
	_, err = f.Write(p)
	return err
}

// Implements json.Unmarshaler
func (config *GonewConfig2) unmarshalJSON(p []byte) error {
	if err := json.Unmarshal(p, config); err != nil {
		return err
	}
	return tryConfigValidate(config)
}
func (config *GonewConfig2) UnmarshalFileJSON(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	p, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	return config.unmarshalJSON(p)
}
