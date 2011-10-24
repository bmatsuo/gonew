{{ template "go._package.t" $ }}
{{ template "go._head.t" $ }}
{{ import .Imports }}

var opt = parseFlags()

{{ func "main" }}
