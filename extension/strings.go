// Copyright 2014, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// strings.go [created: Sat, 18 Jan 2014]

package extension

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

var stringFuncs = Register(String{})

type charClass func(c rune) bool

func (cc charClass) Inverse() charClass {
	return func(c rune) bool { return !cc(c) }
}

type String struct{}

func (_s String) Namespace() string { return "Strings" }
func (_s String) UpperCamel(s string) (string, error) {
	class := charClass(unicode.IsLetter)
	s = strings.TrimFunc(s, class.Inverse())
	ss := strings.FieldsFunc(s, class.Inverse())

	if len(ss) == 1 && ss[0] == "" {
		return "", fmt.Errorf("no letters in string")
	}

	for i := range ss {
		c, n := utf8.DecodeRuneInString(ss[i])
		if c == utf8.RuneError && n == 1 {
			return "", fmt.Errorf("invalid rune %x found in string", c)
		}

		ctitle := strings.ToTitle(string([]rune{c}))

		ss[i] = ctitle + ss[i][n:]
	}

	return strings.Join(ss, ""), nil
}
