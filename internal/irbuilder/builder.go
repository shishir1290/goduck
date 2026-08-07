package irbuilder

import (
	"github.com/shishir1290/goduck/internal/ast"
	"github.com/shishir1290/goduck/internal/ir"
)

type Builder struct {
	program *ast.Program
}

func New(program *ast.Program) *Builder {
	return &Builder{
		program: program,
	}
}

func (b *Builder) Build() *ir.Project {

	project := &ir.Project{}

	if b.program.App != nil {
		project.App.Name = b.program.App.Name
	}

	if b.program.Server != nil {
		project.Server.Port = b.program.Server.Port
	}

	for _, module := range b.program.Modules {

		irModule := ir.Module{
			Name:       module.Name,
			Controller: module.Controller,
			Service:    module.Service,
			Repository: module.Repository,
		}

		for _, route := range module.Routes {

			irModule.Routes = append(
				irModule.Routes,
				ir.Route{
					Method: route.Method,
					Path:   route.Path,
					Action: route.Action,
				},
			)
		}

		project.Modules = append(
			project.Modules,
			irModule,
		)
	}

	return project
}