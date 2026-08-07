package parser

import (
	"github.com/shishir1290/goduck/internal/ast"
	"github.com/shishir1290/goduck/internal/lexer"
)

func (p *Parser) ParseProgram() *ast.Program {

	program := &ast.Program{}

	for p.current.Type != lexer.EOF {

		switch p.current.Type {

		case lexer.APP:
			program.App = p.parseApp()

		case lexer.SERVER:
			program.Server = p.parseServer()

		case lexer.MODULE:

			module := p.parseModule()

			if module != nil {
				program.Modules = append(
					program.Modules,
					module,
				)
			}

		default:
			p.nextToken()
		}
	}

	return program
}