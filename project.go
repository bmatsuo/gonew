package main
/* 
*  File: project.go
*  Author: Bryan Matsuo [bmatsuo@soe.ucsc.edu] 
*  Created: Sat Jul  2 20:28:54 PDT 2011
*/
import (
    "os"
    "log"
    "exec"
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
	if DEBUG {
		log.Print("Creating Makefile")
		log.Printf("    template: %s", templatePath)
        if DEBUG_LEVEL > 0 {
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

func (p Project) MainTemplatePath() string {
    return filepath.Join(GetTemplateRoot(), "gofiles", "pkg.t")
}
func (p Project) CreateMainFile(dict map[string]string) os.Error {
    var mainfile = p.MainFilename()
    var templatePath = p.MainTemplatePath()
    var template = mustache.RenderFile(templatePath, dict, map[string]string{"file":mainfile})
    if DEBUG {
        log.Printf("Creating main file %s", mainfile)
        log.Printf("    template: %s", templatePath)
        if DEBUG_LEVEL > 0 {
            log.Print(dict)
            log.Print("\n", template, "\n")
        }
    }
    var templout = make([]byte, len(template))
    copy(templout, template)
    var errWrite = ioutil.WriteFile(mainfile, templout, FilePermissions)
    return errWrite
}

func (p Project) TestTemplatePath() string {
    return filepath.Join(GetTemplateRoot(), "gofiles", "test.t")
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
        template = mustache.RenderFile(templatePath, dict, map[string]string{"file":testfile})
    )
    if DEBUG {
        log.Printf("Creating main file %s", testfile)
        log.Printf("    template: %s", templatePath)
        if DEBUG_LEVEL > 0 {
            log.Print("\n", template, "\n")
        }
    }
    var templout = make([]byte, len(template))
    copy(templout, template)
    var errWrite = ioutil.WriteFile(testfile, templout, FilePermissions)
    return errWrite
}

func (p Project) ReadmeTemplatePath() string {
    var root = GetTemplateRoot()
    var useMarkdown = p.Host == GitHubHost
    if useMarkdown {
        return filepath.Join(root, "README", p.Type.String() + ".md.t")
    }
    return filepath.Join(root, "README", p.Type.String() + ".t")
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
    var template = mustache.RenderFile(templatePath, dict, map[string]string{"file":readme})
    if DEBUG {
        log.Print("Creating README")
        log.Printf("    template: %s", templatePath)
        if DEBUG_LEVEL > 0 {
            log.Print("\n", template, "\n")
        }
    }
    var templout = make([]byte, len(template))
    copy(templout, template)
    var errWrite = ioutil.WriteFile(readme, templout, FilePermissions)
    return errWrite
}

func (p Project) OtherTemplatePaths() []string {
    var root = GetTemplateRoot()
    var others = make([]string, 0, 1)
    switch p.Repo {
    case GitType:
        others = append(others, filepath.Join(root, "otherfiles", "gitignore.t"))
    }
    if len(others) == 0 {
        return nil
    }
    return others
}
func (p Project) CreateOtherFiles(dict map[string]string) os.Error {
    if DEBUG {
        log.Printf("Creating any other necessary files")
    }
    var templatePaths = p.OtherTemplatePaths()
    if templatePaths == nil {
        return nil
    }
    for _, path := range templatePaths {
        var template = mustache.RenderFile(path, dict)
        if DEBUG {
            log.Printf("    template: %s", path)
            if DEBUG_LEVEL > 0 {
                log.Print("\n", template, "\n")
            }
            var templout = make([]byte, len(template))
            copy(templout, template)
            var errWrite = ioutil.WriteFile(path, templout, FilePermissions)
            if errWrite != nil {
                return errWrite
            }
        }
    }
    return nil
}

func (p Project) InitializeRepo(commit bool) os.Error {
    switch p.Repo {
    case GitType:
        var (
            initcmd = exec.Command("git", "init")
            addcmd = exec.Command("git", "add", ".")
            commitcmd = exec.Command("git", "commit",
                    "-a", "-m", "Empty project generated by gonew.")
        )
        errInit := initcmd.Run()
        if errInit != nil {
            return errInit
        }
        errAdd := addcmd.Run()
        if errAdd != nil {
            return errAdd
        }
        if commit {
            errCommit := commitcmd.Run()
            if errCommit != nil {
                return errCommit
            }
        }
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

func (p Project) YearString() string {
    return time.LocalTime().Format("2006")
}

// fix the formatting of this method.
func (p Project) DateString() string {
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
    td["year"]   = p.YearString()
    td["date"]   = p.DateString()
    return td
}

func (p Project) Create() os.Error {
    var dict = p.GenerateDictionary()
    var errMkdir, errChdir, errRepo, errChdirBack os.Error
    var errMake, errMain, errTest, errReadme, errOther os.Error

    // Make the directory and change the working directory.
    errMkdir = os.Mkdir(p.Name, DirPermissions)
    if errMkdir != nil {
        return errMkdir
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
    errRepo = p.InitializeRepo(true)
    if errRepo != nil {
        return errRepo
    }

    // Change the working directory back.
    errChdirBack = os.Chdir("..")
    return errChdirBack
}
