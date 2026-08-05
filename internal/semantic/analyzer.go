package semantic

import (
	"github.com/shishir1290/goduck/internal/ast"
)

type Analyzer struct {
	program *ast.Program
	errors  []error
}

func New(program *ast.Program) *Analyzer {
	return &Analyzer{
		program: program,
		errors:  []error{},
	}
}

func (a *Analyzer) Analyze() []error {

	a.checkApp()
	a.checkServer()

	for _, route := range a.program.Routes {
		a.checkRoute(route)
	}

	return a.errors
}

func (a *Analyzer) addError(msg string) {
	a.errors = append(a.errors, Error{
		Message: msg,
	})
}