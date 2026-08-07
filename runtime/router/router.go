package router

import "github.com/shishir1290/goduck/runtime/handler"

type Router struct {
	tree *Tree
}

func New() *Router {
	return &Router{
		tree: NewTree(),
	}
}

func (r *Router) GET(path string, h handler.HandlerFunc) {

	r.tree.Insert(&Route{
		Method:  "GET",
		Path:    path,
		Handler: h,
	})
}

func (r *Router) POST(path string, h handler.HandlerFunc) {

	r.tree.Insert(&Route{
		Method:  "POST",
		Path:    path,
		Handler: h,
	})
}

func (r *Router) PUT(path string, h handler.HandlerFunc) {

	r.tree.Insert(&Route{
		Method:  "PUT",
		Path:    path,
		Handler: h,
	})
}

func (r *Router) PATCH(path string, h handler.HandlerFunc) {

	r.tree.Insert(&Route{
		Method:  "PATCH",
		Path:    path,
		Handler: h,
	})
}

func (r *Router) DELETE(path string, h handler.HandlerFunc) {

	r.tree.Insert(&Route{
		Method:  "DELETE",
		Path:    path,
		Handler: h,
	})
}

func (r *Router) Find(
	method,
	path string,
) (*Route, map[string]string) {

	return r.tree.Find(method, path)
}

func (r *Router) AddRoute(route *Route) {

	r.tree.Insert(route)

}

func (r *Router) Routes() []*Route {
	var routes []*Route
	var walk func(n *Node)
	walk = func(n *Node) {
		if n == nil {
			return
		}
		for _, route := range n.Routes {
			routes = append(routes, route)
		}
		if n.Parameter != nil {
			walk(n.Parameter)
		}
		for _, child := range n.Children {
			walk(child)
		}
	}
	walk(r.tree.root)
	return routes
}