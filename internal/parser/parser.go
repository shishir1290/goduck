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

func (p *Parser) isIdentifier(t lexer.TokenType) bool {
	return t == lexer.IDENTIFIER ||
		t == lexer.APP ||
		t == lexer.SERVER ||
		t == lexer.GET ||
		t == lexer.POST ||
		t == lexer.PORT ||
		t == lexer.PUT ||
		t == lexer.PATCH ||
		t == lexer.DELETE
}