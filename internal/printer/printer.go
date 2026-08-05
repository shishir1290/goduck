package printer

import (
	"fmt"

	"github.com/shishir1290/goduck/internal/ast"
)

func Print(program *ast.Program) {

	fmt.Println("========== AST ==========")

	if program.App != nil {
		fmt.Printf("App: %s\n", program.App.Name)
	}

	if program.Server != nil {
		fmt.Printf("Server:\n")
		fmt.Printf("  Port: %d\n", program.Server.Port)
	}

	for _, route := range program.Routes {

		fmt.Println("Route:")
		fmt.Printf("  Path: %s\n", route.Path)
		fmt.Printf("  Method: %s\n", route.Method)
		fmt.Printf("  Controller: %s\n", route.Controller)
		fmt.Printf("  Action: %s\n", route.Action)
	}

	fmt.Println("=========================")
}