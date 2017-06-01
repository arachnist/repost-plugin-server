package rps

import (
	"context"
	"sync"
)

type Plugin struct {
	Name      string
	Call      func(context.Context, map[string][]string, Request) Response
	Variables []string
}

var registry struct {
	lock    sync.RWMutex
	plugins []Plugin
}

func Register(p Plugin) {
	registry.lock.Lock()
	defer registry.lock.Unlock()

	registry.plugins = append(registry.plugins, p)
}

func Plugins() []Plugin {
	return registry.plugins
}
