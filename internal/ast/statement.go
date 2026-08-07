package ast

type Statement interface {
	Node
	statementNode()
}

type BlockStatement struct {
	Statements []Statement
}

func (*BlockStatement) node() {}
func (*BlockStatement) statementNode() {}

type ReturnStatement struct {
	Value Expression
}

func (*ReturnStatement) node() {}
func (*ReturnStatement) statementNode() {}

type VariableDeclaration struct {
	Name  string
	Type  *TypeNode
	Value Expression
	Const bool
}

func (*VariableDeclaration) node() {}
func (*VariableDeclaration) statementNode() {}

type ExpressionStatement struct {
	Expression Expression
}

func (*ExpressionStatement) node() {}
func (*ExpressionStatement) statementNode() {}
