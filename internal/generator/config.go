package generator

func (g *Generator) generateConfig() {

	files := []string{
		"config.go.tmpl",
		"app.go.tmpl",
		"server.go.tmpl",
		"database.go.tmpl",
	}

	for _, file := range files {

		content, err := render(
			file,
			map[string]any{
				"App":  g.program.App,
				"Port": g.program.Server.Port,
			},
		)

		if err != nil {
			panic(err)
		}

		name := file[:len(file)-5]

		g.project.AddFile(
			"config/"+name,
			content,
		)
	}
}