// Copyright 2011, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

/*
 *  Filename:    file.go
 *  Package:     main
 *  Author:      Bryan Matsuo <bmatsuo@soe.ucsc.edu>
 *  Created:     Sun Jul  3 16:57:42 PDT 2011
 *  Description: 
 */
import (
	"fmt"
	"strings"
	"unicode"
)

func TestName(filename string) string {
	var test = filename
	if strings.HasSuffix(test, ".go") {
		test = test[:len(test)-3]
	}
	if strings.HasSuffix(test, "_test") {
		test = test[:len(test)-5]
	}
	return strings.Title(
		strings.Map(
			func(c rune) rune {
				if unicode.IsNumber(c) || unicode.IsLetter(c) {
					return c
				}
				return -1
			},
			test))
}

type FileType uint

const (
	GoFile FileType = iota
	ReadmeFile
	MakeFile
	LicenseFile
	OtherFile
)

//  ProjectFile satisfies the Context interface.
type ProjectFile struct {
	Name      string
	Desc      string
	Type      FileType
	Template  string
	DebugDesc string
	p         Project
}

func (f ProjectFile) Filename() string { return f.Name }

func (f ProjectFile) Package() string {
	if f.p.Type == PkgType {
		return f.p.Target
	}
	return "main"
}

func (f ProjectFile) Imports() []string        { return f.p.ImportLibs }
func (f ProjectFile) Description() string      { return f.Desc }
func (f ProjectFile) DebugDescription() string { return f.DebugDesc }
func (f ProjectFile) Tests() []string          { return []string{TestName(f.p.Target)} }
func (f ProjectFile) Project() Project         { return f.p }
func (f ProjectFile) LicenseType() LicenseType { return f.p.License }
func (f ProjectFile) FileType() FileType       { return f.Type }
func (f ProjectFile) Create() (err error)      { return CreateFile(f, f.Template) }

type TestFile struct {
	File
	Name    string
	Pkg     string
	License LicenseType
	Repo    RepoType
	Host    RepoHost
}

func (t TestFile) Filename() string         { return t.Name + "_test.go" }
func (t TestFile) Package() string          { return t.Pkg }
func (t TestFile) Imports() []string        { return nil }
func (t TestFile) Description() string      { return fmt.Sprintf("For testing %s", t.File.Filename()) }
func (t TestFile) DebugDescription() string { return "test for " + t.File.Filename() }
func (t TestFile) FileType() FileType       { return GoFile }
func (t TestFile) Tests() []string          { return t.File.Tests() }
func (t TestFile) Project() Project         { return Project{} }
func (t TestFile) LicenseType() LicenseType { return t.License }
func (t TestFile) Create() error            { return CreateFile(t, "test.t") }

func (t TestFile) TemplatePath() []string { return []string{"testfiles", "pkg.t"} }

type File struct {
	Name       string
	User       string
	Pkg        string
	ImportLibs []string
	Desc       string
	License    LicenseType
	Repo       RepoType
	Host       RepoHost
}

func (f File) Filename() string         { return f.Name + ".go" }
func (f File) Package() string          { return f.Pkg }
func (f File) Imports() []string        { return f.ImportLibs }
func (f File) Description() string      { return f.Desc }
func (f File) DebugDescription() string { return "library" }
func (f File) FileType() FileType       { return GoFile }
func (f File) Tests() []string          { return []string{TestName(f.Name)} }
func (f File) Project() Project         { return Project{} }
func (f File) LicenseType() LicenseType { return f.License }
func (f File) Create() error            { return CreateFile(f, "go.lib.t") }

func (f File) TemplatePath() []string { return []string{"gofiles", "lib.t"} }
func (f File) TestFile() TestFile {
	return TestFile{File: f, Name: f.Name, Pkg: f.Pkg, Repo: f.Repo, Host: f.Host}
}
