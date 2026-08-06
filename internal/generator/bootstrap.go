package generator

func (g *Generator) generateBootstrap() {

	content, err := render(
		"bootstrap.go.tmpl",
		nil,
	)

	if err != nil {
		panic(err)
	}

	g.project.AddFile(
		"bootstrap/app.go",
		content,
	)
}