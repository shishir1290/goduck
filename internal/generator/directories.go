package generator

func (g *Generator) generateDirectories() {

	dirs := []string{

		"app",

		"app/Controllers",
		"app/Models",
		"app/Services",
		"app/Repositories",
		"app/Middleware",
		"app/Requests",

		"bootstrap",

		"config",

		"routes",

		"storage",

		"public",
	}

	for _, dir := range dirs {
		g.project.AddDirectory(dir)
	}
}