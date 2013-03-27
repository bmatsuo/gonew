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

// A file context. Its values provide data for gonew templates.
type Context interface {
	// The file's type.
	FileType() FileType
	// The filename being written out to.
	Filename() string
	// The package (basename) the file belongs to.
	Package() string

	// The Project associated with the file.
	Project() Project

	// The file's description for documentation.
	Description() string
	// The names of test functions. Should only be used for generating test files.
	Tests() []string

	// TODO this is probably not necessary.
	// The LicenseType for generating file licenses.
	LicenseType() LicenseType
	// TODO this is probably not necessary
	// Packages source files (non-test, non-options) should import.
	Imports() []string

	// A description of the context for printing debug information. Not useful in templates.
	DebugDescription() string
}
