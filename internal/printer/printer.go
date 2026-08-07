package printer

import (
	"fmt"

	"github.com/shishir1290/goduck/internal/ast"
)

func Print(program *ast.Program) {

	fmt.Println("========== AST ==========")

	if program.App != nil {
		fmt.Printf(
			"App: %s\n",
			program.App.Name,
		)
	}

	if program.Server != nil {
		fmt.Println("Server:")
		fmt.Printf(
			"  Port: %d\n",
			program.Server.Port,
		)
	}

	for _, module := range program.Modules {

		fmt.Println("Module:")
		fmt.Printf(
			"  Name: %s\n",
			module.Name,
		)

		if module.Controller != "" {
			fmt.Printf(
				"  Controller: %s\n",
				module.Controller,
			)
		}

		if module.Service != "" {
			fmt.Printf(
				"  Service: %s\n",
				module.Service,
			)
		}

		if module.Repository != "" {
			fmt.Printf(
				"  Repository: %s\n",
				module.Repository,
			)
		}

		for _, route := range module.Routes {

			fmt.Println("  Route:")

			fmt.Printf(
				"    Path: %s\n",
				route.Path,
			)

			fmt.Printf(
				"    Method: %s\n",
				route.Method,
			)

			fmt.Printf(
				"    Action: %s\n",
				route.Action,
			)
		}
	}

	fmt.Println("=========================")
}