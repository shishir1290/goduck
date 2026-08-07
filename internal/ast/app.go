package ast

type App struct {
	Name       string
	Properties []*FieldDeclaration
}

func (*App) node() {}
