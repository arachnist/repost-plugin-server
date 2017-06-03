package rps

import (
	"sync"

	"github.com/arachnist/repost-plugin-server/types"
)

type server struct {
	Config  *Config
	plugins map[string]types.Plugin
	Mux     *ServeMux
	lock    sync.RWMutex
}

func New(basedir string) *server {
	return &server{
		Config:  NewConfig(basedir),
		plugins: make(map[string]types.Plugin),
		Mux:     NewServeMux(),
	}
}

func (r *server) register(p types.Plugin) {
	r.lock.Lock()
	defer r.lock.Unlock()

	if r.plugins == nil {
		r.plugins = make(map[string]types.Plugin)
	}

	r.plugins[p.Name] = p
}

func (r *server) deregister(name string) {
	r.lock.Lock()
	defer r.lock.Unlock()

	delete(r.plugins, name)
	r.Mux.Deregister("/api/v1/rps/" + name)
}

func (r *server) list() map[string]types.Plugin {
	return r.plugins
}
