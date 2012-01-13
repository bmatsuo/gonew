// Copyright 2011, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main
/* 
*  File: project.go
*  Author: Bryan Matsuo [bmatsuo@soe.ucsc.edu] 
*  Created: Sat Jul  2 20:28:54 PDT 2011
 */
import (
	"errors"
	"os"
	"fmt"
	//"log"
)

var (
	NoUserNameError        = errors.New("Missing remote repository username.")
	NoRemoteError          = errors.New("Missing remote repository url")
	DirPermissions  uint32 = 0755
	FilePermissions uint32 = 0644
)

type ProjectType int

const (
	NilProjectType ProjectType = iota
	CmdType
	PkgType
	//LibType
	//MainType
)

func (ptype ProjectType) String() string {
	switch ptype {
	case CmdType:
		return "cmd"
	case PkgType:
		return "pkg"
	}
	return ""
}

func DefaultTarget(pname string) string {
	// TODO strip special characters from the prjoect name.
	return pname
}

type Project struct {
	Name       string
	Target     string
	ImportLibs []string
	User       string
	Remote     string
	License    LicenseType
	Type       ProjectType
	Host       RepoHost
	Repo       RepoType
	Markdown   bool
}

func (p Project) Create() error {
	//var dict = p.GenerateDictionary()

	// Make the directory and change the working directory.
	if DEBUG || VERBOSE {
		fmt.Print("Creating project directory.\n")
	}

	var err error
	if err = os.Mkdir(p.Name, DirPermissions); err != nil {
		return err
	}

	if DEBUG || VERBOSE {
		fmt.Print("Entering project directory.\n")
	}

	if err = os.Chdir(p.Name); err != nil {
		return err
	}
	if err = p.CreateFiles(); err != nil {
		return err
	}
	if userepo {
		if err = p.InitializeRepo(true, true, true); err != nil {
			return err
		}
	}

	if DEBUG || VERBOSE {
		fmt.Print("Leaving project directory.\n")
	}

	return os.Chdir("..")
}
func (p Project) CreateFiles() error {
	var err error
	for _, f := range p.Files() {
		if err = f.Create(); err != nil {
			return err
		}
	}
	return nil
}

func (p Project) Files() []ProjectFile {
	ps := make([]ProjectFile, 0, 7)
	ps = append(ps, p.Makefile())
	ps = append(ps, p.MainFile())
	ps = append(ps, p.TestFile())
	ps = append(ps, p.ReadmeFile())
	if p.License != NilLicenseType {
		ps = append(ps, p.LicenseFile())
	}
	if p.IsCommand() {
		ps = append(ps, p.OptionsFile())
		ps = append(ps, p.DocFile())
	}
	ps = append(ps, p.OtherFiles()...)
	return ps
}

func (p Project) MainFilename() string {
	return p.Target + ".go"
}
func (p Project) TestFilename() string {
	switch p.Type {
	case CmdType:
		return "main_test.go"
	case PkgType:
		return p.Target + "_test.go"
	}
	return p.Target + "_test.go"
}
func (p Project) ReadmeFilename() string {
	if p.ReadmeIsMarkdown() {
		return "README.md"
	}
	return "README"
}

func (p Project) Makefile() ProjectFile {
	return ProjectFile{
		Name:      "Makefile",
		Desc:      fmt.Sprintf("Makefile for %s", p.Name),
		Type:      MakeFile,
		Template:  "Makefile.t",
		DebugDesc: "makefile",
		p:         p,
	}
}
func (p Project) MainFile() ProjectFile {
	return ProjectFile{
		Name:      p.MainFilename(),
		Desc:      fmt.Sprintf("Main source file in %s", p.Name),
		Type:      GoFile,
		Template:  fmt.Sprintf("go.%s.t", p.Type.String()),
		DebugDesc: "main file",
		p:         p,
	}
}
func (p Project) TestFile() ProjectFile {
	return ProjectFile{
		Name:      p.TestFilename(),
		Desc:      fmt.Sprintf("Main test file for %s", p.Name),
		Type:      GoFile,
		Template:  "test.t",
		DebugDesc: "main test",
		p:         p,
	}
}
func (p Project) ReadmeFile() ProjectFile {
	return ProjectFile{
		Name:      p.ReadmeFilename(),
		Desc:      fmt.Sprintf("%s is the best program for...", p.Name),
		Type:      ReadmeFile,
		Template:  "README.t",
		DebugDesc: "README",
		p:         p,
	}
}

func (p Project) OptionsFile() ProjectFile {
	return ProjectFile{
		Name:      "options.go",
		Desc:      fmt.Sprintf("Option parsing for %s", p.Name),
		Type:      GoFile,
		Template:  "go.options.t",
		DebugDesc: "option parse file",
		p:         p,
	}
}

func (p Project) DocFile() ProjectFile {
	return ProjectFile{
		Name:      "doc.go",
		Desc:      fmt.Sprintf("Godoc documentation for %s", p.Name),
		Type:      GoFile,
		Template:  "go.doc.t",
		DebugDesc: "documentation",
		p:         p,
	}
}

func (p Project) LicenseFile() ProjectFile {
	return ProjectFile{
		Name:      "LICENSE",
		Desc:      fmt.Sprintf("License for using %s", p.Name),
		Type:      OtherFile,
		Template:  p.License.FullTemplateName(),
		DebugDesc: "license file",
		p:         p,
	}
}

func (p Project) OtherFiles() []ProjectFile {
	var others = make([]ProjectFile, 0, 1)
	switch p.Repo {
	case GitType:
		others = append(others, ProjectFile{
			Name:      ".gitignore",
			Desc:      "Junk files to ignore",
			Type:      OtherFile,
			Template:  "other.gitignore.t",
			DebugDesc: "ignore file",
			p:         p,
		})
	case HgType:
		others = append(others, ProjectFile{
			Name:      ".hgignore",
			Desc:      "Junk files to ignore",
			Type:      OtherFile,
			Template:  "other.hgignore.t",
			DebugDesc: "ignore file",
			p:         p,
		})
	}
	if len(others) == 0 {
		return nil
	}
	return others
}

//  Returns the remote repo address. For instance, "github.com/ghuser/go-project"
//  if p.Host is GitHubHost. Returns a placeholderstring if p.Host is not defined.
func (p Project) HostRepoString() string {
	switch p.Host {
	case GitHubHost:
		if p.User == "" {
			return ""
		}
		return "github.com/" + p.User + "/" + p.Target
	}
	return "<INSERT REPO HOST HERE>"
}

//  Initialize the repository after templates have been successfully
//  generated.
func (p Project) InitializeRepo(add, commit, push bool) error {
	switch p.Repo {
	case GitType:
		git := GitRepo{}
		git.Initialize(add, commit)
	case HgType:
		hg := HgRepo{}
		hg.Initialize(add, commit)
	}
	switch p.Host {
	case GitHubHost:
		origin := p.Remote
		if origin == "" {
			return nil
		}
		github := GitHubRepo{p}
		if err := github.Init(origin); err != nil {
			return err
		}
		if push {
			if err := github.Push(); err != nil {
				return err
			}
		}
		return nil
	}
	return nil
}

//  Returns true if the repo host uses Markdown enabled README files.
func (p Project) ReadmeIsMarkdown() bool {
	if userepo {
		return p.Host == GitHubHost
	}
	return p.Markdown
}
//  The target is a command.
func (p Project) IsCommand() bool { return p.Type == CmdType }
//  The target is an importable package
func (p Project) IsPackage() bool { return p.Type == PkgType }
