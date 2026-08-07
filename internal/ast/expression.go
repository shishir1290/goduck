package ast

type Expression interface {
	node()
}

type Identifier struct {
	Name string
}

func (*Identifier) node() {}

type BinaryExpression struct {
	Left     Expression
	Operator string
	Right    Expression
}

func (*BinaryExpression) node() {}

type CallExpression struct {
	Function  Expression
	Arguments []Expression
}

func (*CallExpression) node() {}

type MemberExpression struct {
	Object   Expression
	Property string
}

func (*MemberExpression) node() {}
