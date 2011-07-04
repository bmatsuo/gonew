package main
/* 
*  File: config.go
*  Author: Bryan Matsuo [bmatsuo@soe.ucsc.edu] 
*  Created: Sat Jul  2 23:09:50 PDT 2011
*/
import (
    "os"
    "fmt"
    "log"
    "bytes"
    "bufio"
    "path/filepath"
    //"goconf.googlecode.com/hg"    // This does not work for some reason.
    //"conf"
    "github.com/kless/goconfig/config"
)

var ConfigFilename = filepath.Join(os.Getenv("HOME"), ".gonewrc")

type GonewConfig struct {
    Name     string
    Email    string
    HostUser string
    Repo     RepoType
    Host     RepoHost
}

var AppConfig GonewConfig

func ReadConfig() os.Error {
    AppConfig = GonewConfig{"", "", "", NilRepoType, NilRepoHost}
    conf, err := config.ReadDefault(ConfigFilename)
    if err != nil {
        return err
    }
    var (
        repostr string
        hoststr string
    )
    AppConfig.Name, err = conf.String("variables", "name")
    AppConfig.Email, err = conf.String("variables", "email")
    AppConfig.HostUser, err = conf.String("general", "hostuser")
    repostr, err = conf.String("general", "repo")
    switch repostr {
    case "":
        AppConfig.Repo = NilRepoType
    case "git":
        AppConfig.Repo = GitType
    //case "mercurial":
    //...
    default:
        AppConfig.Repo = NilRepoType
    }
    hoststr, err = conf.String("general", "host")
    switch hoststr {
    case "":
        AppConfig.Host = NilRepoHost
    case "github":
        AppConfig.Host = GitHubHost
        AppConfig.Repo = GitType
    //case "googlecode":
    //...
    default:
        AppConfig.Host = NilRepoHost
    }
    return nil
}

func TouchConfig() os.Error {
    stat, err := os.Stat(ConfigFilename)
    var patherr *os.PathError
    switch err.(type) {
    case nil:
        patherr = nil
    case *os.PathError:
        patherr = err.(*os.PathError)
    }
    if patherr != nil && patherr.Error != os.ENOENT {
        fmt.Fprintf(os.Stderr, "Error stat'ing ~/.gonewrc. %v", patherr)
        return patherr
    } else if stat == nil || (patherr != nil && patherr.Error == os.ENOENT) {
        if DEBUG || VERBOSE {
            log.Print("Config not found. Prompting user for info.")
        }
        return MakeConfig()
    } else {
        if DEBUG {
            log.Print("~/.gonewrc found.")
        }
    }
    return nil
}

func MakeConfig() os.Error {
    var (
        c       = config.NewDefault()
        scanner = bufio.NewReader(os.Stdin)
        errScan os.Error
        buff    []byte
    )
    fmt.Printf("Enter your name: ")
    buff, _, errScan = scanner.ReadLine()
    if errScan != nil {
        return errScan
    }
    c.AddOption("variables", "name", string(bytes.TrimRight(buff, "\n")))
    fmt.Printf("Enter your email address: ")
    buff, _, errScan = scanner.ReadLine()
    if errScan != nil {
        return errScan
    }
    c.AddOption("variables", "email", string(bytes.TrimRight(buff, "\n")))
    var (
        repoName string
        repoOk   bool
    )
    for !repoOk {
        fmt.Printf("Enter a repository type ('git', or none): ")
        buff, _, errScan = scanner.ReadLine()
        if errScan != nil {
            return errScan
        }
        repoName = string(bytes.TrimRight(buff, "\n"))
        switch repoName {
        case "":
            fallthrough
        case "git":
            repoOk = true
        default:
            fmt.Printf("I didn't understand repo type %s\n", repoName)
        }
    }
    c.AddOption("general", "repo", repoName)
    var (
        hostName string
        hostOk   bool
    )
    for !hostOk {
        fmt.Printf("Enter a repo host ('github', or none): ")
        buff, _, errScan = scanner.ReadLine()
        if errScan != nil {
            return errScan
        }
        hostName = string(bytes.TrimRight(buff, "\n"))
        switch hostName {
        case "":
            fallthrough
        case "github":
            hostOk = true
        default:
            fmt.Printf("I didn't understand repo host %s\n", hostName)
        }
    }
    c.AddOption("general", "host", hostName)
    var hostuser = ""
    var hostuserOk = hostName != "github"
    for !hostuserOk {
        fmt.Printf("Enter a %s username: ", hostName)
        buff, _, errScan = scanner.ReadLine()
        if errScan != nil {
            return errScan
        }
        hostuser = string(bytes.TrimRight(buff, "\n"))
        if hostuser == "" {
            fmt.Printf("Invalid %s username '%s'\n", hostName, hostuser)
        } else {
            hostuserOk = true
        }
    }
    c.AddOption("general", "hostuser", hostuser)
    return c.WriteFile(ConfigFilename, FilePermissions, "Generated configuration for gonew.")
}
