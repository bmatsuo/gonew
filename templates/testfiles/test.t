{{ template "go._package.t" $ }}
{{ template "go._head.t" $ }}
import (
    "testing"
)

{{ range .Tests }}
{{ with printf "Test%s" . }}{{ func . "t *testing.T" }}{{ end }}
{{ end }}
