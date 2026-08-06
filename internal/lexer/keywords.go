package lexer

var keywords = map[string]TokenType{
	"app": APP,

	"server": SERVER,

	// "route": ROUTE,

	"GET": GET,

	"POST": POST,

	"port": PORT,

	"PUT": PUT,

	"PATCH": PATCH,

	"DELETE": DELETE,
}

func LookupIdentifier(ident string) TokenType {

	if tok, ok := keywords[ident]; ok {
		return tok
	}

	return IDENTIFIER
}