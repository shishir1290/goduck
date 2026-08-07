package ast

type ClassDeclaration struct {
	Decorators []*Decorator
	Name       string
	Fields     []*FieldDeclaration
	Methods    []*MethodDeclaration
}

func (*ClassDeclaration) node() {}
