// Copyright 2011, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

/*
 *  Filename:    license.go
 *  Package:     main
 *  Author:      Bryan Matsuo <bmatsuo@soe.ucsc.edu>
 *  Created:     Mon Jul  4 00:53:08 PDT 2011
 *  Description: 
 */
import (
	"fmt"
)

type LicenseType int

const (
	NilLicenseType LicenseType = iota
	NewBSD
	// Apache
	// GNUGPLv3
	// GNULGPLv3
	// ...
)

var licstrings = []string{NilLicenseType: "no license", NewBSD: "New BSD"}
var licprefix = []string{NilLicenseType: "", NewBSD: "newbsd"}

// A string describing the license, or "no license" if lt == NilLicenseType.
func (lt LicenseType) String() string { return licstrings[lt] }

// The prefix that must be present on all template names.
func (lt LicenseType) TemplatePrefix() string { return fmt.Sprintf("%s.%s", "license", licprefix[lt]) }

// The template of the LICENSE file.
func (lt LicenseType) Template() (name string) {
	if lt == NilLicenseType {
		return
	}
	return lt.TemplatePrefix() + TemplateFileExt
}

// The template to be placed as a file header.
func (lt LicenseType) HeaderTemlate(typ FileType) (name string) {
	if lt == NilLicenseType {
		return
	}

	switch typ {
	//case MakeFile:
	//	name = lt.TemplatePrefix() + ".makefile.t"
	case GoFile:
		name = fmt.Sprintf("%s.%s%s", lt.TemplatePrefix(), "gohead", TemplateFileExt)
	}
	return
}

// The template to be placed as a file footer.
func (lt LicenseType) FooterTemplate(typ FileType) (name string) {
	if lt == NilLicenseType {
		return
	}

	switch typ {
	case ReadmeFile:
		name = fmt.Sprintf("%s.%s%s", lt.TemplatePrefix(), "readme", TemplateFileExt)
	}
	return
}

func (lt LicenseType) TemplateName(ftype FileType) string {
	if lt == NilLicenseType {
		return ""
	}

	t := lt.TemplatePrefix()
	switch ftype {
	case ReadmeFile:
		t += ".readme.t"
	case MakeFile:
		t += ".makefile.t"
	case GoFile:
		t += ".gohead.t"
	}
	return t
}

//  Returns -1 if the license appears at the top of the file, 1 if at the
//  bottom, and 0 if there should be no license.
func (lt LicenseType) Position(ftype FileType) (pos int) {
	pos = -1
	switch lt {
	case NewBSD:
		switch ftype {
		case ReadmeFile:
			pos = 1
		case MakeFile:
			pos = 0
		case GoFile:
			pos = -1
		case OtherFile:
			pos = 0
		}
	}
	return
}
