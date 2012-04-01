// Copyright 2011, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

/* 
*  File: gonew_main.go
*  Author: Bryan Matsuo [bmatsuo@soe.ucsc.edu] 
*  Created: Sat Jul  2 19:17:53 PDT 2011
*  Usage: gonew [options]
 */
import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"unicode"
)

var (
	GonewRoot    string            // The Gonew source directory
	TemplateRoot string            // The "templates" subdirectory of GonewRoot.
	Templates    TemplateHierarchy // The (sets of) templates used by Gonew.
)

func ArgumentError(msg string) {
	fmt.Fprintf(os.Stderr, "%s\n", msg)
	fs := setupFlags()
	fs.Usage()
}

var (
	usage          = `
gonew [options] cmd NAME
gonew [options] pkg NAME
gonew [options] lib NAME PKG
`
	printUsageHead = func() { fmt.Fprint(os.Stderr, usage, "\n") }
	userepo        = true
	norepo         = false
	usehost        = true
	nohost         = false
	VERBOSE        = false
	DEBUG          = false
	DEBUG_LEVEL    = -1
	name           string
	importlist     string
	ptype          string
	repo           string
	host           string
	user           string
	target         string
	license        string
	remote         string
	help           bool
)

func Debug(level int, msg string) {
	if DEBUG && DEBUG_LEVEL >= level {
		log.Print(msg)
	}
}

func Verbose(msg string) {
	if DEBUG || VERBOSE {
		fmt.Print(msg)
	}
}

func setupFlags() *flag.FlagSet {
	var fs = flag.NewFlagSet("gonew", flag.ExitOnError)
	fs.StringVar(&importlist,
		"import", "", "Packages to import in source .go files.")
	fs.StringVar(&repo,
		"repo", "", "Repository type (e.g. 'git').")
	fs.StringVar(&host,
		"host", "", "Repository host if any (e.g. 'github').")
	fs.StringVar(&user,
		"user", "", "Repo host username.")
	fs.StringVar(&remote,
		"remote", "", "Remote repository url to initialize and push to.")
	fs.StringVar(&target,
		"target", "", "Project name (executable/package name). Default based on NAME.")
	fs.StringVar(&license,
		"license", "", "Project license (e.g. 'newbsd').")
	fs.BoolVar(&AppConfig.Markdown,
		"markdown", false, "Markdown-enabled README.")
	fs.BoolVar(&(AppConfig.MakeTests),
		"test", AppConfig.MakeTests, "Produce test files with Go files.")
	fs.BoolVar(&(nohost), "nohost", false, "Don't use repository host.")
	fs.BoolVar(&(norepo), "norepo", false, "Don't start a repository.")
	fs.BoolVar(&VERBOSE,
		"v", false, "Verbose output.")
	fs.IntVar(&DEBUG_LEVEL,
		"debug", -1, "Change the amout of debug output.")
	fs.BoolVar(&help,
		"help", false, "Show this message.")
	fs.Usage = func() {
		printUsageHead()
		fs.PrintDefaults()
	}
	return fs
}

// Returns a path to the directory containing Gonew's source (and the templates/ directory).
// This function first searches locates the 'gonew' executable in PATH. Then locates the directory
// it was built from (either GOROOT or a subdirectory of GOPATH).
func FindGonew() (dir string, err error) {
	var bin string
	if bin, err = exec.LookPath("gonew"); err != nil {
		return
	}
	var bindir string
	if bindir, err = filepath.Abs(filepath.Dir(bin)); err != nil {
		return
	}
	var gobin string
	var gopath string
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, "GOBIN=") {
			if gobin, err = filepath.Abs(env[6:]); err != nil {
				return
			}
		}
		if strings.HasPrefix(env, "GOPATH=") {
			if gopath, err = filepath.Abs(env[7:]); err != nil {
				return
			}
		}
		if len(gobin) > 0 && len(gopath) > 0 {
			break
		}
	}
	if bindir == gobin {
		rootdir := filepath.Join(runtime.GOROOT(), "github.com", "bmtasuo", "gonew")
		var stat os.FileInfo
		switch stat, err = os.Stat(rootdir); {
		case err != nil:
			if err != syscall.ENOENT {
				return
			}
		case stat != nil:
			dir = rootdir
			return
		default:
			panic("unreachable")
		}
	}
	var stat os.FileInfo
	if filepath.Base(bindir) != "bin" {
		err = fmt.Errorf("%s not under a GOPATH", bindir)
	} else if stat, err = os.Stat(filepath.Join(filepath.Dir(bindir), "src")); err != nil {
		err = fmt.Errorf("%s not under a GOPATH: %s", bindir)
	} else if !stat.IsDir() {
		err = fmt.Errorf("%s not under a GOPATH", bindir)
	} else if stat, err = os.Stat(filepath.Join(filepath.Join(filepath.Dir(bindir), "src", "github.com", "bmatsuo", "gonew"))); err != nil {
		if stat, err = os.Stat(filepath.Join(filepath.Join(filepath.Dir(bindir), "src", "gonew"))); err == nil {
			dir = filepath.Join(filepath.Join(filepath.Dir(bindir), "src", "gonew"))
		}
	} else {
		dir = filepath.Join(filepath.Join(filepath.Dir(bindir), "src", "github.com", "bmatsuo", "gonew"))
	}
	if err != nil {
		return
	}
	return
}

type Request int

const (
	NilRequest Request = iota
	ProjectRequest
	LibraryRequest
)

var RequestedFile File
var RequestedProject Project

func parseArgs() Request {
	var fs = setupFlags()
	fs.Parse(os.Args[1:])
	userepo = !norepo
	usehost = !nohost
	if DEBUG_LEVEL >= 0 {
		DEBUG = true
	}
	if help {
		fs.Usage()
		os.Exit(0)
	}
	var narg = fs.NArg()
	if narg < 1 {
		ArgumentError("Missing TYPE argument")
		os.Exit(1)
	}
	if narg < 2 {
		ArgumentError("Missing NAME argument")
		os.Exit(1)
	}
	ptype = fs.Arg(0)
	name = fs.Arg(1)

	if target == "" {
		target = DefaultTarget(name)
	}

	imports := make([]string, 0, 1+strings.Count(importlist, ":"))
	for _, s := range strings.Split(importlist, ":") {
		if pkg := strings.TrimFunc(s, unicode.IsSpace); pkg != "" {
			imports = append(imports, pkg)
		}
	}
	if len(imports) == 0 {
		imports = nil
	}

	var (
		file = File{
			Name: name, Pkg: "main",
			Repo: AppConfig.Repo, License: AppConfig.License,
			User: AppConfig.HostUser, Host: AppConfig.Host,
			ImportLibs: imports}
		project = Project{
			Name: name, Target: target,
			Type: NilProjectType, License: AppConfig.License, Remote: remote,
			Host: AppConfig.Host, User: AppConfig.HostUser,
			Repo:       AppConfig.Repo,
			ImportLibs: imports,
			Markdown:   AppConfig.Markdown}
		produceProject = true
		licObj         = NilLicense
		repoObj        = NilRepo
		hostObj        = NilHost
	)
	switch ptype {
	case "cmd":
		project.Type = CmdType
	case "pkg":
		project.Type = PkgType
	case "lib":
		produceProject = false
	default:
		fmt.Fprintf(os.Stderr, "Unknown TYPE argument %s\n", ptype)
		os.Exit(1)
	}
	switch license {
	case "":
		break
	case "newbsd":
		licObj = NewBSDLicense
	default:
		fmt.Fprintf(os.Stderr, "Unknown license %s\n", ptype)
		os.Exit(1)
	}
	if userepo {
		switch repo {
		case "":
			break
		case "git":
			repoObj = Git
		case "mercurial":
			repoObj = Hg
		default:
			fmt.Fprintf(os.Stderr, "Unknown repository type %s\n", repo)
			os.Exit(1)
		}
		switch host {
		case "":
			break
		case "github":
			hostObj = GitHubHost
			repoObj = Git
		/*
		   case "googlecode":
		       hostObj = GoogleCodeType
		       repoObj = HgType
		*/
		default:
			fmt.Fprintf(os.Stderr, "Unknown respository host %s\n", host)
			os.Exit(1)
		}
	}
	if produceProject {
		// TODO check target for improper characters.
		if user != "" {
			project.User = user
		}
		if licObj != NilLicense {
			project.License = licObj
		}
		if !usehost {
			project.Host = NilHost
		} else if hostObj != NilHost {
			project.Host = hostObj
		}
		if userepo && repoObj != NilRepo {
			project.Repo = repoObj
		}
		RequestedProject = project
		return ProjectRequest
	} else {
		if narg < 3 {
			ArgumentError("Missing PKG argument")
			os.Exit(1)
		}
		file.Pkg = fs.Arg(2)
		if user != "" {
			file.User = user
		}
		if licObj != NilLicense {
			file.License = licObj
		}
		if !usehost {
			file.Host = NilHost
		} else if hostObj != NilHost {
			file.Host = hostObj
		}
		if userepo && repoObj != NilRepo {
			file.Repo = repoObj
		}
		RequestedFile = file
		return LibraryRequest
	}
	return NilRequest
}

func init() {
	var err error
	if err = TouchConfig(); err != nil {
		fmt.Fprint(os.Stderr, err.Error(), "\n")
		os.Exit(1)
	}
	Verbose("Parsing config file.\n")
	ReadConfig()

	Verbose("Locating source directory.\n")
	root, err := FindGonew()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error finding gonew source directory: %v\n", err)
		os.Exit(1)
	}
	GonewRoot = root
	TemplateRoot = filepath.Join(root, "templates")

	Verbose("Searching for templates.\n")
	if Templates, err = FindTemplates(); err != nil {
		panic(err)
	}
}

func main() {
	switch request := parseArgs(); request {
	case ProjectRequest:
		if DEBUG {
			fmt.Printf("Project requested %v\n", RequestedProject)
		} else if VERBOSE {
			fmt.Printf("Generating project %s\n", RequestedProject.Name)
		}
		if err := RequestedProject.Create(); err != nil {
			fmt.Fprint(os.Stderr, err.Error(), "\n")
			os.Exit(1)
		}
	case LibraryRequest:
		if DEBUG {
			fmt.Printf("Library requested %v\n", RequestedFile)
		} else if VERBOSE {
			fmt.Printf("Generating library %s (package %s)\n",
				RequestedFile.Name+".go", RequestedFile.Pkg)
		}
		if err := RequestedFile.Create(); err != nil {
			fmt.Fprint(os.Stderr, err.Error(), "\n")
			os.Exit(1)
		}
		if AppConfig.MakeTests {
			if err := RequestedFile.TestFile().Create(); err != nil {
				fmt.Fprint(os.Stderr, err.Error(), "\n")
				os.Exit(1)
			}
		}
	}
}
