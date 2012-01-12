{{ template "go._package.t" $ }}
{{ template "go._head.t" $ }}
{{ import .Imports "//log" "//fmt" }}

var opt Options

{{ init "opt = parseFlags()" }}

{{ main }}
