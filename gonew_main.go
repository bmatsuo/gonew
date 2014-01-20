// Copyright 2011, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

/*
*  File: gonew_main.go
*  Created: Sat Jul  2 19:17:53 PDT 2011
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
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"syscall"
	"text/template"
	"time"
	"unicode"
)

// The directory containing Gonew's source code.
var GonewRoot string // The Gonew source directory
func FindGonew() error {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		return fmt.Errorf("GOPATH is not set")
	}
	gopath = strings.SplitN(gopath, ":", 2)[0]
	GonewRoot = filepath.Join(gopath, "src", "github.com", "bmatsuo", "gonew")
	stat, err := os.Stat(GonewRoot)
	if err == nil && !stat.IsDir() {
		err = fmt.Errorf("file is not a directory: %s", GonewRoot)
	}
	return err
}

func check(err error, v ...interface{}) error {
	if err != nil {
		if len(v) == 0 {
			fmt.Println(err)
		} else {
			fmt.Printf("%s: %v", fmt.Sprint(v...), err)
		}
	}
	return err
}

func checkFatal(err error, v ...interface{}) {
	if check(err, v...) != nil {
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
		checkFatal(err, "hook cwd template")
		// fmt.Println("cd", cwd)
		for _, _cmd := range hook.Commands {
			cmd, err := tenv.RenderTextAsString(ts, "cmd_", _cmd)
			checkFatal(err, "hook template")
			// fmt.Println("bash", "-c", cmd)
			shell := exec.Command("bash", "-c", cmd)
			shell.Dir = cwd
			shell.Stdin = os.Stdin
			shell.Stdout = os.Stdout
			shell.Stderr = os.Stderr
			checkFatal(shell.Run(), "hook") // TODO clean exit
		}
	}
}

type File struct {
	path    string
	content []byte
}

func funcs(env *config.Environment) template.FuncMap {
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

func parseOptions() *options {
	opts := new(options)
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	fs.StringVar(&opts.env, "env", "", "specify a user environment")
	fs.StringVar(&opts.pkg, "pkg", "", "specify a package name")
	fs.StringVar(&opts.config, "config", "", "specify config path")
	fs.Parse(os.Args[1:])

	args := fs.Args()
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "usage:", os.Args[0], "[options] [project] target")
		os.Exit(1)
	}
	if len(args) == 1 {
		opts.target = args[0]
	} else {
		opts.project, opts.target = args[0], args[1]
	}
	if opts.pkg == "" {
		opts.pkg = opts.target
	}

	return opts
}

func readLine(r *bufio.Reader, prompt string) (string, error) {
	fmt.Print(prompt)
	p, _, err := r.ReadLine()
	line := strings.TrimFunc(string(p), unicode.IsSpace)
	return line, err
}

func initConfig(path string) (conf *config.Gonew, err error) {
	if path == "" {
		home := os.Getenv("HOME")
		path = filepath.Join(home, ".config", "gonew.json")
	}
	conf = new(config.Gonew)
	err = conf.UnmarshalFileJSON(path)
	if err == nil {
		return
	}
	switch perr, ok := err.(*os.PathError); {
	case !ok:
		return
	case perr.Err == syscall.ENOENT || perr.Err == os.ErrNotExist:
		fmt.Fprintf(os.Stderr, "configuration not found at %q\n", path)
		fmt.Fprintln(os.Stderr)
		fmt.Fprintf(os.Stderr, "if you are migrating from an older version of Gonew check out the migration guide\n")
		fmt.Fprintf(os.Stderr, "\thttps://github.com/bmatsuo/gonew/blob/v2/MIGRATION.md\n")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintf(os.Stderr, "otherwise, please take a moment to fill in the user information below\n")
		fmt.Fprintln(os.Stderr)

		var name string
		var email string
		var baseImportPath string
		bufr := bufio.NewReader(os.Stdin)
		name, err = readLine(bufr, "Your name: ")
		checkFatal(err)
		email, err = readLine(bufr, "Your email: ")
		checkFatal(err)
		baseImportPath, err = readLine(bufr, "Base import path (e.g. github.com/bmatsuo): ")
		checkFatal(err)

		examplePath := filepath.Join(GonewRoot, "gonew.json.example")
		checkFatal(conf.UnmarshalFileJSON(examplePath), "example config")
		conf.Environments = config.Environments{
			"default": &config.Environment{
				BaseImportPath: baseImportPath,
				User: &config.EnvironmentUserConfig{
					Name:  name,
					Email: email,
				},
			},
		}
		conf.Default.Environment = "default"
		err = conf.MarshalFileJSON(path)
	}
	return
}

func main() {
	checkFatal(FindGonew(), "root not found")

	// parse command line options/args
	opts := parseOptions()
	// read the config file
	conf, err := initConfig(opts.config)
	checkFatal(err, "config")

	// project metadata
	projectName := opts.target
	packageName := opts.pkg
	envName := opts.env
	if envName == "" {
		envName = conf.Default.Environment
	}
	projType := opts.project
	if projType == "" {
		projType = conf.Default.Project
	}

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
	checkFatal(ts.Funcs(funcs(env)), "templates")

	// read templates
	src := templates.SourceDirectory(filepath.Join(GonewRoot, "templates"))
	checkFatal(ts.Source(src), "templates")
	for i := len(conf.ExternalTemplates) - 1; i >= 0; i-- {
		src := templates.SourceDirectory(conf.ExternalTemplates[i])
		checkFatal(ts.Source(src), "external templates")
	}

	if projConfig.Hooks != nil {
		executeHooks(ts, projTemplEnv, projConfig.Hooks.Pre...)
	}

	// generate files. buffer all output then write.
	files := make([]*File, 0, len(projConfig.Files))
	for name, file := range projConfig.Files {
		_relpath, err := projTemplEnv.RenderTextAsString(ts, "pre_", file.Path)
		checkFatal(err, name)
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
			f := &File{relpath, fileBuf.Bytes()}
			files = append(files, f)
		} else {
			// TODO clean exit
		}
	}
	for _, file := range files {
		dir := filepath.Dir(file.path)

		// fmt.Println("mkdir", "-p", dir)
		err := os.MkdirAll(dir, 0755|os.ModeDir)
		checkFatal(err, file) // TODO clean exit

		// fmt.Println("cat", ">", file.path)
		// fmt.Println(string(file.content))
		writeMode := os.O_WRONLY | os.O_CREATE | os.O_EXCL // must create
		handle, err := os.OpenFile(file.path, writeMode, 0644)
		checkFatal(err, file) // TODO clean exit
		_, err = handle.Write(file.content)
		checkFatal(err, file) // TODO clean exit
		err = handle.Close()
	}

	if projConfig.Hooks != nil {
		// fmt.Println("POST")
		executeHooks(ts, projTemplEnv, projConfig.Hooks.Post...)
	}
}
