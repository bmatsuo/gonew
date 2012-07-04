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

type configPropertyError struct {
	Property string
	Index    interface{}
	Err      error
}

func (err configPropertyError) Error() string {
	prefix := err.Property
	if err.Index != nil {
		prefix = fmt.Sprintf("%s[%#v]", err.Index)
	}
	switch err.Err.(type) {
	case configPropertyError:
		return fmt.Sprintf("%s.%v", prefix, err.Err)
	}
	return fmt.Sprintf("%s: %v", prefix, err.Err)
}
func newConfigPropertyError(property string, err error) error {
	return configPropertyError{property, nil, err}
}
func newConfigPropertyIndexError(property string, index interface{}, err error) error {
	return configPropertyError{property, index, err}
}
func newConfigMissingPropertyError(property string) error {
	return configPropertyError{property, nil, fmt.Errorf("no value")}
}
func newConfigInvalidPropertyValueError(property string, value interface{}) error {
	return configPropertyError{property, nil, fmt.Errorf("invalid value: %#v", value)}
}
func newConfigInvalidPropertyNameError(property string, name string) error {
	return configPropertyError{property, nil, fmt.Errorf("invalid name: %q", name)}
}

func tryConfigValidate(v interface{}) error {
	switch v.(type) {
	case configValidator:
		return v.(configValidator).Validate()
	}
	return nil
}

type GonewConfig2 struct {
	Environments      map[string]*EnvironmentConfig
	ExternalTemplates []string
	Projects          map[string]*ProjectConfig
}

func configValidateRoot(v interface{}) error {
	if err := tryConfigValidate(v); err != nil {
		return configPropertyError{"$", nil, err}
	}
	return nil
}
func (config GonewConfig2) Validate() error {
	if config.Environments == nil {
		return newConfigMissingPropertyError("Environments")
	}
	for k, env := range config.Environments {
		if strings.IndexFunc(k, unicode.IsSpace) > -1 {
			return configPropertyError{"Environments", nil,
				newConfigInvalidPropertyNameError("Environments", k)}
		}
		if err := tryConfigValidate(env); err != nil {
			return configPropertyError{"Environments", k, err}
		}
		if env.Inherits != nil {
			for i, k2 := range env.Inherits {
				if _, ok := config.Environments[k2]; !ok {
					return configPropertyError{"Environments", k,
						configPropertyError{".Inherits", i,
							fmt.Errorf("missing Environment: %q", k2)}}
				}
			}
		}
	}

	for i, ext := range config.ExternalTemplates {
		if !strings.HasPrefix(ext, "/") {
			return configPropertyError{"ExternalTemplates", i,
				newConfigValidationError("relative path " + ext)}
		}
		info, err := os.Stat(ext)
		if err != nil {
			return configPropertyError{"ExternalTemplates", i,
				newConfigValidationError(err.Error())}
		}
		if !info.IsDir() {
			return configPropertyError{"ExternalTemplates", i,
				newConfigValidationError("not a directory " + ext)}
		}
	}

	if config.Projects == nil {
		return newConfigMissingPropertyError("Projects")
	}
	for k, project := range config.Projects {
		if strings.IndexFunc(k, unicode.IsSpace) > -1 {
			return configPropertyError{"Projects", nil,
				newConfigInvalidPropertyNameError("Projects", k)}
		}
		if err := tryConfigValidate(project); err != nil {
			return configPropertyError{"Projects", k, err}
		}
		if project.Inherits != nil {
			for i, k2 := range project.Inherits {
				if _, ok := config.Projects[k2]; !ok {
					return configPropertyError{"Projects", k,
						configPropertyError{".Inherits", i,
							fmt.Errorf("missing Project: %q", k2)}}
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
	return configValidateRoot(config)
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
