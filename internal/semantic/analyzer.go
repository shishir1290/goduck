package semantic

import "github.com/shishir1290/goduck/internal/ast"

type Analyzer struct {
	program *ast.Program
}

func New(program *ast.Program) *Analyzer {
	return &Analyzer{
		program: program,
	}
}

func (a *Analyzer) Analyze() error {

	if err := a.checkApp(); err != nil {
		return err
	}

	if err := a.checkServer(); err != nil {
		return err
	}

	if err := a.checkModules(); err != nil {
		return err
	}

	if err := a.checkRoutes(); err != nil {
		return err
	}

	return nil
}