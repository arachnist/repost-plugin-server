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

type server struct {
	Config  *Config
	plugins map[string]Plugin
	mux     *ServeMux
	lock    sync.RWMutex
}

func New(basedir string) *server {
	// rps.cache = make(map[string]cacheEntry)
	return &server{
		Config:  NewConfig(basedir),
		plugins: make(map[string]Plugin),
		mux:     NewServeMux(),
	}
}

func (r *server) register(p Plugin) {
	r.lock.Lock()
	defer r.lock.Unlock()

	if r.plugins == nil {
		r.plugins = make(map[string]Plugin)
	}

	r.plugins[p.Name] = p
}

func (r *server) deregister(name string) {
	r.lock.Lock()
	defer r.lock.Unlock()

	delete(r.plugins, name)
	r.mux.Deregister("/api/v1/rps/" + name)
}

func (r *server) list() map[string]Plugin {
	return r.plugins
}
