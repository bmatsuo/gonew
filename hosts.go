package main
/*
 *  Filename:    hosts.go
 *  Package:     main
 *  Author:      Bryan Matsuo <bmatsuo@soe.ucsc.edu>
 *  Created:     Sun Jul  3 20:28:52 PDT 2011
 *  Description: 
 */
import (
    "os"
    "log"
    "exec"
)

type RepoHost int
const (
    NilRepoHost RepoHost = iota
    GitHubHost
    //GoogleHost
    //...
)

type RemoteRepository interface {
    Project()          Project
    Init(orign string) os.Error // Setup local repository for push.
    Push()             os.Error // Push changes to the remote host.
    Repo()             RepoType
    Host()             RepoHost
}

func VerifyRemote(remote RemoteRepository) {
    if remote.Repo() != remote.Project().Repo {
        if DEBUG {
            log.Printf("Remote/local repo type mismatch %s %s",
                    remote.Repo(), remote.Project().Repo)
        }
        panic("typemismatch")
    }
}

type GitHubRepository struct {
    P Project
}

func (github GitHubRepository) Project() Project {
    return github.P
}
func (github GitHubRepository) Repo() RepoType {
    return GitType
}
func (github GitHubRepository) Host() RepoHost {
    return GitHubHost
}
func (github GitHubRepository) Init(origin string) os.Error {
    VerifyRemote(github)
    return exec.Command("git", "remote", "add", "origin", origin).Run()
}
func (github GitHubRepository) Push() os.Error {
    VerifyRemote(github)
    return exec.Command("git", "push", "-u", "origin", "master").Run()
}
