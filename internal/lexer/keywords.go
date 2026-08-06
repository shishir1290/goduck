package lexer

import "strings"

var keywords = map[string]TokenType{
	"app": APP,

	"server": SERVER,

	// "route": ROUTE,

	"get": GET,

	"post": POST,

	"port": PORT,

	"put": PUT,

	"patch": PATCH,

	"delete": DELETE,
}

func LookupIdentifier(ident string) TokenType {

	if tok, ok := keywords[strings.ToLower(ident)]; ok {
		return tok
	}

	return IDENTIFIER
}