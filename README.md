*gonew version 0.0_10*

Gonew generates new Go project directories.

About gonew
===========

The idea behind gonew is to quickly and easily generate new Go projects
and repositories that can be installed immediately via Goinstall.

Project directories produced by gonew contain stub files and initialized
repositories (only git/github supported now). Gonew can be used to create
both new packages and new commands.

Gonew has a configuration file stored at ~/.gonewrc. It is generated the
first time you run gonew. Command line options can be used to override
some details of the configuration file.

Prerequisites
=============

You must have Go installed (http://golang.org/) and the $GOROOT
environment variable must be set to the Go source root directory.

Documentation
=============
Install
-------

Installation *must* be done with goinstall. Otherwise, the templates will
not be found.

    goinstall github.com/bmatsuo/gonew

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

(C) 2011 Bryan Matsuo 

TODO - add licensing information!
