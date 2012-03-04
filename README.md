[install go]: http://golang.org/doc/install.html "Install Go" 
[go environment]: http://golang.org/doc/install.html#environment "Go environment"

[issues]: https://github.com/bmatsuo/gonew/issues "Github issues"

About gonew
===========

The idea behind gonew is to quickly and easily generate new Go projects
and repositoriesl. Gonew is inspired by Perl's `h2xs` application.

*Note: Gonew has Mercurial support, but no Google Code support. I'm github for
life (sorry). To help add support more hosts, read the "Help out" section and
contact me. I'll help write the necessary code.*

Features
========

- Create packages, commands, and single files (w/ test files).

- Create a project directory and fill it with stub files generated from templates.

- Initialize a repository and commit the stubs.

- If a new repository has already been created on a host like Github, Gonew can
push the initial commit for you.

Prerequisites
=============

[Install Go](http://golang.org/doc/install.html) 

It's recommended you set the environment variable [GOROOT][go environment] to be
the root directory of your local Go repository.


Documentation
=============
Install
-------

**Option 1**

    go get github.com/bmatsuo/gonew

**Option 2**

    git clone git@github.com:bmatsuo/gonew.git
    go install gonew

Examples
--------

Create a new package project

    gonew -target=mp3lib pkg go-mp3lib

Create a new library and test file belonging to a given package.

    gonew lib decode mp3lib

Create a new command line utility, and initialize an empty (newly
created) github repository with the local project repository.

    gonew -remote=git@github.com:bmatsuo/goplay.git cmd goplay

General Documentation
---------------------

For information about command line options.

    gonew -h

For more detailed information

    godoc github.com/bmatsuo/gonew

Templates
---------

Gonew uses customizable templates. You can browse the repository to see the
existing templates. Specify a directory of custom templates in the configuration
file `~/.gonewrc`. The configuration variable is `templates`, in section
`[general]`.

**Caveat 1** To override a default template, the *filename* of the corresponding
custom template must be identical. Templates with arbitrary filenames can be used
through reference in custom templates with standard names.

**Caveat 2** Custom templates cannot reference default templates. In order to
reference a default gonew template, the template must be copied/linked into the
custom template directory.

Help out
========

Licenses
--------

Licenses can be added very easily

If you have experience with licenses other than New BSD and want to help add
them to Gonew, [create a new issue][issues] and provide links to the relavent
license documentation (and the license text itself).

Versioning systems
------------------

The Git support in Gonew is tested and stable. Mercurial support has been
implemented as well.

People can help add other versioning systems by [creating an issue][issues] and
providing some general information about initializing new repositories,
committing, and pushing to hosts.

Hosts
-----

Currently, Github is the only host usable with Gonew.

To help, [create a new issue][issues] and provide info about versioning systems
used by the host (at least the default one), repository initialization, and
generated import paths.

Templates
---------

If you have any suggestions regarding the contents of the default Gonew templates
please [create a new issue][issues].

Author
======

Bryan Matsuo <bryan.matsuo@gmail.com>

Copyright & License
===================

Copyright (c) 2011, Bryan Matsuo.
All rights reserved.

Use of this source code is governed by a BSD-style license that can be
found in the LICENSE file.
