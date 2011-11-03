About {{.Project.Name}}
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

[Install Go](http://golang.org/doc/install.html). 

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

and view the docs [here](http://localhost:6060/pkg/{{.Project.HostRepoString}}/)

Author
{{ if .Project.ReadmeIsMarkdown }}======{{ end }}

{{name}} <{{email}}>

Copyright & License
{{ if .Project.ReadmeIsMarkdown }}==================={{ end }}

Copyright (c) {{year}}, {{name}}.
All rights reserved.

