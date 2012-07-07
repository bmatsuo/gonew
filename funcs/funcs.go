// Copyright 2012, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package funcs

/*  Filename:    funcs.go
 *  Author:      Bryan Matsuo <bryan.matsuo [at] gmail.com>
 *  Created:     2012-07-06 21:03:15.802747 -0700 PDT
 *  Description: 
 */

import (
	"text/template"
)

type Interface interface {
	Namespace() string
	FuncMap() template.FuncMap
}

var Funcs = template.FuncMap{
}

func Register(funcs Interface) Interface {
	funcmap := funcs.FuncMap()
	if funcmap == nil {
		return funcs
	}
	if Funcs == nil {
		Funcs = make(template.FuncMap, len(funcmap))
	}
	for name, fn := range funcmap {
		Funcs[funcs.Namespace() + name] = fn
	}
	return funcs
}
