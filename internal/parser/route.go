package parser

import (
	"strings"

	"github.com/shishir1290/goduck/internal/ast"
	"github.com/shishir1290/goduck/internal/lexer"
)

func (p *Parser) parseRoute() *ast.Route {

	route := &ast.Route{}

	// HTTP method
	route.Method = strings.ToUpper(
		p.current.Literal,
	)

	// Path
	p.nextToken()

	if p.current.Type != lexer.STRING {
		p.addError("expected route path")
		return nil
	}

	route.Path = p.current.Literal

	// ->
	p.nextToken()

	if p.current.Type != lexer.ARROW {
		p.addError("expected '->'")
		return nil
	}

	// Action
	p.nextToken()

	if p.current.Type != lexer.IDENTIFIER {
		p.addError("expected action")
		return nil
	}

	route.Action = p.current.Literal

	// Move to next token
	p.nextToken()

	return route
}