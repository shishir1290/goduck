package generator

func (g *Generator) generateGoMod() {

	content, err := render(
		"go.mod.tmpl",
		map[string]any{
			"Module": g.project.Name,
		},
	)

	if err != nil {
		panic(err)
	}

	g.project.AddFile(
		"go.mod",
		content,
	)
}