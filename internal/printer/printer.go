package printer

import (
	"fmt"

	"github.com/shishir1290/goduck/internal/ast"
)

func Print(program *ast.Program) {
	fmt.Println("========== AST ==========")
	if program.App != nil {
		fmt.Printf("App: %s\n", program.App.Name)
		for _, prop := range program.App.Properties {
			fmt.Printf("  Prop: %s\n", prop.Name)
		}
	}

	for _, class := range program.Classes {
		fmt.Printf("Class: %s\n", class.Name)
		for _, dec := range class.Decorators {
			fmt.Printf("  Decorator: %s\n", dec.Name)
		}
		for _, field := range class.Fields {
			fmt.Printf("  Field: %s (Type: %s)\n", field.Name, field.Type.Name)
		}
		for _, method := range class.Methods {
			fmt.Printf("  Method: %s\n", method.Name)
			for _, dec := range method.Decorators {
				fmt.Printf("    Decorator: %s\n", dec.Name)
			}
		}
	}
	fmt.Println("=========================")
}