package router

import (
	"github.com/shishir1290/goduck/runtime/handler"
	"github.com/shishir1290/goduck/runtime/middleware"
)

type Route struct {
	Method string

	Path string

	Handler handler.HandlerFunc

	Middlewares []middleware.Middleware
}