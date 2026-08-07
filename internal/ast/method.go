package ast

type MethodDeclaration struct {
	Decorators []*Decorator
	IsAsync    bool
	Name       string
	Parameters []*Parameter
	ReturnType *TypeNode
	Body       *BlockStatement
}

func (*MethodDeclaration) node() {}

type Parameter struct {
	Name string
	Type *TypeNode
}

func (*Parameter) node() {}
