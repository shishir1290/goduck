package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/shishir1290/goduck/runtime/container"
	"github.com/shishir1290/goduck/runtime/context"
	"github.com/shishir1290/goduck/runtime/database"
	"github.com/shishir1290/goduck/runtime/engine"
	"github.com/shishir1290/goduck/runtime/handler"
	"github.com/shishir1290/goduck/runtime/middleware"
	"github.com/shishir1290/goduck/runtime/router"
	"github.com/shishir1290/goduck/runtime/static"
)

type Server struct {

	engine *engine.Engine

	port int

	container *container.Container
}

func New(port int) *Server {

	return &Server{
		engine:    engine.New(),
		container: container.New(),
		port:      port,
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

func (s *Server) Group(
	prefix string,
) *router.Group {
	return s.engine.Router.Group(prefix)
}

func (s *Server) Run() error {

	fmt.Printf("Goduck running on :%d\n", s.port)

	if s.container.Has((*database.Database)(nil)) {
		fmt.Println("Database connection: Connected")
	}
	if s.container.HasTypeWithName("socket") {
		fmt.Println("Socket server: Active")
	}

	routes := s.engine.Router.Routes()
	if len(routes) > 0 {
		fmt.Println("\nMapped Routes:")
		for _, r := range routes {
			fmt.Printf("  %-8s http://localhost:%d%s\n", "["+r.Method+"]", s.port, r.Path)
		}
		fmt.Println()
	}

	return http.ListenAndServe(
		fmt.Sprintf(":%d", s.port),
		s,
	)
}


func (s *Server) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
) {
	// Enable CORS by default
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS, HEAD")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	route, params := s.engine.Router.Find(
		r.Method,
		r.URL.Path,
	)

	if route == nil {
		// Check static files if no route matched
		for _, st := range s.engine.Statics {
			prefix := st.Prefix
			if !strings.HasSuffix(prefix, "/") {
				prefix += "/"
			}
			path := r.URL.Path
			if !strings.HasSuffix(path, "/") {
				path += "/"
			}
			if strings.HasPrefix(path, prefix) {
				st.Handler().ServeHTTP(w, r)
				return
			}
		}

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

func (s *Server) Container() *container.Container {
	return s.container
}