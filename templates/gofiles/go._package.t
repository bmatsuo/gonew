package {{ if .IsCommand }}main{{ else }}{{ .gotarget }}{{ end }}
