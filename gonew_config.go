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
	"go-validate"
	"io/ioutil"
	"os"
	"strings"
)

type ExternalTemplateConfig string

func (config ExternalTemplateConfig) Validate() (err error) {
	path := string(config)
	var info os.FileInfo
	if !strings.HasPrefix(path, "/") {
		return errors.New("relative path " + path)
	} else if info, err = os.Stat(path); err != nil {
		return
	} else if !info.IsDir() {
		return errors.New("not a directory " + path)
	}
	return
}

type GonewConfig2 struct {
	Environments      EnvironmentsConfig
	ExternalTemplates []ExternalTemplateConfig
	Projects          ProjectsConfig
}

func (config GonewConfig2) Validate() (err error) {
	err = validate.PropertyFunc("Environments", func() (err error) {
		if config.Environments == nil {
			return fmt.Errorf("missing")
		}
		if len(config.Environments) == 0 {
			return fmt.Errorf("empty")
		}
		return
	})
	if err == nil {
		err = validate.Property("Environments", config.Environments)
	}
	if err != nil {
		return
	}

	err = validate.PropertyFunc("ExternalTemplates", func() (err error) {
		for i, ext := range config.ExternalTemplates {
			if err = validate.Index(i, ext); err != nil {
				return
			}
		}
		return
	})
	if err != nil {
		return
	}

	err = validate.PropertyFunc("Projects", func() (err error) {
		if config.Projects == nil {
			err = errors.New("missing")
		}
		return
	})
	if err == nil {
		err = validate.Property("Projects", config.Projects)
	}
	return
}

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

func (config *GonewConfig2) unmarshalJSON(p []byte) error {
	if err := json.Unmarshal(p, config); err != nil {
		return err
	}
	return validate.Property("$", config)
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
