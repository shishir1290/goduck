package generator

func (g *Generator) generateMain() {

	content, err := render(
		"main.go.tmpl",
		map[string]any{
			"Module": g.project.Name,
			"Port":   g.program.Server.Port,
		},
	)

	if err != nil {
		panic(err)
	}

	g.project.AddFile(
		"main.go",
		content,
	)
}