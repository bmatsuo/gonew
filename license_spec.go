// Copyright 2012, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

/*  Filename:    license_spec.go
 *  Author:      Bryan Matsuo <bryan.matsuo [at] gmail.com>
 *  Created:     2012-06-30 22:46:45.509005 -0700 PDT
 *  Description: 
 */

import ()

var Licenses = make(map[string]*LicenseSpec)

func registerLicense(license *LicenseSpec) *LicenseSpec {
	if _, ok := Licenses[license.Type]; ok {
		panic("duplicate license type: " + license.Type)
	}
	Licenses[license.Type] = license
	return license
}

var NewBSDSpec = registerLicense(&LicenseSpec{
	Type: "newbsd",
	TemplateSpecs: map[string]*LicenseTemplateSpec{
		"go":      {Head: []string{"license.newbsd.gohead.t"}},
		"readme":  {Foot: []string{"license.newbsd.readme.t"}},
		"license": {Head: []string{"license.newbsd.t"}},
	},
})

type LicenseSpec struct {
	Type          string
	TemplateSpecs map[string]*LicenseTemplateSpec
}

func (license *LicenseSpec) Head(fileType string) []string {
	templateSpec, ok := license.TemplateSpecs[fileType]
	if !ok {
		return nil
	}
	return templateSpec.Head
}

func (license *LicenseSpec) Foot(fileType string) []string {
	templateSpec, ok := license.TemplateSpecs[fileType]
	if !ok {
		return nil
	}
	return templateSpec.Foot
}

type LicenseTemplateSpec struct{ Head, Foot []string }
