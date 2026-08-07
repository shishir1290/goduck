package ast

type Program struct {
	App       *App
	Classes   []*ClassDeclaration
	Variables []*VariableDeclaration
}

func (*Program) node() {}
