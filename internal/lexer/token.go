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

	LBRACKET = "["
	RBRACKET = "]"

	LT = "<"
	GT = ">"

	DOT   = "."
	COMMA = ","

	ARROW = "->"

	APP    = "APP"
	SERVER = "SERVER"
	ROUTE  = "ROUTE"

	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	PATCH  = "PATCH"
	DELETE = "DELETE"

	PORT = "PORT"

	CONTROLLER = "CONTROLLER"
	SERVICE    = "SERVICE"
	REPOSITORY = "REPOSITORY"

	// Variable
	LET   = "LET"
	CONST = "CONST"

	IMPORT = "IMPORT"
	FROM   = "FROM"

	// Types
	STRING_TYPE  = "STRING_TYPE"
	NUMBER_TYPE  = "NUMBER_TYPE"
	BOOLEAN_TYPE = "BOOLEAN_TYPE"
	VOID_TYPE    = "VOID_TYPE"
	// OOP keywords
	CLASS     = "CLASS"
	FUNC      = "FUNC"
	ASYNC     = "ASYNC"
	AWAIT     = "AWAIT"
	RETURN    = "RETURN"
	DECORATOR = "DECORATOR"
	PLUS      = "+"

	// Assignment
	ASSIGN = "="

	// Statement
	SEMICOLON = ";"
	COLON     = ":"
)
