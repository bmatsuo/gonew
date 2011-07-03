/*
Gonew generates new Go project directories. Its produced project
directories contain stub files and initialized repositories (only
git/github supported now). It can be used to create new packages and
commands.

The gonew configuration file is stored at ~/.gonewrc. It is generated the
first time you run gonew. Command line options can be used to override
some details of the configuration file.

Usage:

    gonew [options] TYPE NAME

Arguments:

    TYPE
            The type of the new project ("pkg" and "cmd" supported).

    NAME
            The name of the new project/repo.

Options:

    -host=""
            Repository host if any (currently, "github" is the only
            supported host).

    -repo="git"
            Repository type (currently, "git" is the only supported
            repository type).

    -target=""
            Makefile target. The executable name in case the argument
            TYPE is "cmd", package name in case of "pkg". The default
            value based on the argument NAME.

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
*/
package documentation
