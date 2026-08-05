package semantic

import (
	"strings"

	"github.com/shishir1290/goduck/internal/ast"
)

func (a *Analyzer) checkRoute(route *ast.Route) {

	if strings.TrimSpace(route.Path) == "" {
		a.addError("route path cannot be empty")
	}

	switch strings.ToUpper(route.Method) {

	case "GET", "POST":

	default:
		a.addError("unsupported HTTP method: " + route.Method)
	}

	if route.Controller == "" {
		a.addError("controller name cannot be empty")
	}

	if route.Action == "" {
		a.addError("controller action cannot be empty")
	}
}