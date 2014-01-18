[templates]: https://github.com/bmatsuo/gonew/tree/master/templates "templates"
[godoc.org]: http://godoc.org/github.com/bmatsuo/gonew/config "godoc.org"

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
specified in your config file, or in [templates][] directory hierarchy.

Hooks
-----

Gonew can be fully integrated with version control systems like git and
mercurial, as well as with hosts like Github and Google Code by using
`"Pre"` and `"Post"` hooks in your configuration file.
