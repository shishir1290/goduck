package parser

import (
	"strings"

	"github.com/shishir1290/goduck/internal/ast"
	"github.com/shishir1290/goduck/internal/lexer"
)

func (p *Parser) parseRoute() *ast.Route {

	route := &ast.Route{}

	// Current token is GET/POST/PUT/PATCH/DELETE
	route.Method = strings.ToLower(p.current.Literal)

	// Move to path
	p.nextToken()

	if p.current.Type != lexer.STRING {
		p.addError("expected route path")
		return nil
	}

	route.Path = p.current.Literal

	// Move to ->
	p.nextToken()

	if p.current.Type != lexer.ARROW {
		p.addError("expected '->'")
		return nil
	}

	// Move to controller
	p.nextToken()

	if !p.isIdentifier(p.current.Type) {
		p.addError("expected controller name")
		return nil
	}

	route.Controller = p.current.Literal

	// Move to .
	p.nextToken()

	if p.current.Type != lexer.DOT {
		p.addError("expected '.'")
		return nil
	}

	// Move to action
	p.nextToken()

	if !p.isIdentifier(p.current.Type) {
		p.addError("expected action name")
		return nil
	}

	route.Action = p.current.Literal

	// Advance to the next token for ParseProgram()
	p.nextToken()

	return route
}