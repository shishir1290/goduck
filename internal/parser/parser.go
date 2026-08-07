package parser

import (
	"fmt"

	"github.com/shishir1290/goduck/internal/ast"
	"github.com/shishir1290/goduck/internal/lexer"
)

type Parser struct {
	lexer *lexer.Lexer

	current lexer.Token
	peek    lexer.Token

	errors []string

	prefixParseFns map[lexer.TokenType]func() ast.Expression
	infixParseFns  map[lexer.TokenType]func(ast.Expression) ast.Expression
}

const (
	_ int = iota
	LOWEST
	SUM     // +
	CALL    // ()
	MEMBER  // .
)

var precedences = map[lexer.TokenType]int{
	lexer.PLUS:   SUM,
	lexer.LPAREN: CALL,
	lexer.DOT:    MEMBER,
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		lexer:  l,
		errors: []string{},
	}

	p.prefixParseFns = make(map[lexer.TokenType]func() ast.Expression)
	p.prefixParseFns[lexer.IDENTIFIER] = p.parseIdentifier
	p.prefixParseFns[lexer.NUMBER] = p.parseNumberLiteral
	p.prefixParseFns[lexer.STRING] = p.parseStringLiteral
	p.prefixParseFns[lexer.LPAREN] = p.parseGroupedExpression

	p.infixParseFns = make(map[lexer.TokenType]func(ast.Expression) ast.Expression)
	p.infixParseFns[lexer.PLUS] = p.parseBinaryExpression
	p.infixParseFns[lexer.LPAREN] = p.parseCallExpression
	p.infixParseFns[lexer.DOT] = p.parseMemberExpression

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.current = p.peek
	p.peek = p.lexer.NextToken()
}

func (p *Parser) peekPrecedence() int {
	if prec, ok := precedences[p.peek.Type]; ok {
		return prec
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if prec, ok := precedences[p.current.Type]; ok {
		return prec
	}
	return LOWEST
}

func (p *Parser) expectPeek(t lexer.TokenType) bool {
	if p.peek.Type == t {
		p.nextToken()
		return true
	}
	p.addError(fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peek.Type))
	return false
}

// Top-level parser
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{
		Classes:   []*ast.ClassDeclaration{},
		Variables: []*ast.VariableDeclaration{},
	}

	var decorators []*ast.Decorator

	for p.current.Type != lexer.EOF {
		switch p.current.Type {
		case lexer.APP:
			program.App = p.parseApp()
		case lexer.IMPORT:
			p.parseImportStatement()
		case lexer.DECORATOR:
			dec := p.parseDecorator()
			if dec != nil {
				decorators = append(decorators, dec)
			}
		case lexer.CLASS, lexer.FUNC, lexer.SERVICE:
			classDec := p.parseClass(decorators)
			decorators = nil // Reset decorators after consuming
			if classDec != nil {
				program.Classes = append(program.Classes, classDec)
			}
		case lexer.IDENTIFIER:
			if p.current.Literal == "service" {
				classDec := p.parseClass(decorators)
				decorators = nil
				if classDec != nil {
					program.Classes = append(program.Classes, classDec)
				}
			} else {
				p.nextToken()
			}
		case lexer.LET, lexer.CONST:
			varDec := p.parseVariableDeclaration()
			if varDec != nil {
				program.Variables = append(program.Variables, varDec)
			}
		default:
			p.nextToken()
		}
	}

	return program
}

func (p *Parser) parseImportStatement() {
	// Consume "import"
	p.nextToken()

	// Consume until we find a semicolon or EOF
	for p.current.Type != lexer.SEMICOLON && p.current.Type != lexer.EOF {
		p.nextToken()
	}

	if p.current.Type == lexer.SEMICOLON {
		p.nextToken()
	}
}

func (p *Parser) parseApp() *ast.App {
	// app
	p.nextToken()

	if p.current.Type != lexer.IDENTIFIER {
		p.addError("expected app name")
		return nil
	}

	appNode := &ast.App{
		Name:       p.current.Literal,
		Properties: []*ast.FieldDeclaration{},
	}

	// {
	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}
	p.nextToken()

	// Parse fields inside App
	for p.current.Type != lexer.RBRACE && p.current.Type != lexer.EOF {
		if p.current.Type == lexer.IDENTIFIER || p.current.Type == lexer.PORT {
			field := p.parseFieldDeclaration()
			if field != nil {
				appNode.Properties = append(appNode.Properties, field)
			}
		} else {
			p.nextToken()
		}
	}

	if p.current.Type == lexer.RBRACE {
		p.nextToken()
	}

	return appNode
}

func (p *Parser) parseDecorator() *ast.Decorator {
	dec := &ast.Decorator{
		Name:      p.current.Literal,
		Arguments: []ast.Expression{},
	}

	p.nextToken() // move past decorator name

	if p.current.Type == lexer.LPAREN {
		p.nextToken() // past (
		for p.current.Type != lexer.RPAREN && p.current.Type != lexer.EOF {
			expr := p.parseExpression(LOWEST)
			if expr != nil {
				dec.Arguments = append(dec.Arguments, expr)
			}
			if p.peek.Type == lexer.COMMA {
				p.nextToken() // current is ,
				p.nextToken() // current is start of next arg
			} else if p.peek.Type == lexer.RPAREN {
				p.nextToken() // current is )
			} else {
				p.nextToken()
			}
		}
		if p.current.Type == lexer.RPAREN {
			p.nextToken() // past )
		}
	}

	return dec
}

func (p *Parser) parseClass(decorators []*ast.Decorator) *ast.ClassDeclaration {
	keywordType := p.current.Type
	keywordLiteral := p.current.Literal
	p.nextToken()

	if p.current.Type != lexer.IDENTIFIER {
		p.addError("expected name after class/func/service keyword")
		return nil
	}

	if keywordType == lexer.SERVICE || keywordLiteral == "service" {
		hasServiceDec := false
		for _, dec := range decorators {
			if dec.Name == "@Service" {
				hasServiceDec = true
			}
		}
		if !hasServiceDec {
			decorators = append(decorators, &ast.Decorator{
				Name:      "@Service",
				Arguments: []ast.Expression{},
			})
		}
	}

	classDec := &ast.ClassDeclaration{
		Decorators: decorators,
		Name:       p.current.Literal,
		Fields:     []*ast.FieldDeclaration{},
		Methods:    []*ast.MethodDeclaration{},
	}

	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}
	p.nextToken() // past {

	var memberDecorators []*ast.Decorator

	for p.current.Type != lexer.RBRACE && p.current.Type != lexer.EOF {
		switch p.current.Type {
		case lexer.DECORATOR:
			dec := p.parseDecorator()
			if dec != nil {
				memberDecorators = append(memberDecorators, dec)
			}
		case lexer.ASYNC, lexer.IDENTIFIER:
			isAsync := false
			if p.current.Type == lexer.ASYNC {
				isAsync = true
				p.nextToken()
			}

			if p.current.Type != lexer.IDENTIFIER {
				p.addError("expected member name")
				p.nextToken()
				continue
			}

			memberName := p.current.Literal

			// Peek ahead to see if it's a field or method
			if p.peek.Type == lexer.COLON {
				p.nextToken() // current is identifier
				field := p.parseFieldDeclarationRest(memberName)
				if field != nil {
					classDec.Fields = append(classDec.Fields, field)
				}
				memberDecorators = nil
			} else if p.peek.Type == lexer.LPAREN {
				p.nextToken() // current is LPAREN '('
				method := p.parseMethodDeclarationRest(memberDecorators, isAsync, memberName)
				if method != nil {
					classDec.Methods = append(classDec.Methods, method)
				}
				memberDecorators = nil
			} else {
				p.addError("unexpected token in class member: " + p.peek.Literal)
				p.nextToken()
			}
		default:
			p.nextToken()
		}
	}

	if p.current.Type == lexer.RBRACE {
		p.nextToken()
	}

	return classDec
}

func (p *Parser) parseFieldDeclaration() *ast.FieldDeclaration {
	name := p.current.Literal
	p.nextToken()
	return p.parseFieldDeclarationRest(name)
}

func (p *Parser) parseFieldDeclarationRest(name string) *ast.FieldDeclaration {
	if p.current.Type != lexer.COLON {
		p.addError("expected ':' after field name")
		return nil
	}
	p.nextToken() // past :

	typeNode := p.parseTypeNode()
	field := &ast.FieldDeclaration{
		Name: name,
		Type: typeNode,
	}

	if p.current.Type == lexer.ASSIGN {
		p.nextToken() // past =
		field.Value = p.parseExpression(LOWEST)
	}

	if p.current.Type == lexer.SEMICOLON {
		p.nextToken()
	}

	return field
}

func (p *Parser) parseMethodDeclarationRest(decorators []*ast.Decorator, isAsync bool, name string) *ast.MethodDeclaration {
	method := &ast.MethodDeclaration{
		Decorators: decorators,
		IsAsync:    isAsync,
		Name:       name,
		Parameters: []*ast.Parameter{},
	}

	if p.current.Type != lexer.LPAREN {
		p.addError("expected '(' before parameters")
		return nil
	}
	p.nextToken() // past (

	for p.current.Type != lexer.RPAREN && p.current.Type != lexer.EOF {
		if p.current.Type != lexer.IDENTIFIER {
			p.addError("expected parameter name")
			return nil
		}
		paramName := p.current.Literal
		p.nextToken()

		if p.current.Type != lexer.COLON {
			p.addError("expected ':' after parameter name")
			return nil
		}
		p.nextToken() // past :

		paramType := p.parseTypeNode()
		method.Parameters = append(method.Parameters, &ast.Parameter{
			Name: paramName,
			Type: paramType,
		})

		if p.current.Type == lexer.COMMA {
			p.nextToken()
		}
	}

	if p.current.Type != lexer.RPAREN {
		p.addError("expected ')' after parameters")
		return nil
	}
	p.nextToken() // past )

	// Optional return type
	if p.current.Type == lexer.COLON {
		p.nextToken() // past :
		method.ReturnType = p.parseTypeNode()
	} else {
		method.ReturnType = &ast.TypeNode{Name: "void"}
	}

	if p.current.Type != lexer.LBRACE {
		p.addError("expected '{' for method body")
		return nil
	}

	method.Body = p.parseBlockStatement()

	return method
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Statements: []ast.Statement{}}
	p.nextToken() // past {

	for p.current.Type != lexer.RBRACE && p.current.Type != lexer.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
	}

	if p.current.Type == lexer.RBRACE {
		p.nextToken() // past }
	}

	return block
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.current.Type {
	case lexer.LET, lexer.CONST:
		return p.parseVariableDeclaration()
	case lexer.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseVariableDeclaration() *ast.VariableDeclaration {
	isConst := p.current.Type == lexer.CONST
	p.nextToken() // past let/const

	if p.current.Type != lexer.IDENTIFIER {
		p.addError("expected identifier in variable declaration")
		return nil
	}
	name := p.current.Literal
	p.nextToken()

	var typeNode *ast.TypeNode
	if p.current.Type == lexer.COLON {
		p.nextToken() // past :
		typeNode = p.parseTypeNode()
	}

	var val ast.Expression
	if p.current.Type == lexer.ASSIGN {
		p.nextToken() // past =
		val = p.parseExpression(LOWEST)
	}

	if p.peek.Type == lexer.SEMICOLON {
		p.nextToken()
	}
	p.nextToken()

	return &ast.VariableDeclaration{
		Name:  name,
		Type:  typeNode,
		Value: val,
		Const: isConst,
	}
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{}
	p.nextToken() // past return

	if p.current.Type != lexer.SEMICOLON && p.current.Type != lexer.RBRACE {
		stmt.Value = p.parseExpression(LOWEST)
	}

	if p.peek.Type == lexer.SEMICOLON {
		p.nextToken()
	}
	p.nextToken()

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{}
	stmt.Expression = p.parseExpression(LOWEST)

	if p.peek.Type == lexer.SEMICOLON {
		p.nextToken()
	}
	p.nextToken()

	return stmt
}

func (p *Parser) parseTypeNode() *ast.TypeNode {
	typeName := ""

	if p.current.Type == lexer.STRING_TYPE {
		typeName = "string"
		p.nextToken()
	} else if p.current.Type == lexer.NUMBER_TYPE {
		typeName = "number"
		p.nextToken()
	} else if p.current.Type == lexer.BOOLEAN_TYPE {
		typeName = "boolean"
		p.nextToken()
	} else if p.current.Type == lexer.VOID_TYPE {
		typeName = "void"
		p.nextToken()
	} else if p.current.Type == lexer.IDENTIFIER {
		typeName = p.current.Literal
		p.nextToken()
		// Support generic types like Promise<User>
		if p.current.Type == lexer.LT {
			typeName += "<"
			p.nextToken() // past <
			if p.current.Type == lexer.IDENTIFIER {
				typeName += p.current.Literal
				p.nextToken()
			}
			if p.current.Type == lexer.GT {
				typeName += ">"
				p.nextToken()
			}
		}
	} else {
		p.addError("expected type name, got: " + p.current.Literal)
		return &ast.TypeNode{Name: "any"}
	}

	isArray := false
	if p.current.Type == lexer.LBRACKET {
		p.nextToken() // past [
		if p.current.Type == lexer.RBRACKET {
			isArray = true
			p.nextToken() // past ]
		}
	}

	return &ast.TypeNode{
		Name:    typeName,
		IsArray: isArray,
	}
}

// Expression Parser (Pratt Parser)
func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.current.Type]
	if prefix == nil {
		p.addError("no prefix parse function for " + string(p.current.Type))
		p.nextToken() // prevent infinite loop
		return nil
	}
	leftExp := prefix()

	for p.peek.Type != lexer.SEMICOLON && p.peek.Type != lexer.EOF && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peek.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	ident := &ast.Identifier{Name: p.current.Literal}
	return ident
}

func (p *Parser) parseNumberLiteral() ast.Expression {
	lit := &ast.Literal{Value: p.current.Literal, Type: "number"}
	return lit
}

func (p *Parser) parseStringLiteral() ast.Expression {
	lit := &ast.Literal{Value: p.current.Literal, Type: "string"}
	return lit
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken() // past (
	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}
	return exp
}

func (p *Parser) parseBinaryExpression(left ast.Expression) ast.Expression {
	exp := &ast.BinaryExpression{
		Left:     left,
		Operator: p.current.Literal,
	}
	precedence := p.curPrecedence()
	p.nextToken()
	exp.Right = p.parseExpression(precedence)
	return exp
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{
		Function:  function,
		Arguments: []ast.Expression{},
	}
	p.nextToken() // past (

	for p.current.Type != lexer.RPAREN && p.current.Type != lexer.EOF {
		arg := p.parseExpression(LOWEST)
		if arg != nil {
			exp.Arguments = append(exp.Arguments, arg)
		}
		if p.peek.Type == lexer.COMMA {
			p.nextToken() // current is ,
			p.nextToken() // current is start of next arg
		} else if p.peek.Type == lexer.RPAREN {
			p.nextToken() // current is )
		} else {
			p.nextToken()
		}
	}

	return exp
}

func (p *Parser) parseMemberExpression(left ast.Expression) ast.Expression {
	p.nextToken() // past .
	if p.current.Type != lexer.IDENTIFIER {
		p.addError("expected identifier after '.' in member expression")
		return nil
	}
	exp := &ast.MemberExpression{
		Object:   left,
		Property: p.current.Literal,
	}
	return exp
}
