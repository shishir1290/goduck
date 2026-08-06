package container

import "reflect"

func (c *Container) Register(
	target any,
	provider Provider,
) {

	t := reflect.TypeOf(target)

	c.providers[t] = provider
}

func (c *Container) Resolve(
	target any,
) any {

	t := reflect.TypeOf(target)

	// singleton

	if instance, ok := c.instances[t]; ok {
		return instance
	}

	provider := c.providers[t]

	instance := provider(c)

	c.instances[t] = instance

	return instance
}