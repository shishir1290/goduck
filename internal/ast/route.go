package ast

type Route struct {
	Path   string
	Method string
	Action string
}

func (*Route) node() {}