package parser

import (
	"strconv"

	"github.com/shishir1290/goduck/internal/ast"
	"github.com/shishir1290/goduck/internal/lexer"
)

func (p *Parser) parseServer() *ast.Server {

	server := &ast.Server{}

	p.nextToken()

	if p.current.Type != lexer.LBRACE {
		p.addError("expected '{' after server")
		return server
	}

	p.nextToken()

	for p.current.Type != lexer.RBRACE &&
		p.current.Type != lexer.EOF {

		if p.current.Type == lexer.PORT {

			p.nextToken()

			port, _ := strconv.Atoi(p.current.Literal)

			server.Port = port
		}

		p.nextToken()
	}

	p.nextToken()

	return server
}