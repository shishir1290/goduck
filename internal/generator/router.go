package generator

func (g *Generator) generateRouter() {

	content, err := render(
		"router.go.tmpl",
		map[string]any{
			"Module":  g.project.Name,
			"Modules": g.program.Modules,
		},
	)

	if err != nil {
		panic(err)
	}

	g.project.AddFile(
		"routes/web.go",
		content,
	)
}