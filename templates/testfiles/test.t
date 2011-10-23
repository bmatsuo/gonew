{{ template "go._package.t" $ }}
{{ template "go._head.t" $ }}
import (
    "testing"
)

func Test{{.test}}(T *testing.T) {
}
