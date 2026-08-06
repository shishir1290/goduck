package parser

import "github.com/shishir1290/goduck/internal/lexer"

type Parser struct {
    lexer *lexer.Lexer

    current lexer.Token
    peek    lexer.Token

    errors []string
}

func New(l *lexer.Lexer) *Parser {

	p := &Parser{
		lexer: l,
	}

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.current = p.peek
	p.peek = p.lexer.NextToken()
}