// Copyright 2011, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

/* 
*  File: config.go
*  Author: Bryan Matsuo [bmatsuo@soe.ucsc.edu] 
*  Created: Sat Jul  2 23:09:50 PDT 2011
 */
import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/robfig/config"
	"os"
	"path/filepath"
	"syscall"
)

var ConfigFilename = filepath.Join(os.Getenv("HOME"), ".gonewrc")

type GonewConfig struct {
	MakeTests bool
	Name      string
	Email     string
	HostUser  string
	AltRoot   string
	License   LicenseType
	Repo      RepoType
	Host      RepoHost
	Markdown  bool
}

var AppConfig GonewConfig = GonewConfig{
	true, "", "", "", "", NilLicense, NilRepo, NilHost, false}

func ReadConfig() error {
	conf, err := config.ReadDefault(ConfigFilename)
	if err != nil {
		return err
	}
	var (
		repostr string
		hoststr string
		license string
	)
	AppConfig.Name, err = conf.String("variables", "name")
	AppConfig.Email, err = conf.String("variables", "email")
	AppConfig.HostUser, err = conf.String("general", "hostuser")
	AppConfig.AltRoot, err = conf.String("general", "templates")
	AppConfig.Markdown, err = conf.Bool("general", "markdown")

	license, err = conf.String("general", "license")
	switch license {
	case "":
		AppConfig.License = NilLicense
	case "newbsd":
		AppConfig.License = NewBSDLicense
	}

	repostr, err = conf.String("general", "repo")
	switch repostr {
	case "":
		AppConfig.Repo = NilRepo
	case "git":
		AppConfig.Repo = Git
	//case "hg":
	//...
	default:
		AppConfig.Repo = NilRepo
	}

	hoststr, err = conf.String("general", "host")
	switch hoststr {
	case "":
		AppConfig.Host = NilHost
	case "github":
		AppConfig.Host = GitHubHost
		AppConfig.Repo = Git
	//case "googlecode":
	//...
	default:
		AppConfig.Host = NilHost
	}

	return nil
}

func TouchConfig() error {
	var patherr *os.PathError

	stat, err := os.Stat(ConfigFilename)
	switch err.(type) {
	case *os.PathError:
		patherr = err.(*os.PathError)
	}

	if patherr != nil && patherr.Err != syscall.ENOENT {
		fmt.Fprintf(os.Stderr, "Error stat'ing ~/.gonewrc. %v", patherr)
		return patherr
	} else if stat == nil || (patherr != nil && patherr.Err == syscall.ENOENT) {
		fmt.Fprintln(os.Stderr, "~/.gonewrc now found. Please initialize it now.")
		return MakeConfig()
	} else {
		Debug(0, "~/.gonewrc found.")
	}
	return nil
}

func MakeConfig() error {
	c := config.NewDefault()
	scanner := bufio.NewReader(os.Stdin)

	var (
		err  error
		buff []byte
	)
	fmt.Printf("Enter your name: ")
	if buff, _, err = scanner.ReadLine(); err != nil {
		return err
	}
	c.AddOption("variables", "name", string(bytes.TrimRight(buff, "\n")))

	fmt.Printf("Enter your email address: ")
	if buff, _, err = scanner.ReadLine(); err != nil {
		return err
	}
	c.AddOption("variables", "email", string(bytes.TrimRight(buff, "\n")))

	var repoName string
	for repoOk := false; !repoOk; {
		fmt.Printf("Enter a repository type ('git', or none): ")
		if buff, _, err = scanner.ReadLine(); err != nil {
			return err
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

	var hostName string
	for hostOk := false; !hostOk; {
		fmt.Printf("Enter a repo host ('github', or none): ")
		if buff, _, err = scanner.ReadLine(); err != nil {
			return err
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

	var hostuser string
	for hostuserOk := hostName != "github"; !hostuserOk; {
		fmt.Printf("Enter a %s username: ", hostName)
		if buff, _, err = scanner.ReadLine(); err != nil {
			return err
		}
		hostuser = string(bytes.TrimRight(buff, "\n"))
		if hostuser == "" {
			fmt.Printf("Invalid %s username '%s'\n", hostName, hostuser)
		} else {
			hostuserOk = true
		}
	}
	c.AddOption("general", "hostuser", hostuser)

	var license string
	for licenseOK := false; !licenseOK; {
		fmt.Print("Enter a license type ('newbsd', or none): ")
		if buff, _, err = scanner.ReadLine(); err != nil {
			return err
		}
		license = string(bytes.TrimRight(buff, "\n"))
		switch license {
		case "":
			fallthrough
		case "newbsd":
			licenseOK = true
		default:
			fmt.Printf("Invalid %s username '%s'\n", hostName, hostuser)
		}
	}
	c.AddOption("general", "license", license)

	return c.WriteFile(ConfigFilename, FilePermissions, "Generated configuration for gonew.")
}
