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
	"os"
)

//  A WriteMode determines how a Context should be written.
type WriteMode uint

//  A Context can either be written by being appended to a file, or creating
//  the file and truncating the file.
const (
	AppendMode WriteMode = iota
	CreateMode
)

var writeModeStrings = []string{
	AppendMode: "Appending",
	CreateMode: "Creating",
}

func (wm WriteMode) String() string { return writeModeStrings[wm] }

//  Write a Context to the file it specifies using a specified WriteMode.
//  WriteContext uses template named tname in the templates ExecutorSet. The
//  value desc is only used for printing debugging information.
func WriteContext(context Context, mode WriteMode, tname, desc string) error {
	filename := context.Filename()
	Verbose(fmt.Sprintf("WriteContext: %s %s %s\n", mode.String(), desc, filename))
	Debug(1, fmt.Sprintf("Executing template %s", tname))

	// Execute the template.
	p, errExec := SetExecuted(Templates, tname, context)
	if errExec != nil {
		return errExec
	}
	Debug(2, fmt.Sprintf("\n%s\n", p))

	// Open the output file.
	openMode := os.O_WRONLY
	if mode == AppendMode {
		openMode |= os.O_APPEND
	} else {
		openMode |= os.O_CREATE | os.O_TRUNC
	}
	fout, err := os.OpenFile(filename, openMode, FilePermissions)
	if err != nil {
		return err
	}

	// Write out the executed template and close the file
	if _, err := fout.Write(p); err != nil {
		return err
	}
	if err := fout.Close(); err != nil {
		return err
	}
	return nil
}

//  Write the file specified in c using the template tname.
func CreateFile(c Context, tname string) (err error) {
	mainWriteMode := CreateMode
	debugdesc := c.DebugDescription()

	// Analyze license and file type.
	ltype := c.LicenseType()
	license := ltype.TemplateName(c.FileType())
	pos := ltype.Position(c.FileType())
	showLicense := pos != 0 && license != ""

	// Write out a license header for the file.
	if showLicense && pos < 0 {
		err = WriteContext(c, CreateMode, license, debugdesc+" license")
		if err != nil {
			return
		}
		mainWriteMode = AppendMode
	}

	// Write the main file contents.
	err = WriteContext(c, mainWriteMode, tname, debugdesc)
	if err != nil {
		return
	}

	// Write out a license footer for the file.
	if showLicense && pos > 0 {
		err = WriteContext(c, AppendMode, license, debugdesc+" license")
		if err != nil {
			return
		}
	}

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
