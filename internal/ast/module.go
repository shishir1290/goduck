package ast

import "strings"

type Module struct {
	Name string

	Controller string

	Service string

	Repository string

	Routes []*Route
}

func (*Module) node() {}

func (m *Module) PackageName() string {
	return strings.ToLower(m.Name)
}