package parser

import "github.com/shishir1290/goduck/internal/ast"

func (p *Parser) parseApp() *ast.App {

	node := &ast.App{}

	p.nextToken()

	node.Name = p.current.Literal

	p.nextToken()

	return node
}