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
		project.AppName = b.program.App.Name
	}
	project.Classes = b.program.Classes
	return project
}