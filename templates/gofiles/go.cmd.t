{{ template "go._package.t" $ }}
{{ template "go._head.t" $ }}
{{ import }}

var opt = parseFlags()

{{ func "main" }}
