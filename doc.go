// Copyright 2011, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Gonew generates new Go project directories. Its produced project
directories contain stub files and initialized repositories (only
git/github supported now). It can be used to create new packages and
commands.

The gonew configuration file is stored at ~/.gonewrc. It is generated the
first time you run gonew. Command line options can be used to override
some details of the configuration file.

Usage:

    gonew [options] cmd NAME
    gonew [options] pkg NAME
    gonew [options] lib NAME PKG

Arguments:

    NAME
            The name of the new project/repo.

    PKG
            The package a library (.go file) belongs to.

Options:

    -target=""
            Makefile target. The executable name in case the argument
            TYPE is "cmd", package name in case of "pkg". The default
            value based on the argument NAME.

    -repo="git"
            Repository type (currently, "git" is the only supported
            repository type).

    -remote=""
            When passed a url to a remote repository, attempt to
            initialize the remote repository to the new project's
            repository. The url passed in must agree with the host
            specified in the config file (or by -host).

    -host=""
            Repository host if any (currently, "github" is the only
            supported host). The value supplied must agree with the
            value of -repo.

    -user=""
            Username for the repo host (necessary for "github").

    -v
            Print verbose output to the stdout (this intersects with
            some -debug output).

    -debug=-1
            When set to a non-negative value, debugging output will be
            printed.

    -help
            Print a usage message


Configuration:

The configuration for gonew is simple. The configuration can provide
default hosts, usernames, and repositories. However, it also contains the
setting of the {{name}} and {{email}} template variables.

The configuration file for gonew (~/.gonewrc) is generated on the spot if
one does not exist. So you do not need to worry about editing it for the
most part.

If you wish to write/edit the configuration file. An example configuration
file can be found at the path

    $GOROOT/src/pkg/github.com/bmatsuo/gonew/gonewrc.example

Examples:

    gonew -target=mp3lib pkg go-mp3lib
    gonew lib decode mp3lib
    gonew -remote=git@github.com:bmatsuo/goplay.git cmd goplay

*/
package documentation
