package generator

func (g *Generator) generateProvider() {

	content, err := render(
		"provider.go.tmpl",
		nil,
	)

	if err != nil {
		panic(err)
	}

	g.project.AddFile(
		"app/Providers/app_provider.go",
		content,
	)
}