// Copyright 2012, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package funcs

/*  Filename:    time.go
 *  Author:      Bryan Matsuo <bryan.matsuo [at] gmail.com>
 *  Created:     2012-07-06 21:10:38.380735 -0700 PDT
 *  Description: 
 */

import (
	"text/template"
	"time"
)

var timeFuncs = Register(Time(time.Now()))

type Time time.Time

func (t Time) Namespace() string { return "time" }
func (t Time) FuncMap() template.FuncMap {
	return template.FuncMap{
		"Now":         func(format string) string { return time.Time(t).Format(format) },
		"String":      func() string { return time.Time(t).String() },
		"ANSIC":       func() string { return time.Time(t).Format(time.ANSIC) },
		"UnixDate":    func() string { return time.Time(t).Format(time.UnixDate) },
		"RubyDate":    func() string { return time.Time(t).Format(time.RubyDate) },
		"RFC822":      func() string { return time.Time(t).Format(time.RFC822) },
		"RFC822Z":     func() string { return time.Time(t).Format(time.RFC822Z) },
		"RFC850":      func() string { return time.Time(t).Format(time.RFC850) },
		"RFC1123":     func() string { return time.Time(t).Format(time.RFC1123) },
		"RFC1123Z":    func() string { return time.Time(t).Format(time.RFC1123Z) },
		"RFC3339":     func() string { return time.Time(t).Format(time.RFC3339) },
		"RFC3339Nano": func() string { return time.Time(t).Format(time.RFC3339Nano) },
		"Kitchen":     func() string { return time.Time(t).Format(time.Kitchen) },
		"Stamp":       func() string { return time.Time(t).Format(time.Stamp) },
		"StampMilli":  func() string { return time.Time(t).Format(time.StampMilli) },
		"StampMicro":  func() string { return time.Time(t).Format(time.StampMicro) },
		"StampNano":   func() string { return time.Time(t).Format(time.StampNano) },
	}
}
