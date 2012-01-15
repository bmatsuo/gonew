// Copyright 2011, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main
/*
 *  Filename:    hosts.go
 *  Package:     main
 *  Author:      Bryan Matsuo <bmatsuo@soe.ucsc.edu>
 *  Created:     Sun Jul  3 20:28:52 PDT 2011
 *  Description: 
 */
import (
	"os/exec"
	"fmt"
)

type RepoHost int

const (
	NilHost RepoHost = iota
	GitHubHost
	//GoogleCodeHost
	//...
)

var hoststrings = []string{
	NilHost: "no host",
	GitHubHost:  "GitHub",
}

func (rh RepoHost) String() string {
	return hoststrings[rh]
}

type RemoteRepository interface {
	// The new project the remote repository is being used for.
	Project() Project
	// Setup local repository for push.
	Init(orign string) error
	// Push changes to the remote host.
	Push() error
	// Return the base repository type used by the host.
	Type() RepoType
	// Return the RepoHost that the object implements.
	Host() RepoHost
	// Return true if the host accepts README.md (markdown).
	UseMarkdown() bool
}

func VerifyRemote(remote RemoteRepository) {
	if remote.Type() != remote.Project().Repo {
		panic(fmt.Errorf("Remote/local repo type mismatch %s %s",
			remote.Type().String(), remote.Project().Repo))
	}
}

type GitHubRepo struct {
	P Project
}

func (github GitHubRepo) Project() Project  { return github.P }
func (github GitHubRepo) Type() RepoType    { return Git }
func (github GitHubRepo) Host() RepoHost    { return GitHubHost }
func (github GitHubRepo) UseMarkdown() bool { return true }
func (github GitHubRepo) Init(origin string) error {
	VerifyRemote(github)
	return exec.Command("git", "remote", "add", "origin", origin).Run()
}
func (github GitHubRepo) Push() error {
	VerifyRemote(github)
	return exec.Command("git", "push", "-u", "origin", "master").Run()
}
