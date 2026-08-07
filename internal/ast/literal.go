package ast

type Literal struct {
	Value string
	Type  string
}

func (*Literal) node() {}
