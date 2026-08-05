package ast

type Route struct {
	Path       string
	Method     string
	Controller string
	Action     string
}

func (*Route) node() {}