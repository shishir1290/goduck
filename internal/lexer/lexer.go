package lexer

type Lexer struct {
	input string

	position     int
	readPosition int
	ch           byte
}

func New(input string) *Lexer {

	l := &Lexer{
		input: input,
	}

	l.readChar()

	return l
}

func (l *Lexer) readChar() {

	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) peekChar() byte {

	if l.readPosition >= len(l.input) {
		return 0
	}

	return l.input[l.readPosition]
}

func (l *Lexer) skipWhitespace() {

	for l.ch == ' ' ||
		l.ch == '\n' ||
		l.ch == '\t' ||
		l.ch == '\r' {

		l.readChar()
	}
}

func isLetter(ch byte) bool {

	return ('a' <= ch && ch <= 'z') ||
		('A' <= ch && ch <= 'Z') ||
		ch == '_'
}

func isDigit(ch byte) bool {

	return '0' <= ch && ch <= '9'
}

func (l *Lexer) readIdentifier() string {

	position := l.position

	for isLetter(l.ch) {
		l.readChar()
	}

	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {

	position := l.position

	for isDigit(l.ch) {
		l.readChar()
	}

	return l.input[position:l.position]
}

func (l *Lexer) readString() string {

	l.readChar()

	position := l.position

	for l.ch != '"' && l.ch != 0 {
		l.readChar()
	}

	return l.input[position:l.position]
}

func newToken(tokenType TokenType, ch byte) Token {
	return Token{
		Type:    tokenType,
		Literal: string(ch),
	}
}

func (l *Lexer) NextToken() Token {

	var tok Token

	l.skipWhitespace()

	switch l.ch {

	case '{':
		tok = newToken(LBRACE, l.ch)

	case '}':
		tok = newToken(RBRACE, l.ch)

	case '(':
		tok = newToken(LPAREN, l.ch)

	case ')':
		tok = newToken(RPAREN, l.ch)

	case '.':
		tok = newToken(DOT, l.ch)

	case ',':
		tok = newToken(COMMA, l.ch)

	case '"':
		tok.Type = STRING
		tok.Literal = l.readString()

	case '-':
		if l.peekChar() == '>' {
			ch := l.ch
			l.readChar()
			tok = Token{
				Type:    ARROW,
				Literal: string(ch) + string(l.ch),
			}
		} else {
			tok = newToken(ILLEGAL, l.ch)
		}

	case 0:
		tok = Token{
			Type:    EOF,
			Literal: "",
		}

	default:

		if isLetter(l.ch) {

			literal := l.readIdentifier()

			tok.Type = LookupIdentifier(literal)
			tok.Literal = literal

			return tok
		}

		if isDigit(l.ch) {

			tok.Type = NUMBER
			tok.Literal = l.readNumber()

			return tok
		}

		tok = newToken(ILLEGAL, l.ch)
	}

	l.readChar()

	return tok
}