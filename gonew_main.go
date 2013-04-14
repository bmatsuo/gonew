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
	"github.com/bmatsuo/gonew/config"
	"github.com/bmatsuo/gonew/project"
	"github.com/bmatsuo/gonew/templates"

	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"syscall"
	"text/template"
	"time"
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
	usage = `
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

func Verbose(msg ...string) {
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
		return
	}
	gosrc := filepath.Join(filepath.Dir(bindir), "src")
	stat, err = os.Stat(gosrc)
	if err != nil || !stat.IsDir() {
		err = fmt.Errorf("%s does not exist", gosrc)
		return
	}

	dirCandidates := [...]string{
		filepath.Join(gosrc, "pkg", "github.com", "bmatsuo", "gonew"),
		filepath.Join(gosrc, "github.com", "bmatsuo", "gonew"),
		filepath.Join(gosrc, "gonew"),
	}
	for _, dir = range dirCandidates {
		stat, err = os.Stat(dir)
		if err == nil && stat.IsDir() {
			Verbose("Found gonew source directory: ", dir, "\n")
			return
		}
	}
	err = fmt.Errorf("Couldn't find gonew source directory")
	return
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

func check(err error) error {
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func checkFatal(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func logJson(v ...interface{}) {
	w := make([]interface{}, 0, len(v))
	w = append(w, v[:len(v)-1]...)
	p, _ := json.MarshalIndent(v[len(v)-1], " ", "\t")
	w = append(w, string(p))
	fmt.Println(w...)
}

func executeHooks(ts templates.Interface, tenv templates.Environment, hooks ...*config.HookConfig) {
	for _, hook := range hooks {
		cwd, err := tenv.RenderTextAsString(ts, "cwd_", hook.Cwd)
		checkFatal(err)
		// fmt.Println("cd", cwd)
		for _, _cmd := range hook.Commands {
			cmd, err := tenv.RenderTextAsString(ts, "cmd_", _cmd)
			checkFatal(err)
			// fmt.Println("bash", "-c", cmd)
			shell := exec.Command("bash", "-c", cmd)
			shell.Dir = cwd
			shell.Stdin = os.Stdin
			shell.Stdout = os.Stdout
			shell.Stderr = os.Stderr
			checkFatal(shell.Run()) // TODO clean exit
		}
	}
}

type File2 struct {
	path    string
	content []byte
}

func funcsV2(env *config.EnvironmentConfig) template.FuncMap {
	return template.FuncMap{
		"name":  func() string { return env.User.Name },
		"email": func() string { return env.User.Email },

		"year": func() string { return time.Now().Format("2006") },
		"time": func(format ...string) string {
			if len(format) == 0 {
				format = append(format, time.RFC1123)
			}
			return time.Now().Format(format[0])
		},
		"date": func(format ...string) string {
			if len(format) == 0 {
				format = append(format, "Jan 02, 2006")
			}
			return time.Now().Format(format[0])
		},

		"import": func(pkgs ...string) string {
			if len(pkgs) == 0 {
				return `import ()`
			}
			if len(pkgs) == 1 {
				return `import "` + pkgs[0] + `"`
			}
			s := "import (\n"
			for _, pkg := range pkgs {
				s += "\t" + pkg + "\n"
			}
			s += ")"
			return s
		},
		"equal": func(v1, v2 interface{}) bool {
			return reflect.DeepEqual(reflect.ValueOf(v1), reflect.ValueOf(v2))
		},
	}
}

type options struct {
	env     string
	project string
	target  string
	pkg     string
	config  string
}

func parseOptionsV2() *options {
	opts := new(options)
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	fs.StringVar(&opts.env, "env", "", "specify a user environment")
	fs.StringVar(&opts.pkg, "pkg", "", "specify a package name")
	fs.StringVar(&opts.config, "config", "", "specify config path")
	fs.Parse(os.Args[1:])

	args := fs.Args()
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "usage:", os.Args[0], "[options] project target")
		os.Exit(1)
	}
	opts.project, opts.target = args[0], args[1]
	if opts.pkg == "" {
		opts.pkg = opts.target
		if strings.HasPrefix(opts.pkg, "go-") {
			opts.pkg = opts.pkg[3:]
		}
		if strings.HasSuffix(opts.pkg, ".go") {
			opts.pkg = opts.pkg[:len(opts.pkg)-3]
		}
	}

	return opts
}

func initConfigV2(path string) (*config.GonewConfig2, error) {
	if path == "" {
		home := os.Getenv("HOME")
		path = filepath.Join(home, ".config", "gonew.json")
	}
	conf := new(config.GonewConfig2)
	err := conf.UnmarshalFileJSON(path)
	if err != nil {
		if err, ok := err.(*os.PathError); ok {
			if err.Err == syscall.ENOENT || err.Err == os.ErrNotExist {
				fmt.Fprintln(os.Stderr, "couldn't find config")

				bufr := bufio.NewReader(os.Stdin)
				fmt.Print("Enter your name: ")
				p, _, err := bufr.ReadLine()
				checkFatal(err)
				name := strings.TrimFunc(string(p), unicode.IsSpace)
				fmt.Print("Enter your email: ")
				p, _, err = bufr.ReadLine()
				checkFatal(err)
				email := strings.TrimFunc(string(p), unicode.IsSpace)
				fmt.Print("Base import path (e.g. github.com/bmatsuo): ")
				p, _, err = bufr.ReadLine()
				checkFatal(err)
				baseImportPath := strings.TrimFunc(string(p), unicode.IsSpace)
				// TODO locate gonew.json.example
				examplePath := "gonew.json.example"
				checkFatal(conf.UnmarshalFileJSON(examplePath))
				conf.Environments = config.EnvironmentsConfig{
					"norm": &config.EnvironmentConfig{
						BaseImportPath: baseImportPath,
						User: &config.EnvironmentUserConfig{
							Name:  name,
							Email: email,
						},
					},
				}

				err = conf.MarshalFileJSON(path)
				if err != nil {
					return nil, err
				}
			} else {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return conf, nil
}

func mainv2() {

	// parse command line options/args
	opts := parseOptionsV2()
	projectName := opts.target
	packageName := opts.pkg
	envName := opts.env
	projType := opts.project

	// read the config file
	conf, err := initConfigV2(opts.config)
	checkFatal(err)

	// initialize project
	env, err := conf.Environment(envName)
	checkFatal(err)
	project.BaseImportPath = env.BaseImportPath
	proj := project.New(projectName, packageName, env)
	projContext := project.Context("", "", proj)
	projTemplEnv := templates.Env(projContext)
	projConfig, err := conf.Project(projType)
	checkFatal(err)

	// initialize template environment
	ts := templates.New(".t2")
	checkFatal(ts.Funcs(funcsV2(env)))

	// read templates
	src := templates.SourceDirectory("/Users/bryan/Go/src/github.com/bmatsuo/gonew/templates")
	checkFatal(ts.Source(src))
	for i := len(conf.ExternalTemplates) - 1; i >= 0; i-- {
		src := templates.SourceDirectory(conf.ExternalTemplates[i])
		checkFatal(ts.Source(src))
	}

	if projConfig.Hooks != nil {
		// fmt.Println("PRE")
		executeHooks(ts, projTemplEnv, projConfig.Hooks.Pre...)
	}

	// generate files. buffer all output then write.
	files := make([]*File2, 0, len(projConfig.Files))
	for _, file := range projConfig.Files {
		_relpath, err := projTemplEnv.RenderTextAsString(ts, "pre_", file.Path)
		checkFatal(err)
		relpath := string(_relpath)
		filename := filepath.Base(relpath)
		filetype := file.Type

		fileContext := project.Context(filename, filetype, proj)
		fileTemplEnv := templates.Env(fileContext)
		fileBuf := new(bytes.Buffer)
		for _, t := range file.Templates {
			if nil != check(fileTemplEnv.Render(fileBuf, ts, t)) {
				fileBuf = nil
				break
			}
		}

		if fileBuf != nil {
			f := &File2{relpath, fileBuf.Bytes()}
			files = append(files, f)
		} else {
			// TODO clean exit
		}
	}
	for _, file := range files {
		dir := filepath.Dir(file.path)

		// fmt.Println("mkdir", "-p", dir)
		err := os.MkdirAll(dir, 0755|os.ModeDir)
		checkFatal(err) // TODO clean exit

		// fmt.Println("cat", ">", file.path)
		// fmt.Println(string(file.content))
		handle, err := os.OpenFile(file.path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
		checkFatal(err) // TODO clean exit
		_, err = handle.Write(file.content)
		checkFatal(err) // TODO clean exit
		err = handle.Close()
	}

	if projConfig.Hooks != nil {
		// fmt.Println("POST")
		executeHooks(ts, projTemplEnv, projConfig.Hooks.Post...)
	}
}

func main() {
	mainv2()
}
