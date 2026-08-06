package server

import (
	"fmt"
	"net/http"

	"github.com/shishir1290/goduck/runtime/context"
	"github.com/shishir1290/goduck/runtime/engine"
	"github.com/shishir1290/goduck/runtime/handler"
	"github.com/shishir1290/goduck/runtime/middleware"
	"github.com/shishir1290/goduck/runtime/static"
)

type Server struct {

	engine *engine.Engine

	port int
}

func New(port int) *Server {
	return &Server{
		engine: engine.New(),
		port:   port,
	}
}

func (s *Server) GET(
	path string,
	h handler.HandlerFunc,
) {
	s.engine.Router.GET(path, h)
}

func (s *Server) POST(
	path string,
	handler handler.HandlerFunc,
) {
	s.engine.Router.POST(path, handler)
}

func (s *Server) PUT(
	path string,
	handler handler.HandlerFunc,
) {
	s.engine.Router.PUT(path, handler)
}

func (s *Server) PATCH(
	path string,
	handler handler.HandlerFunc,
) {
	s.engine.Router.PATCH(path, handler)
}

func (s *Server) DELETE(
	path string,
	handler handler.HandlerFunc,
) {
	s.engine.Router.DELETE(path, handler)
}

func (s *Server) Run() error {

	fmt.Printf("Goduck running on :%d\n", s.port)

	return http.ListenAndServe(
		fmt.Sprintf(":%d", s.port),
		s,
	)
}


func (s *Server) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
) {

	route, params := s.engine.Router.Find(
		r.Method,
		r.URL.Path,
	)

	if route == nil {

		http.NotFound(w, r)

		return
	}

	ctx := context.New(w, r)

	for k, v := range params {
		ctx.SetParam(k, v)
	}

	all := append(
		s.engine.Middlewares,
		route.Middlewares...,
	)

	h := middleware.Apply(
		route.Handler,
		all,
	)

	h(ctx)
}

func (s *Server) Use(
	m middleware.Middleware,
) {
	s.engine.Middlewares = append(
		s.engine.Middlewares,
		m,
	)
}

func (s *Server) Static(
	prefix string,
	root string,
) {

	s.engine.Statics = append(
		s.engine.Statics,
		static.New(prefix, root),
	)
}