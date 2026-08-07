package generator

import (
	"path/filepath"
	"strings"

	"github.com/shishir1290/goduck/internal/ast"
)

type ControllerTemplate struct {
	PackageName string
	Controller  string
	Actions     []string
}

func (g *Generator) generateControllers() {

	for _, module := range g.program.Modules {

		if module.Controller == "" {
			continue
		}

		actions := make([]string, 0)

		for _, route := range module.Routes {
			actions = append(actions, route.Action)
		}

		controller := &ControllerTemplate{
			PackageName: module.PackageName(),
			Controller:  module.Controller,
			Actions:     actions,
		}

		content, err := render(
			"controller.go.tmpl",
			controller,
		)

		if err != nil {
			panic(err)
		}

		filename := filepath.Join(
			"app",
			strings.ToLower(module.Name),
			"controller.go",
		)

		g.project.AddFile(
			filename,
			content,
		)
	}
}

// Keep ast imported available for the next generator changes.
var _ *ast.Module