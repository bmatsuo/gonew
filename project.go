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
    "os"
    "fmt"
    "log"
    "time"
)

var (
    NoUserNameError        = os.NewError("Missing remote repository username.")
    NoRemoteError          = os.NewError("Missing remote repository url")
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
    Name    string
    Target  string
    User    string
    Remote  string
    License LicenseType
    Type    ProjectType
    Host    RepoHost
    Repo    RepoType
}

func (p Project) Create() os.Error {
    var dict = p.GenerateDictionary()

    // Make the directory and change the working directory.
    if DEBUG || VERBOSE {
        fmt.Print("Creating project directory.\n")
    }

    var err os.Error
    if err = os.Mkdir(p.Name, DirPermissions); err != nil {
        return err
    }

    if DEBUG || VERBOSE {
        fmt.Print("Entering project directory.\n")
    }

    if err = os.Chdir(p.Name); err != nil {
        return err
    }
    if err = p.CreateFiles(dict); err != nil {
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
func (p Project) CreateFiles(dict map[string]string) os.Error {
    var err os.Error
    if err = p.CreateMakefile(dict); err != nil {
        return err
    }
    if err = p.CreateMainFile(dict); err != nil {
        return err
    }
    if err = p.CreateOptionsFile(dict); err != nil {
        return err
    }
    if err = p.CreateDocFile(dict); err != nil {
        return err
    }
    if err = p.CreateTestFile(dict); err != nil {
        return err
    }
    if err = p.CreateReadme(dict); err != nil {
        return err
    }
    if err = p.CreateOtherFiles(dict); err != nil {
        return err
    }
    if err = p.CreateLicense(dict); err != nil {
        return err
    }
    return nil
}

func (p Project) GenerateDictionary() map[string]string {
    var td = make(map[string]string, 9)
    td["project"] = p.Name
    td["name"] = AppConfig.Name
    td["email"] = AppConfig.Email
    td["gotarget"] = p.Target
    td["main"] = p.MainFilename()
    td["type"] = p.Type.String()
    td["repo"] = p.HostRepoString()
    td["year"] = YearString()
    td["date"] = DateString()
    return td
}

func (p Project) CreateMakefile(dict map[string]string) os.Error {
    var (
        templatePath = p.MakefileTemplatePath()
        errWrite     = WriteTemplate("Makefile", "makefile", dict, templatePath...)
    )
    return errWrite
}
func (p Project) CreateMainFile(dict map[string]string) os.Error {
    var (
        mainfile      = p.MainFilename()
        ltemplateName = p.GofileLicenseHeadTemplateName()
        templatePath  = p.MainTemplatePath()
    )
    if ltemplateName == "" {
        if DEBUG {
            log.Print("No license template found for %s", p.License.String())
        }
        return WriteTemplate(mainfile, "main file", dict, templatePath...)
    }
    var ltemplatePath = []string{"licenses", ltemplateName}
    errWrite := WriteTemplate(mainfile, "main file license", dict, ltemplatePath...)
    if errWrite != nil {
        return errWrite
    }
    return AppendTemplate(mainfile, "main file contents", dict, templatePath...)
}
func (p Project) CreateTestFile(dict map[string]string) os.Error {
    if !AppConfig.MakeTest {
        if DEBUG || VERBOSE {
            fmt.Printf("Skipping test file generation.")
        }
        return nil
    }
    var (
        testfile      = p.TestFilename()
        ltemplateName = p.GofileLicenseHeadTemplateName()
        templatePath  = p.TestTemplatePath()
    )
    if ltemplateName == "" {
        if DEBUG {
            log.Print("No license template found for %s", p.License.String())
        }
        return WriteTemplate(testfile, "test file", dict, templatePath...)
    }
    var ltemplatePath = []string{"licenses", ltemplateName}
    errWrite := WriteTemplate(testfile, "test file license", dict, ltemplatePath...)
    if errWrite != nil {
        return errWrite
    }
    return AppendTemplate(testfile, "test file contents", dict, templatePath...)
}
func (p Project) CreateReadme(dict map[string]string) os.Error {
    var (
        templatePath = p.ReadmeTemplatePath()
        readme       = p.ReadmeFilename()
    )
    errWrite := WriteTemplate(readme, "README", dict, templatePath...)
    if errWrite != nil {
        return errWrite
    }
    if p.License == NilLicenseType {
        return nil
    }
    ltemplatePath := []string{"licenses", p.ReadmeLicenseTailTemplateName()}
    return AppendTemplate(readme, "README license tail", dict, ltemplatePath...)
}
func (p Project) CreateLicense(dict map[string]string) os.Error {
    var (
        templatePath = []string{"licenses", p.LicenseTemplateName()}
        license      = "LICENSE"
    )
    return WriteTemplate(license, "license", dict, templatePath...)
}
func (p Project) CreateOptionsFile(dict map[string]string) os.Error {
    if p.Type == PkgType {
        return nil
    }
    var (
        doc          = "options.go"
        templatePath = p.OptionsTemplatePath()
    )
    return WriteTemplate(doc, "option parsing file", dict, templatePath...)
}
func (p Project) CreateDocFile(dict map[string]string) os.Error {
    if p.Type == PkgType {
        return nil
    }
    var (
        doc          = "doc.go"
        templatePath = p.DocTemplatePath()
    )
    return WriteTemplate(doc, "documentation file", dict, templatePath...)
}
func (p Project) CreateOtherFiles(dict map[string]string) os.Error {
    templatePaths := p.OtherTemplatePaths()
    if templatePaths == nil {
        return nil
    }
    for _, path := range templatePaths {
        errWrite := WriteTemplate(path[0], "other file", dict, path[1:]...)
        if errWrite != nil {
            return errWrite
        }
    }
    return nil
}
func (p Project) InitializeRepo(add, commit, push bool) os.Error {
    switch p.Repo {
    case GitType:
        git := GitRepository{}
        git.Initialize(add, commit)
    case MercurialType:
        hg := MercurialRepository{}
        hg.Initialize(add, commit)
    }
    switch p.Host {
    case GitHubHost:
        origin := p.Remote
        if origin == "" {
            return nil
        }
        github := GitHubRepository{p}
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

func (p Project) LicenseTemplateName() string {
    lstring := p.License.TemplateNamePrefix()
    if lstring == "" {
        return ""
    }
    return lstring + ".t"
}
func (p Project) GofileLicenseHeadTemplateName() string {
    lstring := p.License.TemplateNamePrefix()
    if lstring == "" {
        return ""
    }
    return lstring + ".gohead.t"
}
func (p Project) ReadmeLicenseTailTemplateName() string {
    lstring := p.License.TemplateNamePrefix()
    if lstring == "" {
        return ""
    }
    if p.ReadmeIsMarkdown() {
        return lstring + ".readme.md.t"
    }
    return lstring + ".readme.t"
}
func (p Project) MainTemplateName() string {
    switch p.Type {
    case PkgType:
        return "pkg.t"
    case CmdType:
        return "cmd.t"
    }
    return ""
}
func (p Project) TestTemplateName() string {
    switch p.Type {
    case PkgType:
        return "pkg.t"
    case CmdType:
        return "cmd.t"
    }
    return ""
}
func (p Project) ReadmeTemplateName() string {
    switch p.Type {
    case PkgType:
        if p.ReadmeIsMarkdown() {
            return "pkg.md.t"
        } else {
            return "pkg.t"
        }
    case CmdType:
        if p.ReadmeIsMarkdown() {
            return "cmd.md.t"
        } else {
            return "cmd.t"
        }
    }
    return ""
}

func (p Project) MakefileTemplatePath() []string {
    return []string{"makefiles", p.Type.String() + ".t"}
}
func (p Project) MainTemplatePath() []string {
    return []string{"gofiles", p.MainTemplateName()}
}
func (p Project) TestTemplatePath() []string {
    return []string{"testfiles", p.TestTemplateName()}
}
func (p Project) ReadmeTemplatePath() []string {
    return []string{"README", p.ReadmeTemplateName()}
}
func (p Project) OptionsTemplatePath() []string {
    return []string{"gofiles", "options.t"}
}
func (p Project) DocTemplatePath() []string {
    return []string{"gofiles", "doc.t"}
}
func (p Project) OtherTemplatePaths() [][]string {
    var others = make([][]string, 0, 1)
    switch p.Repo {
    case GitType:
        others = append(others, []string{".gitignore", "otherfiles", "gitignore.t"})
    case MercurialType:
        others = append(others, []string{".hgignore", "otherfiles", "hgignore.t"})
    }
    if len(others) == 0 {
        return nil
    }
    return others
}

func (p Project) HostString() string {
    switch p.Host {
    case GitHubHost:
        if p.User == "" {
            return ""
        }
        return "github.com/" + p.User
    }
    return "<INSERT REPO HOST HERE>"
}
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
func YearString() string {
    return time.LocalTime().Format("2006")
}
func DateString() string {
    return time.LocalTime().String()
}
func (p Project) ReadmeIsMarkdown() bool {
    return userepo && p.Host == GitHubHost
}
