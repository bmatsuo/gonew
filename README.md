*gonew version 0.0_6*

Gonew generates new Go project directories.

About gonew
===========

Project directories produced by gonew contain stub files and initialized
repositories (only git/github supported now). It can be used to create new
packages and commands.

The gonew configuration file is stored at ~/.gonewrc. It is generated the
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

Installation must be done with goinstall.

    goinstall github.com/bmatsuo/gonew

General Documentation
---------------------

For information about command line options.

    gonew -help

For more detailed information

    godoc github.com/bmatsuo/gonew

Help out
========

I only use git for version control and have never seriously used
mercurial, svn, ... If you can help write handlers for repos and hosts
of those kind, please let me know. Your help would be greatly
appericiated.

If you have any suggestions regarding the contents of the gonew templates
please make an issue on github (https://github.com/bmatsuo/gonew/issues).

Author
======

Bryan Matsuo <bmatsuo@soe.ucsc.edu>

Copyright & License
===================

(C) 2011 Bryan Matsuo 

TODO - add licensing information!
