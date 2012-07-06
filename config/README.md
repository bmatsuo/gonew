[install go]: http://golang.org/doc/install.html "Install Go" 
[go environment]: http://golang.org/doc/install.html#environment "Go environment"

[issues]: https://github.com/bmatsuo/gonew/issues "Github issues"
[the templates package]: https://github.com/bmatsuo/gonew/tree/master/templates "The Templates Package directory"
[on gopkgdoc]: http://gopkgdoc.appspot.com/pkg/github.com/bmatsuo/gonew/config "on GoPkgDoc"

About gonew/config
==================

Gonew/config uses a json configuration. There is an
[example file](https://github.com/bmatsuo/gonew/tree/master/gonew.json.example)
in the base repository.

Documentation
=============

Because gonew unmarshals its configuration from to a struct, its most convenient
to read the documentation with `go doc` or [on GoPkgDoc][].

Common Points
=============

Templates
---------

Project file templates are specified in your config file. Each file specified a
list of templates (excluding license templates) which compose the file. These
templates are found in either one of the `"ExternalTemplates"` directories
specified in your config file, or in [The Templates Package directory][],

Version Control Systems
-----------------------

Specifying `"Environments[_].User.VersionControl"` in your configuration file
doesn't cause gonew to automatically initialize repositories in project
directories. It is mainly used for README generation templates.

Gonew can be fully integrated with versioning systems by using `"Pre"` and
`"Post"` hooks in your configuration file.

Hosts
-----

Specifying `Environments[_].User.Host` in your configuration file doesn't cause
gonew to automatically initialize repos on remote hosts. It is mainly used for
README generation templates.

Gonew can be fully integrated with hosts like Github and Google Code by using
`"Pre"` and `"Post"` hooks in your configuration file.
