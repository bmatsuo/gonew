{{ if .Project.ReadmeIsMarkdown }}
[install go]: http://golang.org/install.html "Install Go"
[the godoc url]: http://localhost:6060/pkg/{{.Project.HostRepoString}}/ "the Godoc URL"

{{ end }}About {{.Project.Name}}
{{ if .Project.ReadmeIsMarkdown }}============={{ end }}

{{.Description}}

Documentation
{{ if .Project.ReadmeIsMarkdown }}============={{ end }}
{{ if .Project.IsCommand }}
Usage
{{ if .Project.ReadmeIsMarkdown }}-----{{ end }}

Run {{.Project.Name}} with the command

    {{.Project.Target}} [options]
{{ end }}
Prerequisites
{{ if .Project.ReadmeIsMarkdown }}-------------{{ end }}

<<<<<<< HEAD
{{ if .Project.ReadmeIsMarkdown }}[{{ end }}Install Go{{ if .Project.ReadmeIsMarkdown }}][]{{ else }}(http://golang.org/install.html){{ end }}.
=======
[Install Go](http://golang.org/doc/install.html). 
>>>>>>> 75f2280e7bf9fba5b614cb1fe9bc13d9bae70495

Installation
{{ if .Project.ReadmeIsMarkdown }}-------------{{ end }}

Use goinstall to install {{.Project.Name}}

    goinstall {{.Project.HostRepoString}}

General Documentation
{{ if .Project.ReadmeIsMarkdown }}---------------------{{ end }}

Use godoc to vew the documentation for {{.Project.Name}}

    godoc {{.Project.HostRepoString}}

Or alternatively, use a godoc http server

    godoc -http=:6060
<<<<<<< HEAD
{{ if .Project.ReadmeIsMarkdown }}
and visit [the Godoc URL][]
{{ else }}
and view the Godoc URL http://localhost:6060/pkg/{{.Project.HostRepoString}}/.
=======

and view the docs [here](http://localhost:6060/pkg/{{.Project.HostRepoString}}/)
>>>>>>> 75f2280e7bf9fba5b614cb1fe9bc13d9bae70495

Author
{{ if .Project.ReadmeIsMarkdown }}======{{ end }}

{{name}} {{ if .Project.ReadmeIsMarkdown }}&lt;{{else}}<{{end}}{{email}}{{ if .Project.ReadmeIsMarkdown }}&gt;{{else}}>{{end}}

Copyright & License
{{ if .Project.ReadmeIsMarkdown }}==================={{ end }}

Copyright (c) {{year}}, {{name}}.
All rights reserved.
