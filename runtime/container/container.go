package container

import "reflect"

type Container struct {
	providers map[reflect.Type]Provider
	instances map[reflect.Type]any
}

func New() *Container {

	return &Container{
		providers: make(map[reflect.Type]Provider),
		instances: make(map[reflect.Type]any),
	}
}