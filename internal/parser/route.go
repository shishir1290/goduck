package parser

import (
	"github.com/shishir1290/goduck/internal/ast"
	"github.com/shishir1290/goduck/internal/lexer"
)

func (p *Parser) parseRoute() *ast.Route {

	route := &ast.Route{}

	p.nextToken()

	route.Path = p.current.Literal

	p.nextToken()

	if p.current.Type != lexer.LBRACE {
		p.addError("expected '{'")
		return route
	}

	p.nextToken()

	route.Method = p.current.Literal

	p.nextToken()

	if p.current.Type != lexer.ARROW {
		p.addError("expected ->")
		return route
	}

	p.nextToken()

	route.Controller = p.current.Literal

	p.nextToken()

	if p.current.Type != lexer.DOT {
		p.addError("expected '.'")
		return route
	}

	p.nextToken()

	route.Action = p.current.Literal

	for p.current.Type != lexer.RBRACE &&
		p.current.Type != lexer.EOF {

		p.nextToken()
	}

	p.nextToken()

	return route
}