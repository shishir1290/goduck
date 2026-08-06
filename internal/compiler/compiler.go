package compiler

import (
	"fmt"
	"os"

	"github.com/shishir1290/goduck/internal/generator"
	"github.com/shishir1290/goduck/internal/gobuilder"
	"github.com/shishir1290/goduck/internal/irbuilder"
	"github.com/shishir1290/goduck/internal/lexer"
	"github.com/shishir1290/goduck/internal/parser"
	"github.com/shishir1290/goduck/internal/printer"
	"github.com/shishir1290/goduck/internal/semantic"
	"github.com/shishir1290/goduck/internal/writer"
)

func Build(filename string) error {

	fmt.Printf("Reading %s...\n", filename)

	content, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	// Phase 1 - Lexer
	l := lexer.New(string(content))
	fmt.Println("✓ Lexing completed")

	// Phase 2 - Parser
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		fmt.Println("\nParser Errors:")
		for _, err := range p.Errors() {
			fmt.Println("-", err)
		}
		return fmt.Errorf("parser failed")
	}

	fmt.Println("✓ Parsing completed")

	// Phase 3 - Semantic Analysis
	analyzer := semantic.New(program)

	errors := analyzer.Analyze()

	if errors != nil {

		fmt.Println("\nSemantic Errors:")

		fmt.Println("-", errors)

		return fmt.Errorf("semantic analysis failed")
	}

	fmt.Println("✓ Semantic analysis completed")

	builder := irbuilder.New(program)

	projectIR := builder.Build()

	fmt.Println()

	fmt.Println("IR")

	fmt.Printf("%+v\n", projectIR)

	gen := generator.New(program)

	project := gen.Generate()

	w := writer.New(project)

	if err := w.Write(); err != nil {
		return err
	}

	fmt.Println("✓ Project written")

	gb := gobuilder.New(project)

	if err := gb.Build(); err != nil {
		return err
	}

	fmt.Println("✓ Go build completed")

	fmt.Println()

	fmt.Println("Generated Files:")

	for _, file := range project.Files {
		fmt.Println("-", file.Path)
	}

	// Debug
	printer.Print(program)

	fmt.Println("\n✓ Build succeeded.")

	return nil
}