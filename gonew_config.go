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
	"errors"
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
		prefix = fmt.Sprintf("%s[%#v]", prefix, err.Index)
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
func configPropertyDo(property string, validate func() (interface{}, error)) error {
	if key, err := validate(); err != nil {
		return configPropertyError{property, key, err}
	}
	return nil
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
func (config GonewConfig2) Validate() (err error) {
	err = configPropertyDo("Environments", func() (key interface{}, err error) {
		if config.Environments == nil {
			return nil, fmt.Errorf("missing")
		}
		if len(config.Environments) == 0 {
			return nil, fmt.Errorf("empty")
		}
		for k, env := range config.Environments {
			if strings.IndexFunc(k, unicode.IsSpace) > -1 {
				key, err = nil, fmt.Errorf("invalid name: %v", k)
				return
			}
			if err = tryConfigValidate(env); err != nil {
				key = nil
				return
			}
			if env.Inherits != nil {
				for i, k2 := range env.Inherits {
					key, err = k, configPropertyDo("Inherits", func() (key interface{}, err error) {
						if _, ok := config.Environments[k2]; !ok {
							key, err = i, fmt.Errorf("unknown environment: %q", k2)
						}
						return
					})
					if err != nil {
						return
					}
				}
			}
		}
		return
	})
	if err != nil {
		return
	}

	err = configPropertyDo("ExternalTemplates", func() (key interface{}, err error) {
		for i, ext := range config.ExternalTemplates {
			if !strings.HasPrefix(ext, "/") {
				key, err = i, errors.New("relative path " + ext)
				return
			}
			var info os.FileInfo
			info, err = os.Stat(ext)
			if err != nil {
				key = i
				return
			} else if !info.IsDir() {
				key, err = i, errors.New("not a directory " + ext)
			}
		}
		return
	})
	if err != nil {
		return
	}

	err = configPropertyDo("Projects", func() (key interface{}, err error) {
		if config.Projects == nil {
			err = errors.New("missing")
			return
		}
		for k, project := range config.Projects {
			if strings.IndexFunc(k, unicode.IsSpace) > -1 {
				key, err = k, errors.New("invalid name: " + k)
				return
			}
			if err = tryConfigValidate(project); err != nil {
				key = k
				return
			}
			err = configPropertyDo("Inherits", func()(key interface{}, err error) {
				if project.Inherits != nil {
					for i, k2 := range project.Inherits {
						if _, ok := config.Projects[k2]; !ok {
							key, err = i, fmt.Errorf("missing project: %q", k2)
							return
						}
					}
				}
				return
			})
			if err != nil {
				key = k
				return
			}
		}
		return
	})
	return
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
