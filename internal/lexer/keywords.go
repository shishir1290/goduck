package lexer

var keywords = map[string]TokenType{
	"app": APP,

	"server": SERVER,

	"route": ROUTE,

	"get": GET,

	"post": POST,

	"port": PORT,
}

func LookupIdentifier(ident string) TokenType {

	if tok, ok := keywords[ident]; ok {
		return tok
	}

	return IDENTIFIER
}