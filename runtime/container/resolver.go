package container

import (
	"reflect"
	"strings"
)

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
	if provider == nil {
		return nil
	}

	instance := provider(c)

	c.instances[t] = instance

	return instance
}

func (c *Container) Has(target any) bool {
	t := reflect.TypeOf(target)
	_, ok := c.providers[t]
	return ok
}

func (c *Container) HasTypeWithName(name string) bool {
	for t := range c.providers {
		if strings.Contains(strings.ToLower(t.String()), strings.ToLower(name)) {
			return true
		}
	}
	return false
}