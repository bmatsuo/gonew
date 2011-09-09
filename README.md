*gonew version 0.1_4*

Gonew generates new Go project directories.

About gonew
===========

Note: There is Mercurial support in Gonew, but no Google Code support.
I don't know the best practices concerning mercurial. If you can help.
Read the "Help out" section and contact me.

The idea behind gonew is to quickly and easily generate new Go projects
and repositories that can be installed immediately via Goinstall. The
inspiration behind gonew is the h2xs application for Perl.

A Project directory produced by gonew contains stub files and an
initialized repository. Git repositories can be automatically pushed to
new github repositories. Gonew can be used to create both new packages
and new commands.

Prerequisites
=============

You must have Go installed (http://golang.org/) and the $GOROOT
environment variable must be set to the Go source root directory.

Documentation
=============
Install
-------

Install gonew using Goinstall. It's the easiest way to ensure the templates
are correctly installed.

    goinstall github.com/bmatsuo/gonew

For the more adventurous, who like making life harder, clone the git
repository and build gonew yourself.

    git clone git@github.com:bmatsuo/gonew.git
    cd gonew
    gomake install

To finish the installation and make the templates visible to gonew,
create the file ~/.gonewrc (there is an example in the repo) and
edit the "templates" setting so that the path point to the templates/
subdirectory of the repository.

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

    gonew -help

For more detailed information

    godoc github.com/bmatsuo/gonew

Help out
========

If you have experience using licenses other than new BSD and LGPL and
want to help add support for them into Gonew, get in touch with me and
we can quickly get that done.

I only use git/github for version control and have never seriously used
mercurial, svn, etc. or any of their web-based hosts. If you can help
write handlers for repos and hosts of those kind, please let me know.
Your help would be greatly appericiated.

If you have any suggestions regarding the contents of the gonew templates
please make an issue on github (https://github.com/bmatsuo/gonew/issues).

Author
======

Bryan Matsuo <bmatsuo@soe.ucsc.edu>

Copyright & License
===================

Copyright (c) 2011, Bryan Matsuo.
All rights reserved.

Use of this source code is governed by a BSD-style license that can be
found in the LICENSE file.
