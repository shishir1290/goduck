package generator

import (
	"path/filepath"
	"strings"

	"github.com/shishir1290/goduck/internal/ast"
)

type ModuleTemplate struct {
	Name       string
	Package    string
	Controller string
	Service    string
	Repository string
	Routes     []*ast.Route
}

func (g *Generator) generateModules() {

	for _, module := range g.program.Modules {

		template := &ModuleTemplate{
			Name:       module.Name,
			Package:    strings.ToLower(module.Name),
			Controller: module.Controller,
			Service:    module.Service,
			Repository: module.Repository,
			Routes:     module.Routes,
		}

		moduleDir := filepath.Join(
			"app",
			strings.ToLower(module.Name),
		)

		// Controller
		content, err := render(
			"controller.go.tmpl",
			template,
		)

		if err != nil {
			panic(err)
		}

		g.project.AddFile(
			filepath.Join(
				moduleDir,
				"controller.go",
			),
			content,
		)

		// Service
		if module.Service != "" {

			content, err := render(
				"service.go.tmpl",
				template,
			)

			if err != nil {
				panic(err)
			}

			g.project.AddFile(
				filepath.Join(
					moduleDir,
					"service.go",
				),
				content,
			)
		}

		// Repository
		if module.Repository != "" {

			content, err := render(
				"repository.go.tmpl",
				template,
			)

			if err != nil {
				panic(err)
			}

			g.project.AddFile(
				filepath.Join(
					moduleDir,
					"repository.go",
				),
				content,
			)
		}
	}
}