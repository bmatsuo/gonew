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

type configInheritanceGraph map[string]map[string]interface{}

type configInheritanceDFSLog map[string]*struct{ start, finish int }

func makeConfigInheritanceDFSLog(g configInheritanceGraph) configInheritanceDFSLog {
	log := make(configInheritanceDFSLog, len(g))
	for v := range g {
		log[v] = new(struct{ start, finish int })
	}
	return log
}

func (g configInheritanceGraph) HasCycles(start string) (bool, []string) {
	b := false
	dfs := &configInheritanceDFS{
		g: g, log: makeConfigInheritanceDFSLog(g),
		onBack: func(v, w string) error { b = true; return nil },
	}
	dfs.visit(1, start)
	mergeOrder := make([]string, 0, len(dfs.finished))
	mergeOrder = append(mergeOrder, dfs.finished...)
	//for i := len(dfs.finished)-1; i >= 0; i-- {
	//	mergeOrder = append(mergeOrder, dfs.finished[i])
	//}
	return b, mergeOrder
}

func (log configInheritanceDFSLog) State(vertex string) string {
	switch vlog, ok := log[vertex]; {
	case !ok:
		return "invalid vertex"
	case vlog.start == 0:
		return "unvisited"
	case vlog.finish == 0:
		return "unfinished"
	}
	return "finished"
}

type configInheritanceDFS struct {
	g        configInheritanceGraph
	log      configInheritanceDFSLog
	finished []string
	started  []string
	onBack   func(from, to string) error
	onTree   func(from, to string) error
	onCross  func(from, to string) error
	onFinish func(vertex string) error
	onStart  func(vertex string) error
}

func (dfs *configInheritanceDFS) visit(t int, v string) int {
	dfs.log[v].start = t
	dfs.started = append(dfs.started, v)
	if dfs.onStart != nil {
		dfs.onStart(v)
	}

	for w := range dfs.g[v] {
		if dfs.log[w].start == 0 {
			if dfs.onTree != nil {
				dfs.onTree(v, w)
			}
			t = dfs.visit(t+1, w)
			continue
		}
		if dfs.log[w].finish == 0 {
			if dfs.onBack != nil {
				dfs.onBack(v, w)
			}
			continue
		}
		if dfs.onCross != nil {
			dfs.onCross(v, w)
		}
	}

	t++
	dfs.log[v].finish = t
	dfs.finished = append(dfs.finished, v)
	if dfs.onFinish != nil {
		dfs.onFinish(v)
	}
	return t
}

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

func (config GonewConfig2) Environment(name string) (*EnvironmentConfig, error) {
	if _, present := config.Environments[name]; !present {
		return nil, errors.New("unknown environment: " + name)
	}
	_, mergeOrder := config.Environments.inheritanceGraph().HasCycles(name)

	env := new(EnvironmentConfig)
	for _, key := range mergeOrder {
		env.Merge(config.Environments[key])
	}
	return env, nil
}

func (config GonewConfig2) Project(name string) (*ProjectConfig, error) {
	if _, present := config.Projects[name]; !present {
		return nil, errors.New("unknown project: " + name)
	}
	_, mergeOrder := config.Projects.inheritanceGraph().HasCycles(name)

	env := new(ProjectConfig)
	for _, key := range mergeOrder {
		env.Merge(config.Projects[key])
	}
	return env, nil
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
