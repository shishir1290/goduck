package lexer

var keywords = map[string]TokenType{

	"app": APP,

	"server": SERVER,

	"module": MODULE,

	"get":    GET,
	"post":   POST,
	"put":    PUT,
	"patch":  PATCH,
	"delete": DELETE,

	"GET":    GET,
	"POST":   POST,
	"PUT":    PUT,
	"PATCH":  PATCH,
	"DELETE": DELETE,

	"port": PORT,

	"controller": CONTROLLER,
	"service":    SERVICE,
	"repository": REPOSITORY,
}

func LookupIdentifier(
	ident string,
) TokenType {

	if tok, ok := keywords[ident]; ok {
		return tok
	}

	return IDENTIFIER
}