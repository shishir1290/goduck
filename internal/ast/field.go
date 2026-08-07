package ast

type FieldDeclaration struct {
	Name  string
	Type  *TypeNode
	Value Expression
}

func (*FieldDeclaration) node() {}
