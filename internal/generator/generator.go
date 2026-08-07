package generator

import "github.com/shishir1290/goduck/internal/ast"

type Generator struct {
	program *ast.Program
	project *Project
}

func New(program *ast.Program) *Generator {

	project := &Project{
		Name: program.App.Name,
	}

	return &Generator{
		program: program,
		project: project,
	}
}

func (g *Generator) Generate() *Project {

	g.generateGoMod()

	g.generateMain()

	g.generateModules()

	g.generateRouter()

	return g.project
}