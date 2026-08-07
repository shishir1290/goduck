package ast

type Program struct {
	App     *App
	Server  *Server
	Modules []*Module
}

func (*Program) node() {}