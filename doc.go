// Copyright 2011, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Gonew generates new Go projects. The produced projects contain stub files and
can optionally initialize repositories and add files to them.

Usage:

    gonew [options] project target

Arguments:

	project: The type of project to generate
	target: The name from which filenames are based

Options:

	-config="": specify config path
	-env="": specify a user environment
	-pkg="": specify a package name

Examples:

    gonew pkg go-mp3lib
    gonew -pkg mp3lib lib decode
    gonew cmdtest goplay

Configuration:

Gonew is configured via a JSON file stored in ~/.config/gonew.json. An example
can be found in gonew.json.example The configuration file specifies
environments, projects, and the locations of externally defined templates. An
environement hold information used in template rendering like user metadata and
import paths for created projects. A project configuration describes the files
contained in a project and script hooks to execute on file creation.
Environments can inherit/override other environments and projects can
inherit/override from other projects.

Custom Templates:

Users can define their own set of custom templates. This is done by adding
entries to the ExternalTemplates array in the configuration file. Templates
can make use of the standard gonew templates (in the "templates" directory).
Templates must have the .t2 file extension to be recognized by Gonew.

Template Functions:

Templates in Gonew have acces to a small library of helper functions Here is
list of all available template functions.

		name: the user's name specified in the environment
		email: the user's email specified in the environment
		year: the year in 4-digit format
		time: the time with an optional format string argument
		date: the date with an optional format string argument
		import: import an arbitrary number of packages into go source files
		equal: compare two values for equality

*/
package documentation
