package router

import "github.com/shishir1290/goduck/runtime/handler"

type Group struct {
	prefix string
	router *Router
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

	g.router.GET(
		g.prefix+path,
		handler,
	)
}

func (g *Group) POST(
	path string,
	handler handler.HandlerFunc,
) {
	g.router.POST(
		g.prefix+path,
		handler,
	)
}

func (g *Group) PUT(
	path string,
	handler handler.HandlerFunc,
) {
	g.router.PUT(
		g.prefix+path,
		handler,
	)
}

func (g *Group) DELETE(
	path string,
	handler handler.HandlerFunc,
) {
	g.router.DELETE(
		g.prefix+path,
		handler,
	)
}

func (g *Group) PATCH(
	path string,
	handler handler.HandlerFunc,
) {
	g.router.PATCH(
		g.prefix+path,
		handler,
	)
}