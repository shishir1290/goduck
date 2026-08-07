package compiler

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/shishir1290/goduck/internal/ast"
	"github.com/shishir1290/goduck/internal/generator"
	"github.com/shishir1290/goduck/internal/gobuilder"
	"github.com/shishir1290/goduck/internal/irbuilder"
	"github.com/shishir1290/goduck/internal/lexer"
	"github.com/shishir1290/goduck/internal/parser"
	"github.com/shishir1290/goduck/internal/printer"
	"github.com/shishir1290/goduck/internal/semantic"
	"github.com/shishir1290/goduck/internal/writer"
)

func Build(filename string) (string, error) {
	var duckFiles []string

	fi, err := os.Stat(filename)
	if err != nil {
		return "", err
	}

	if fi.IsDir() {
		err = filepath.Walk(filename, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && strings.HasSuffix(info.Name(), ".duck") {
				duckFiles = append(duckFiles, path)
			}
			return nil
		})
		if err != nil {
			return "", err
		}
	} else {
		dir := filepath.Dir(filename)
		err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && strings.HasSuffix(info.Name(), ".duck") {
				duckFiles = append(duckFiles, path)
			}
			return nil
		})
		if err != nil {
			return "", err
		}
	}

	if len(duckFiles) == 0 {
		return "", fmt.Errorf("no .duck files found to compile")
	}

	program := &ast.Program{
		Classes:   []*ast.ClassDeclaration{},
		Variables: []*ast.VariableDeclaration{},
	}

	hasErrors := false
	for _, file := range duckFiles {
		fmt.Printf("Lexing & Parsing %s...\n", file)
		content, err := os.ReadFile(file)
		if err != nil {
			return "", err
		}

		l := lexer.New(string(content))
		p := parser.New(l)
		prog := p.ParseProgram()

		if len(p.Errors()) > 0 {
			fmt.Printf("Parser Errors in %s:\n", file)
			for _, pErr := range p.Errors() {
				fmt.Printf("  - %s\n", pErr)
			}
			hasErrors = true
		}

		if prog != nil {
			if prog.App != nil {
				program.App = prog.App
			}
			program.Classes = append(program.Classes, prog.Classes...)
			program.Variables = append(program.Variables, prog.Variables...)
		}
	}

	if hasErrors {
		return "", fmt.Errorf("parser failed")
	}

	fmt.Println("✓ Parsing completed")

	// Phase 3 - Semantic Analysis
	analyzer := semantic.New(program)
	errors := analyzer.Analyze()

	if errors != nil {
		fmt.Println("\nSemantic / Type Check Errors:")
		fmt.Println("-", errors)
		return "", fmt.Errorf("semantic analysis failed")
	}

	fmt.Println("✓ Semantic / Type checking completed")

	builder := irbuilder.New(program)
	projectIR := builder.Build()

	fmt.Println()
	fmt.Println("IR")
	fmt.Printf("%+v\n", projectIR)

	gen := generator.New(program)
	project := gen.Generate()

	w := writer.New(project)
	if err := w.Write(); err != nil {
		return "", err
	}

	fmt.Println("✓ Project written")

	gb := gobuilder.New(project)
	if err := gb.Build(); err != nil {
		return "", err
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
	return project.Name, nil
}