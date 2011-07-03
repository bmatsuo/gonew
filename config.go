package main
/* 
*  File: config.go
*  Author: Bryan Matsuo [bmatsuo@soe.ucsc.edu] 
*  Created: Sat Jul  2 23:09:50 PDT 2011
*/
import (
    //"goconf.googlecode.com/hg"    // This does not work for some reason.
    //"conf"
    "github.com/kless/goconfig/config"
    "os"
    "path/filepath"
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
