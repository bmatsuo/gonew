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

func (license LicenseType) String() string {
    switch license {
    case NewBSD:
        return "New BSD License"
    }
    return ""
}

func (license LicenseType) TemplateNamePrefix() string {
    switch license {
    case NewBSD:
        return "newbsd"
    }
    return ""
}
