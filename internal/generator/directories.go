package generator

func (g *Generator) generateDirectories() {

	dirs := []string{
		"app",

		"app/Controllers",
		"app/Services",
		"app/Repositories",
		"app/Models",
		"app/Middleware",
		"app/Requests",
		"app/Providers",

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