// Copyright 2011, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

/*  Filename:    context.go
 *  Author:      Bryan Matsuo <bryan.matsuo@gmail.com>
 *  Created:     Sat Oct 22 23:58:07 PDT 2011
 *  Description: 
 */

import (
	"fmt"
	"io/ioutil"
)

//  Write the file specified in c using the template tname.
func CreateFile(c Context, tname string) (err error) {
	filename := c.Filename()
	Verbose(fmt.Sprintf("WriteContext: Creating %s %s\n", c.DebugDescription(), filename))

	p, err := generateTemplate(Templates, tname, c)
	if err != nil {
		return err
	}
	ioutil.WriteFile(filename, p, FilePermissions)

	// TODO: check the new file into git under certain conditions...

	return
}

//  A Context values are used as the data for gonew templates. The methods
//  can be called from a template with ".<method>". For example
//      This file is '{{ .Filename }}'
type Context interface {
	// The filename to which the Context is written out to.
	Filename() string
	// The name of the go package the Context belongs to. For commands, this
	// should return "main". For packages, the package/target name.
	Package() string
	// The names of packages that source (non-test, non-options) .go files
	// should import.
	Imports() []string
	// The description of the Context. This should be unique for each file
	// and should describe the file's purpose.
	Description() string
	// The names of tests for the Context. Really only used for generating
	// test files
	Tests() []string
	// The FileType of the Context. This is used when generating license text.
	FileType() FileType
	// The LicenseType for generating file licenses.
	LicenseType() LicenseType
	// A description of the file for printing debug information. This is not
	// very useful in templates.
	DebugDescription() string
	// The Project associated with the file. This is the zero Project for
	// any library files (generated with the command gonew lib NAME PKG)
	Project() Project
}
