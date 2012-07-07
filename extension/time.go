// Copyright 2012, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package extension

/*  Filename:    time.go
 *  Author:      Bryan Matsuo <bryan.matsuo [at] gmail.com>
 *  Created:     2012-07-06 21:10:38.380735 -0700 PDT
 *  Description: 
 */

import (
	"time"
)

var timeFuncs = Register(Time(time.Now()))

type Time time.Time

func (t Time) Namespace() string { return "Time" }
func (t Time) Now(format string) string { return time.Time(t).Format(format) }
func (t Time) String() string { return time.Time(t).String() }
func (t Time) ANSIC() string { return time.Time(t).Format(time.ANSIC) }
func (t Time) UnixDate() string { return time.Time(t).Format(time.UnixDate) }
func (t Time) RubyDate() string { return time.Time(t).Format(time.RubyDate) }
func (t Time) RFC822() string { return time.Time(t).Format(time.RFC822) }
func (t Time) RFC822Z() string { return time.Time(t).Format(time.RFC822Z) }
func (t Time) RFC850() string { return time.Time(t).Format(time.RFC850) }
func (t Time) RFC1123() string { return time.Time(t).Format(time.RFC1123) }
func (t Time) RFC1123Z() string { return time.Time(t).Format(time.RFC1123Z) }
func (t Time) RFC3339() string { return time.Time(t).Format(time.RFC3339) }
func (t Time) RFC3339Nano() string { return time.Time(t).Format(time.RFC3339Nano) }
func (t Time) Kitchen() string { return time.Time(t).Format(time.Kitchen) }
func (t Time) Stamp() string { return time.Time(t).Format(time.Stamp) }
func (t Time) StampMilli() string { return time.Time(t).Format(time.StampMilli) }
func (t Time) StampMicro() string { return time.Time(t).Format(time.StampMicro) }
func (t Time) StampNano() string { return time.Time(t).Format(time.StampNano) }
