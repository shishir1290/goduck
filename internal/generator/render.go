package generator

import (
	"bytes"
	"embed"
	"text/template"
)

//go:embed templates/*.tmpl
var templateFS embed.FS

func render(name string, data any) (string, error) {

	tmpl, err := template.ParseFS(
		templateFS,
		"templates/"+name,
	)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer

	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}