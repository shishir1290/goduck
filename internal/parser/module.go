package parser

import (
	"github.com/shishir1290/goduck/internal/ast"
	"github.com/shishir1290/goduck/internal/lexer"
)

func (p *Parser) parseModule() *ast.Module {

	module := &ast.Module{}

	// module name
	p.nextToken()

	if p.current.Type != lexer.IDENTIFIER {
		p.addError("expected module name")
		return nil
	}

	module.Name = p.current.Literal

	// {
	p.nextToken()

	if p.current.Type != lexer.LBRACE {
		p.addError("expected '{' after module name")
		return nil
	}

	p.nextToken()

	for p.current.Type != lexer.RBRACE &&
		p.current.Type != lexer.EOF {

		switch p.current.Type {

		case lexer.CONTROLLER:

			p.nextToken()

			if p.current.Type != lexer.IDENTIFIER {
				p.addError("expected controller name")
				return nil
			}

			module.Controller = p.current.Literal
			p.nextToken()

		case lexer.SERVICE:

			p.nextToken()

			if p.current.Type != lexer.IDENTIFIER {
				p.addError("expected service name")
				return nil
			}

			module.Service = p.current.Literal
			p.nextToken()

		case lexer.REPOSITORY:

			p.nextToken()

			if p.current.Type != lexer.IDENTIFIER {
				p.addError("expected repository name")
				return nil
			}

			module.Repository = p.current.Literal
			p.nextToken()

		case lexer.GET,
			lexer.POST,
			lexer.PUT,
			lexer.PATCH,
			lexer.DELETE:

			route := p.parseRoute()

			if route != nil {
				module.Routes = append(
					module.Routes,
					route,
				)
			}

		default:

			p.nextToken()
		}
	}

	if p.current.Type == lexer.RBRACE {
		p.nextToken()
	}

	return module
}