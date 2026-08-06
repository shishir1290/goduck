package container

import "reflect"

type Container struct {
    services map[reflect.Type]any
}

func New() *Container {
    return &Container{
        services: make(map[reflect.Type]any),
    }
}

func (c *Container) Singleton(service any) {

    t := reflect.TypeOf(service)

    c.services[t] = service
}

func (c *Container) Resolve(t reflect.Type) any {

    return c.services[t]
}