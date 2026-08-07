package lexer

var keywords = map[string]TokenType{
	"app":    APP,
	"server": SERVER,
	"route":  ROUTE,

	"get":    GET,
	"post":   POST,
	"put":    PUT,
	"patch":  PATCH,
	"delete": DELETE,

	"port": PORT,

	"let":   LET,
	"const": CONST,

	"string":  STRING_TYPE,
	"number":  NUMBER_TYPE,
	"boolean": BOOLEAN_TYPE,
	"void":    VOID_TYPE,

	"class":   CLASS,
	"func":    FUNC,
	"async":   ASYNC,
	"await":   AWAIT,
	"return":  RETURN,

	"import": IMPORT,
	"from":   FROM,
}

func LookupIdentifier(ident string) TokenType {

	if tok, ok := keywords[ident]; ok {
		return tok
	}

	return IDENTIFIER
}