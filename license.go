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
    "strings"
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

var licstrings = []string{
    NilLicenseType: "Nil",
    NewBSD:         "New BSD",
}
func (lt LicenseType) String() string { return licstrings[lt] }

var licprefix = []string{
    NilLicenseType: "",
    NewBSD:         "newbsd",
}
func (lt LicenseType) TemplateNamePrefix() string {
    return strings.Join([]string{"license", licprefix[lt]}, ".")
}


func (lt LicenseType) FullTemplateName() string {
    if lt == NilLicenseType {
        return ""
    }

    return lt.TemplateNamePrefix() + TemplateFileExt
}

func (lt LicenseType) TemplateName(ftype FileType) string {
    if lt == NilLicenseType {
        return ""
    }

    t := lt.TemplateNamePrefix()
    switch ftype {
    case README:
        t += ".readme.t"
    case Makefile:
        t += ".makefile.t"
    case Go:
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
        case README:
            pos = 1
        case Makefile:
            pos = 0
        case Go:
            pos = -1
        case Other:
            pos = 0
        }
    }
    return
}
