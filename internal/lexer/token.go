package lexer

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"

	IDENTIFIER TokenType = "IDENTIFIER"

	NUMBER TokenType = "NUMBER"
	STRING TokenType = "STRING"

	LBRACE = "{"
	RBRACE = "}"

	LPAREN = "("
	RPAREN = ")"

	DOT   = "."
	COMMA = ","

	ARROW = "->"

	APP    = "APP"
	SERVER = "SERVER"
	MODULE = "MODULE"

	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	PATCH  = "PATCH"
	DELETE = "DELETE"

	PORT = "PORT"

	CONTROLLER = "CONTROLLER"
	SERVICE    = "SERVICE"
	REPOSITORY = "REPOSITORY"
)