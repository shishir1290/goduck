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

	project.App.Name = b.program.App.Name

	project.Server.Port = b.program.Server.Port

	for _, route := range b.program.Routes {

		project.Routes = append(project.Routes, ir.Route{

			Method: route.Method,

			Path: route.Path,

			Controller: route.Controller,

			Action: route.Action,
		})
	}

	return project
}