{{ template "go._head.t" $ }}

// Package {{.Project.Target}} does ....
{{ template "go._package.t" $ }}
{{ import .Imports }}
