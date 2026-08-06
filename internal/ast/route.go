package ast

import "strings"

type Route struct {
	Path       string
	Method     string
	Controller string
	Action     string
}

func (*Route) node() {}

func (r *Route) MethodUpper() string {
	return strings.ToUpper(r.Method)
}