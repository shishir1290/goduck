package ast

type Decorator struct {
	Name      string
	Arguments []Expression
}

func (*Decorator) node() {}
