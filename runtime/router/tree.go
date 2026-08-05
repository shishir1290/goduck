package router

import "strings"

type Tree struct {
	root *Node
}

func NewTree() *Tree {
	return &Tree{
		root: NewNode("/"),
	}
}

func (t *Tree) Insert(route *Route) {

	current := t.root

	path := strings.Trim(route.Path, "/")

	if path == "" {

		current.Routes[route.Method] = route

		return
	}

	parts := strings.Split(path, "/")

	for _, part := range parts {

		// parameter

		if strings.HasPrefix(part, ":") {

			if current.Parameter == nil {
				current.Parameter = NewNode(part)
			}

			current = current.Parameter

			continue
		}

		child, ok := current.Children[part]

		if !ok {

			child = NewNode(part)

			current.Children[part] = child
		}

		current = child
	}

	current.Routes[route.Method] = route
}