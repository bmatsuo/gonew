package main
/* 
*  File: gonew_main.go
*  Author: Bryan Matsuo [bmatsuo@soe.ucsc.edu] 
*  Created: Sat Jul  2 19:17:53 PDT 2011
*  Usage: gonew [options]
*/
import (
    "os"
    //"io"
    //"log"
    "fmt"
    "flag"
    "github.com/kr/pretty.go"
)

const (
    DEBUG = true
	DEBUG_LEVEL = 0
)
var (
    usage = "gonew [options] TYPE NAME"
    printUsageHead = func () { fmt.Fprint(os.Stderr, "\n", usage, "\n\n") }
    name   string
    ptype  string
    repo   string
    host   string
    target string
    help   bool
)

func setupFlags() *flag.FlagSet {
    var fs = flag.NewFlagSet("gonew", flag.ExitOnError)
    fs.StringVar(&repo,
        "repo", "git", "Repository type (e.g. 'git').")
    fs.StringVar(&host,
        "host", "", "Remote repository remote host if any (e.g. 'github').")
    fs.StringVar(&target,
        "target", "", "Makefile target. Default based on NAME.")
    fs.BoolVar(&help,
        "help", false, "Show this message.")
    var usageTemp = fs.Usage
    fs.Usage = func () {
        printUsageHead()
        usageTemp()
    }
    return fs
}

func parseArgs() Project {
    var fs = setupFlags()
    fs.Parse(os.Args[1:])
    if help {
        fs.Usage()
        os.Exit(0)
    }
    var narg  = fs.NArg()
    if narg < 1 {
        fmt.Fprint(os.Stderr, "missing TYPE argument")
        os.Exit(1)
    }
    if narg < 2 {
        fmt.Fprint(os.Stderr, "mising NAME argument")
        os.Exit(1)
    }
    ptype = fs.Arg(0)
    name = fs.Arg(1)

    if target == "" {
        target = DefaultTarget(name)
    }

    var project = Project{
            Name:name, Target:target, Type: NilProjectType,
            Host:AppConfig.Host, Repo:AppConfig.Repo}
    switch ptype {
    case "cmd":
        project.Type = CmdType
    case "pkg":
        project.Type = PkgType
    default:
        fmt.Fprintf(os.Stderr, "Unknown TYPE %s\n", ptype)
        os.Exit(1)
    }
    switch host {
    case "":
        break
    case "github":
        project.Host = GitHubHost
        project.Repo = GitType
    default:
        fmt.Fprintf(os.Stderr, "Unknown HOST %s\n", host)
        os.Exit(1)
    }
    if project.Host == NilRepoHost {
        switch repo {
        case "":
            break
        case "git":
            project.Repo = GitType
        default:
            fmt.Fprintf(os.Stderr, "Unknown REPO %s\n", repo)
            os.Exit(1)
        }
    }
    // TODO check target for improper characters.
    return project
}

func main() {
    var project = parseArgs()
    var errTouch = TouchConfig()
    if errTouch != nil {
        fmt.Print(errTouch.String())
        os.Exit(1)
    }
	ReadConfig()
    fmt.Printf("%v\n", pretty.Formatter(project))
    var errCreate = project.Create()
    if errCreate != nil {
        fmt.Fprintf(os.Stderr, errCreate.String())
        os.Exit(1)
    }
}
