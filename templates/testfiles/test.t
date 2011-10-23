{{ template "go._package.t" $ }}
{{ template "go._head.t" $ }}
import (
    "testing"
)

{{ range .Tests }}
{{ with printf "Test%s" . }}{{ func . "T *testing.T" }}{{ end }}
{{ end }}
