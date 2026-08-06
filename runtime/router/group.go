package router

import (
	"github.com/shishir1290/goduck/runtime/handler"
	"github.com/shishir1290/goduck/runtime/middleware"
)

type Group struct {
	prefix string

	router *Router

	middlewares []middleware.Middleware
}

func (r *Router) Group(
	prefix string,
) *Group {

	return &Group{
		prefix: prefix,
		router: r,
	}
}

func (g *Group) GET(
	path string,
	handler handler.HandlerFunc,
) {

	g.router.AddRoute(&Route{

		Method: "GET",

		Path: g.prefix + path,

		Handler: handler,

		Middlewares: g.middlewares,
	})
}

func (g *Group) POST(
	path string,
	handler handler.HandlerFunc,
) {
	g.router.AddRoute(&Route{

		Method: "POST",

		Path: g.prefix + path,

		Handler: handler,

		Middlewares: g.middlewares,
	})
}

func (g *Group) PUT(
	path string,
	handler handler.HandlerFunc,
) {
	g.router.AddRoute(&Route{

		Method: "PUT",

		Path: g.prefix + path,

		Handler: handler,

		Middlewares: g.middlewares,
	})
}

func (g *Group) DELETE(
	path string,
	handler handler.HandlerFunc,
) {
	g.router.AddRoute(&Route{

		Method: "DELETE",

		Path: g.prefix + path,

		Handler: handler,

		Middlewares: g.middlewares,
	})
}

func (g *Group) PATCH(
	path string,
	handler handler.HandlerFunc,
) {
	g.router.AddRoute(&Route{

		Method: "PATCH",

		Path: g.prefix + path,

		Handler: handler,

		Middlewares: g.middlewares,
	})
}

func (g *Group) Use(
	m middleware.Middleware,
) {

	g.middlewares = append(
		g.middlewares,
		m,
	)
}