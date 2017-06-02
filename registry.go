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

type registry struct {
	lock    sync.RWMutex
	plugins map[string]Plugin
}

func (r *registry) register(p Plugin) {
	r.lock.Lock()
	defer r.lock.Unlock()

	if r.plugins == nil {
		r.plugins = make(map[string]Plugin)
	}

	r.plugins[p.Name] = p
}

func (r *registry) deregister(name string) {
	r.lock.Lock()
	defer r.lock.Unlock()

	delete(r.plugins, name)
}

func (r *registry) list() map[string]Plugin {
	return r.plugins
}
