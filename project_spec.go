// Copyright 2012, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

/*  Filename:    project_spec.go
 *  Author:      Bryan Matsuo <bryan.matsuo [at] gmail.com>
 *  Created:     2012-06-30 22:28:04.819539 -0700 PDT
 *  Description: 
 */

import (
	"path/filepath"
)

type UserSpec struct { Name, Email string }

type ProjectSpec struct {
	Prefix         string
	Package        string
	Lib            string
	User           *UserSpec
	VersionControl *VersionControlSpec
	License        *LicenseSpec
	Files          []*FileSpec
}

// Implement ProjectDefinition.
func (project *ProjectSpec) ProjectPrefix() string {
	return filepath.Join(project.Prefix, project.Package)
}

// Implement ProjectDefinition.
func (project *ProjectSpec) ProjectPackage() string { return project.Package }

// Implement ProjectDefinition.
func (project *ProjectSpec) ProjectLib() string { return project.Lib }

// Compile Files for generation.
func (project *ProjectSpec) compile() []*FileSpec {
	files := make([]*FileSpec, len(project.Files))
	for i, file := range project.Files {
		// Create a compiled copy of the file.
		file = &FileSpec{
			Path:      file.ProjectPath(project),
			Type:      file.Type,
			Templates: file.Templates,
		}
		files[i] = file

		// Add license-specific templates.
		if head := project.License.Head(file.Type); head != nil {
			file.Templates = append(head, file.Templates...)
		}
		if foot := project.License.Foot(file.Type); foot != nil {
			file.Templates = append(file.Templates, foot...)
		}
	}
	return files
}
