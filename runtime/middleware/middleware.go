package middleware

import (
	"github.com/shishir1290/goduck/runtime/context"
	"github.com/shishir1290/goduck/runtime/handler"
)

type Middleware func(handler.HandlerFunc) handler.HandlerFunc

func Identity(next handler.HandlerFunc) handler.HandlerFunc {
	return next
}

func Apply(
	final handler.HandlerFunc,
	middlewares []Middleware,
) handler.HandlerFunc {

	h := final

	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}

	return h
}

// Keep the import used until Context gains more helpers.
var _ = context.Context{}