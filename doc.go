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

	-test=false
			Test files will not be produced.

    -target=""
            Makefile target. The executable name in case the argument
            TYPE is "cmd", package name in case of "pkg". The default
            value based on the argument NAME.

    -import=""
            Colon ':' separated list of packages to include in source
            .go files. The packages are not imported in any tests, or
            the options.go file (or doc.go) created for cmd projects.

    -repo=""
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
setting of the {{name}} and {{email}} template function values.

The configuration file for gonew (~/.gonewrc) is generated on the spot if
one does not exist. So you do not need to worry about editing it for the
most part.

If you wish to write/edit your own configuration file. An example can be
found at the path

    $GOROOT/src/pkg/github.com/bmatsuo/gonew/gonewrc.example

Examples:

    gonew -target=mp3lib pkg go-mp3lib
    gonew lib decode mp3lib
    gonew -remote=git@github.com:bmatsuo/goplay.git cmd goplay

Custom Templates:

Custom file templates can be created and stored in a heirarchy separate
from the standard Gonew template heirarchy at

    $GOROOT/src/pkg/github.com/bmatsuo/gonew/templates

Supply a custom template heirarchy path in the ~/.gonewrc file. See the
example config file for more info about how to do this.

The custom heirarchy is parsed as a template.Set and can reference other
templates in the heirarchy (http://golang.org/pkg/template). Templates,
however, cannot reference templates in the default Gonew template
heirarchy. In order to do so, copies or symbolic links of the desired
templates must be made in the custom heirarchy.

Templates must be given a ".t" extension in order to be recognized and
parsed by Gonew. When generating a new project or library, Gonew uses
templates with specific names to generate the files. In order to use
custom templates, a template with the proper name must be found in the
custom heirarchy. Increase the Gonew debugging variable for information
about with templates are being executed.

Template Functions:

All templates used by Gonew have acces to a small library of simple
helper functions. The {{name}} and {{email}} variables have already been
discussed. Here is list of all available template functions.

    import [PACKAGE [...]]
            Produces an import statement which includes the packages
            specified in it arguments. The arguments can be either
            strings or slices of strings.

    func NAME [ARGUMENT [...]]
            Produces the definition of a function with zero return
            values. The arguments should be supplied as identifier-type
            pairs. The name should be a valid go identifier.

    date [FORMAT]
            Produces the current date-time as human readable string when
            called with no argument. A format can be supplied in the form
            of an example string formatting of a specific day (see "time").

    year
            Produces the current year in a four-digit format.

    name
            Produces the Gonew user's name as defined in ~/.gonewrc

    email
            Produces the Gonew user's email as defined in ~/.gonewrc

Template Contexts:

The data supplied to each template a type known as a Context. Contexts
have a defined method interface that can be called from inside a template.
Read about the Context interface using godoc

    godoc $GOROOT/src/pkg/github.com/bmatsuo/gonew Context | less

*/
package documentation
