package ast

type Server struct {
	Port int
}

func (*Server) node() {}