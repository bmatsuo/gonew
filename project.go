package main
/* 
*  File: project.go
*  Author: Bryan Matsuo [bmatsuo@soe.ucsc.edu] 
*  Created: Sat Jul  2 20:28:54 PDT 2011
*/
import (
    "os"
    "log"
    "fmt"
    //"exec"
    "path/filepath"
    "io/ioutil"
    "time"
    "github.com/hoisie/mustache.go"
)

var (
    DirPermissions  uint32 = 0755
    FilePermissions uint32 = 0644
)

type RepoType int
const(
    NilRepoType RepoType = iota
    GitType
    // MercurialType
    // ...
)

type RepoHost int
const (
    NilRepoHost RepoHost = iota
    GitHubHost
    //GoogleHost
    //...
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
    Type   ProjectType
    Host   RepoHost
    Repo   RepoType
}

func GetGoroot() string {
    goroot, err := os.Getenverror("GOROOT")
    if err != nil {
        panic("goroot")
    }
    return goroot
}
func GetTemplateRoot() string {
    var goroot = GetGoroot()
    return filepath.Join(goroot, "src", "pkg",
            "github.com", "bmatsuo", "gonew", "templates")
}

func (p Project) MakefileTemplatePath() string {
    return filepath.Join(GetTemplateRoot(), "makefiles", p.Type.String() + ".t")
}
func (p Project) CreateMakefile(dict map[string]string) os.Error {
    var (
        templatePath = p.MakefileTemplatePath()
        template = mustache.RenderFile(templatePath, dict, map[string]string{"file":"Makefile"})
    )
	if DEBUG || VERBOSE {
		fmt.Print("Creating Makefile\n")
    }
    if DEBUG && DEBUG_LEVEL > 0 {
        log.Printf("template: %s", templatePath)
        if DEBUG_LEVEL > 1 {
	     log.Print("\n", template, "\n")
        }
    }
    var templout = make([]byte, len(template))
    copy(templout, template)
    var errWrite = ioutil.WriteFile("Makefile", templout, FilePermissions)
    return errWrite
}

func (p Project) MainFilename() string {
    return p.Name + ".go"
    /*
    switch p.Type {
    case CmdType:
    case PkgType:
    }
    */
}

func (p Project) MainTemplatePath() []string {
    switch p.Type {
    case CmdType:
        return []string{GetTemplateRoot(), "gofiles", "cmd.t"}
    case PkgType:
        return []string{GetTemplateRoot(), "gofiles", "pkg.t"}
    }
    return []string{""}
}
func (p Project) CreateMainFile(dict map[string]string) os.Error {
    var (
        mainfile = p.MainFilename()
        templatePath = p.MainTemplatePath()
        errWrite = WriteTemplate(mainfile, "main file", dict, templatePath...)
    )
    return errWrite
}

func (p Project) TestTemplatePath() []string {
    if p.Type == CmdType {
        return []string{GetTemplateRoot(), "testfiles", "cmd.t"}
    }
    return []string{GetTemplateRoot(), "testfiles", "pkg.t"}
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

func (p Project) CreateTestFile(dict map[string]string) os.Error {
    var (
        testfile = p.TestFilename()
        templatePath = p.TestTemplatePath()
		errWrite = WriteTemplate(testfile, "test file", dict, templatePath...)
    )
    return errWrite
}

func (p Project) ReadmeTemplatePath() []string {
    var root = GetTemplateRoot()
    var useMarkdown = p.Host == GitHubHost
    if useMarkdown {
        return []string{root, "README", p.Type.String() + ".md.t"}
    }
    return []string{root, "README", p.Type.String() + ".t"}
}
func (p Project) CreateReadme(dict map[string]string) os.Error {
    var (
        templatePath = p.ReadmeTemplatePath()
        readme = "README"
        useMarkdown = p.Host == GitHubHost
    )
    if useMarkdown {
        readme += ".md"
    }
	var errWrite = WriteTemplate(readme, "README", dict, templatePath...)
    return errWrite
}

func (p Project) DocTemplatePath() []string {
    var root = GetTemplateRoot()
    return []string{root, "gofiles", "doc.t"}
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

func (p Project) OtherTemplatePaths() [][]string {
    var root = GetTemplateRoot()
    var others = make([][]string, 0, 1)
    switch p.Repo {
    case GitType:
        others = append(others, []string{".gitignore", root, "otherfiles", "gitignore.t"})
    }
    if len(others) == 0 {
        return nil
    }
    return others
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

func (p Project) InitializeRepo(add, commit bool) os.Error {
    switch p.Repo {
    case GitType:
		var git = GitRepository{}
		git.Initialize(add, commit)
    }
    return nil
}

// fix this method.
func (p Project) HostString() string {
    switch p.Repo {
    case GitType:
        return "github.com/" + AppConfig.HostUser
    }
    return "<INSERT REPO HOST HERE>"
}

func YearString() string {
    return time.LocalTime().Format("2006")
}

// fix the formatting of this method.
func DateString() string {
    return time.LocalTime().String()
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

func (p Project) Create() os.Error {
    var dict = p.GenerateDictionary()
    var errMkdir, errChdir, errRepo, errChdirBack os.Error
    var errMake, errMain, errDoc, errTest, errReadme, errOther os.Error

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

    // Create the project files.
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
    errRepo = p.InitializeRepo(true, true)
    if errRepo != nil {
        return errRepo
    }

    // Change the working directory back.
    if DEBUG || VERBOSE {
        fmt.Print("Leaving project directory.\n")
    }
    errChdirBack = os.Chdir("..")
    return errChdirBack
}
