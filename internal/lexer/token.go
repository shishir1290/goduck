package lexer

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (

	// Special
	ILLEGAL TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"

	// Identifier
	IDENTIFIER = "IDENTIFIER"

	// Literals
	NUMBER = "NUMBER"
	STRING = "STRING"

	// Symbols
	LBRACE = "{"
	RBRACE = "}"
	LPAREN = "("
	RPAREN = ")"

	DOT   = "."
	COMMA = ","

	ARROW = "->"

	// Keywords
	APP = "APP"

	SERVER = "SERVER"

	// ROUTE = "ROUTE"

	GET = "GET"

	POST = "POST"

	PORT = "PORT"

	PUT = "PUT"

	PATCH = "PATCH"

	DELETE = "DELETE"
)