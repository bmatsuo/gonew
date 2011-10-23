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
    "os"
)

type FileType uint

const (
    Go FileType = iota
    README
    Makefile
    License
    Other
)

//  ProjectFile satisfies the Context interface.
type ProjectFile struct {
    Name string
    Desc string
    Type FileType
    Template string
    DebugDesc string
    p Project
}

func (f ProjectFile) Filename() string { return f.Name }

func (f ProjectFile) Package() string {
    if f.p.Type == PkgType {
        return f.p.Target
    }
    return "main"
}

func (f ProjectFile) Description() string { return f.Desc }
func (f ProjectFile) DebugDescription() string { return f.DebugDesc }
func (f ProjectFile) Tests() []string {
    return []string{ f.p.Target }
}
func (f ProjectFile) Project() Project { return f.p }
func (f ProjectFile) LicenseType() LicenseType { return f.p.License }
func (f ProjectFile) FileType() FileType { return f.Type }

func (f ProjectFile) Create(ms TemplateMultiSet) (err os.Error) {
    return CreateFile(f, f.Template, ms)
}

type TestFile struct {
    File
    Name    string
    Pkg     string
    License LicenseType
    Repo    RepoType
    Host    RepoHost
}

func (t TestFile) Filename() string { return t.Name + "_test.go" }
func (t TestFile) Package() string { return t.Pkg }
func (t TestFile) Description() string {
    return fmt.Sprintf("For testing %s", t.File.Filename())
}
func (t TestFile) DebugDescription() string {
    return t.File.Filename() + " test"
}
func (t TestFile) FileType() FileType { return Go }
func (t TestFile) Tests() []string { return t.File.Tests() }
func (t TestFile) Project() Project { return Project{} }
func (t TestFile) LicenseType() LicenseType { return t.License }
func (t TestFile) Create(ms TemplateMultiSet) os.Error {
    return CreateFile(t, "test.t", ms)
}

func (t TestFile) TemplatePath() []string { return []string{"testfiles", "pkg.t"} }

type File struct {
    Name    string
    User    string
    Pkg     string
    Desc    string
	License LicenseType
    Repo    RepoType
    Host    RepoHost
}

func (f File) Filename() string { return f.Name + ".go" }
func (f File) Package() string { return f.Pkg }
func (f File) Description() string { return f.Desc }
func (f File) DebugDescription() string { return "library" }
func (f File) FileType() FileType { return Go }
func (f File) Tests() []string { return []string{f.Name} }
func (f File) Project() Project { return Project{} }
func (f File) LicenseType() LicenseType { return f.License }
func (f File) Create(ms TemplateMultiSet) os.Error {
    return CreateFile(f, "go.lib.t", ms)
}

func (f File) TemplatePath() []string { return []string{"gofiles", "lib.t"} }
func (f File) TestFile() TestFile {
    return TestFile{File:f, Name:f.Name, Pkg:f.Pkg, Repo:f.Repo, Host:f.Host}
}
