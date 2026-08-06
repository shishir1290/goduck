package engine

import (
	"github.com/shishir1290/goduck/runtime/config"
	"github.com/shishir1290/goduck/runtime/container"
	"github.com/shishir1290/goduck/runtime/middleware"
	"github.com/shishir1290/goduck/runtime/router"
	"github.com/shishir1290/goduck/runtime/static"
)

type Engine struct {

    Router *router.Router

    Middlewares []middleware.Middleware

    Container *container.Container

	Config *config.Config

	Statics []*static.Static
}

func New() *Engine {

    return &Engine{
        Router: router.New(),
        Container: container.New(),
		Config: config.Default(),
    }
}