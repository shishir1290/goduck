package router

import "strings"

func (t *Tree) Find(
	method string,
	path string,
) (*Route, map[string]string) {

	current := t.root

	params := make(map[string]string)

	path = strings.Trim(path, "/")

	if path == "" {

		route := current.Routes[method]

		return route, params
	}

	parts := strings.Split(path, "/")

	for _, part := range parts {

		// Exact match

		if child, ok := current.Children[part]; ok {

			current = child

			continue
		}

		// Parameter

		if current.Parameter != nil {

			name := strings.TrimPrefix(
				current.Parameter.Segment,
				":",
			)

			params[name] = part

			current = current.Parameter

			continue
		}

		return nil, nil
	}

	route := current.Routes[method]

	if route == nil {
		return nil, nil
	}

	return route, params
}