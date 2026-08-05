package generator

func (g *Generator) generateMain() {

	content, err := render(
		"main.go.tmpl",
		map[string]any{
			"Port": g.program.Server.Port,
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