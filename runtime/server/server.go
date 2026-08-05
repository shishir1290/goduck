package server

import (
	"fmt"
	"net/http"

	"github.com/shishir1290/goduck/runtime/context"
	"github.com/shishir1290/goduck/runtime/handler"
	"github.com/shishir1290/goduck/runtime/middleware"
	"github.com/shishir1290/goduck/runtime/router"
)

type Server struct {
	router      *router.Router
	port        int
	middlewares []middleware.Middleware
}

func New(port int) *Server {
	return &Server{
		router: router.New(),
		port:   port,
	}
}

func (s *Server) GET(
	path string,
	h handler.HandlerFunc,
) {
	s.router.GET(path, h)
}

func (s *Server) POST(
	path string,
	handler handler.HandlerFunc,
) {
	s.router.POST(path, handler)
}

func (s *Server) PUT(
	path string,
	handler handler.HandlerFunc,
) {
	s.router.PUT(path, handler)
}

func (s *Server) DELETE(
	path string,
	handler handler.HandlerFunc,
) {
	s.router.DELETE(path, handler)
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

	route, params := s.router.Find(
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

	h := middleware.Apply(
		route.Handler,
		s.middlewares,
	)

	h(ctx)
}

func (s *Server) Use(
	m middleware.Middleware,
) {
	s.middlewares = append(
		s.middlewares,
		m,
	)
}