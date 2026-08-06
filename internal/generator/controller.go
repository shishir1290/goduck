package generator

import (
	"path/filepath"
	"strings"
)

type ControllerTemplate struct {
	Controller string
	Actions    []string
}

func (g *Generator) generateControllers() {

	controllers := map[string]*ControllerTemplate{}

	for _, route := range g.program.Routes {

		name := route.Controller

		if _, ok := controllers[name]; !ok {
			controllers[name] = &ControllerTemplate{
				Controller: name,
			}
		}

		controllers[name].Actions = append(
			controllers[name].Actions,
			route.Action,
		)
	}

	for _, controller := range controllers {

		content, err := render(
			"controller.go.tmpl",
			controller,
		)

		if err != nil {
			panic(err)
		}

		filename := filepath.Join(
			"app",
			"Controllers",
			strings.ToLower(controller.Controller)+".go",
		)

		g.project.AddFile(filename, content)
	}
}