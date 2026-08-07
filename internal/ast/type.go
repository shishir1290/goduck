package ast

type TypeNode struct {
	Name    string
	IsArray bool
}

func (*TypeNode) node() {}
