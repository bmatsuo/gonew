package main
/* 
*  File: project.go
*  Author: Bryan Matsuo [bmatsuo@soe.ucsc.edu] 
*  Created: Sat Jul  2 20:28:54 PDT 2011
*/
import (
    "os"
    "fmt"
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
    Name   string
    Target string
    User   string
    Remote string
    Type   ProjectType
    Host   RepoHost
    Repo   RepoType
}

func (p Project) Create() os.Error {
    var dict = p.GenerateDictionary()
    var errMkdir, errChdir, errFiles, errRepo, errChdirBack os.Error

    // Make the directory and change the working directory.
    if DEBUG || VERBOSE {
        fmt.Print("Creating project directory.\n")
    }
    errMkdir = os.Mkdir(p.Name, DirPermissions)
    if errMkdir != nil {
        return errMkdir
    }
    if DEBUG || VERBOSE {
        fmt.Print("Entering project directory.\n")
    }
    errChdir = os.Chdir(p.Name)
    if errChdir != nil {
        return errChdir
    }
    errFiles = p.CreateFiles(dict)
    if errFiles != nil {
        return errFiles
    }
    errRepo = p.InitializeRepo(true, true, true)
    if errRepo != nil {
        return errRepo
    }
    if DEBUG || VERBOSE {
        fmt.Print("Leaving project directory.\n")
    }
    errChdirBack = os.Chdir("..")
    return errChdirBack
}
func (p Project) CreateFiles(dict map[string]string) os.Error {
    var errMake, errMain, errDoc, errTest, errReadme, errOther os.Error
    errMake = p.CreateMakefile(dict)
    if errMake != nil {
        return errMake
    }
    errMain = p.CreateMainFile(dict)
    if errMain != nil {
        return errMain
    }
    errDoc = p.CreateDocFile(dict)
    if errDoc != nil {
        return errDoc
    }
    errTest = p.CreateTestFile(dict)
    if errTest != nil {
        return errTest
    }
    errReadme = p.CreateReadme(dict)
    if errReadme != nil {
        return errReadme
    }
    errOther = p.CreateOtherFiles(dict)
    if errOther != nil {
        return errOther
    }
    return nil
}

func (p Project) GenerateDictionary() map[string]string {
    var td = make(map[string]string, 9)
    td["project"]   = p.Name
    td["name"]   = AppConfig.Name
    td["email"]  = AppConfig.Email
    td["gotarget"] = p.Target
    td["main"]   = p.MainFilename()
    td["type"]   = p.Type.String()
    td["repo"]   = p.HostString()
    td["year"]   = YearString()
    td["date"]   = DateString()
    return td
}

func (p Project) CreateMakefile(dict map[string]string) os.Error {
    var (
        templatePath = p.MakefileTemplatePath()
        errWrite = WriteTemplate("Makefile", "makefile", dict, templatePath...)
    )
    return errWrite
}
func (p Project) CreateMainFile(dict map[string]string) os.Error {
    var (
        mainfile = p.MainFilename()
        templatePath = p.MainTemplatePath()
        errWrite = WriteTemplate(mainfile, "main file", dict, templatePath...)
    )
    return errWrite
}
func (p Project) CreateTestFile(dict map[string]string) os.Error {
    var (
        testfile = p.TestFilename()
        templatePath = p.TestTemplatePath()
		errWrite = WriteTemplate(testfile, "test file", dict, templatePath...)
    )
    return errWrite
}
func (p Project) CreateReadme(dict map[string]string) os.Error {
    var (
        templatePath = p.ReadmeTemplatePath()
        readme = p.ReadmeFilename()
    )
	var errWrite = WriteTemplate(readme, "README", dict, templatePath...)
    return errWrite
}
func (p Project) CreateDocFile(dict map[string]string) os.Error {
    if p.Type == PkgType {
        return nil
    }
	var (
		doc = "doc.go"
        templatePath = p.DocTemplatePath()
		errWrite = WriteTemplate(doc, "documentation files", dict, templatePath...)
	)
    return errWrite
}
func (p Project) CreateOtherFiles(dict map[string]string) os.Error {
    var templatePaths = p.OtherTemplatePaths()
    if templatePaths == nil {
        return nil
    }
    for _, path := range templatePaths {
        var errWrite = WriteTemplate(path[0], "other file", dict, path[1:]...)
        if errWrite != nil {
            return errWrite
        }
    }
    return nil
}
func (p Project) InitializeRepo(add, commit, push bool) os.Error {
    switch p.Repo {
    case GitType:
		var git = GitRepository{}
		git.Initialize(add, commit)
    }
    switch p.Host {
    case GitHubHost:
        var origin = p.Remote
        if origin == "" {
            return nil
        }
        var github = GitHubRepository{p}
        var errGHInit = github.Init(origin)
        if errGHInit != nil {
            return errGHInit
        }
        if push {
            var errPush = github.Push()
            if errPush != nil {
                return errPush
            }
        }
        return nil
    }
    return nil
}

func (p Project) MainFilename() string {
    return p.Name + ".go"
}
func (p Project) TestFilename() string {
    switch p.Type {
    case CmdType:
        return "main_test.go"
    case PkgType:
        return p.Name + "_test.go"
    }
    return p.Name + "_test.go"
}
func (p Project) ReadmeFilename() string {
    if p.ReadmeIsMarkdown() {
        return "README.md"
    }
    return "README"
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
func (p Project) DocTemplatePath() []string {
    return []string{"gofiles", "doc.t"}
}
func (p Project) OtherTemplatePaths() [][]string {
    var others = make([][]string, 0, 1)
    switch p.Repo {
    case GitType:
        others = append(others, []string{".gitignore", "otherfiles", "gitignore.t"})
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
func YearString() string {
    return time.LocalTime().Format("2006")
}
func DateString() string {
    return time.LocalTime().String()
}
func (p Project) ReadmeIsMarkdown() bool {
    return p.Host == GitHubHost
}
