// Copyright 2011, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

/*  Filename:    funcmap.go
 *  Author:      Bryan Matsuo <bryan.matsuo@gmail.com>
 *  Created:     Sun Oct 23 01:45:33 PDT 2011
 *  Description: 
 */

import (
	"template"
	"strings"
	"time"
	"fmt"
	"os"
)

//  The name supplied via ~/.gonewrc or by command line options. Accessible
//  as the template fuction "name"
func GonewUserName() string { return AppConfig.Name }

//  The email supplied via ~/.gonewrc or by command line options. Accessible
//  as the template function "email"
func GonewUserEmail() string { return AppConfig.Email }

//  The current year in a four-digit format. Accessible as the template
//  function "year"
func YearString() string {
	return time.LocalTime().Format("2006")
}

//  The local datetime in the default string format. Accessible as the
//  template function "date".
func DateString(format ...string) string {
	if len(format) > 0 {
		return time.LocalTime().Format(format[0])
	}
	return time.LocalTime().String()
}

//  A template helper function to import multiple libraries into a go file.
func GoImport(pkgs ...interface{}) (stmt string, err os.Error) {
	if len(pkgs) == 0 {
		stmt = "import ()"
		return
	}
	pkgstrs := make([]string, 0, len(pkgs))
	for _, s := range pkgs {
		switch s.(type) {
		case string:
			pkgstrs = append(pkgstrs, s.(string))
		case []string:
			pkgstrs = append(pkgstrs, s.([]string)...)
		case []interface{}:
			for _, s1 := range s.([]interface{}) {
				switch s1.(type) {
				case string:
					pkgstrs = append(pkgstrs, s1.(string))
				default:
					err = os.NewError("import argument slice element is not a string")
					return
				}
			}
		default:
			err = os.NewError("import argument not a (slice of) string(s)")
			return
		}
	}

	for i, pkg := range pkgstrs {
		pkgfmt := `    "%s"`
		if strings.HasPrefix(pkg, "//") {
			pkg = pkg[2:]
			pkgfmt = `    //"%s"`
		}
		pkgstrs[i] = fmt.Sprintf(pkgfmt, pkg)
	}
	pieces := make([]string, len(pkgstrs)+2)
	pieces[0] = "import ("
	copy(pieces[1:], pkgstrs)
	pieces[len(pieces)-1] = ")"
	return strings.Join(pieces, "\n"), nil
}

//  A template helper function to define an empty function with specified
//  arguments and no return values.
func GoFunction(name string, args ...string) string {
	return fmt.Sprintf("func %s(%s) {\n\n}", name, strings.Join(args, ", "))
}

func GoNiladic(name string, stmts ...string) string {
	return fmt.Sprintf("func %s() {\n    %s\n}", name, strings.Join(stmts, "\n    "))
}

// A template helper function to define a command's 'main' function with a list
// of statements.
func GoMain(stmts ...string) string { return GoNiladic("main", stmts...) }

// A template helper function to define a package/command's 'init' function with a list
// of statements.
func GoInit(stmts ...string) string { return GoNiladic("init", stmts...) }

func DefaultFuncMap() template.FuncMap {
	return template.FuncMap{
		"import": GoImport,
		"init":   GoInit,
		"main":   GoMain,
		"func":   GoFunction,
		"date":   DateString,
		"year":   YearString,
		"name":   GonewUserName,
		"email":  GonewUserEmail,
	}
}
