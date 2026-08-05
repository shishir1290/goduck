package ast

type Program struct {
	App    *App
	Server *Server
	Routes []*Route
}

func (*Program) node() {}