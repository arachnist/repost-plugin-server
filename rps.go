package rps

import (
	"context"
	"errors"
	"path"
	"path/filepath"
	"plugin"
	"strings"
	"sync"

	"golang.org/x/net/trace"

	"github.com/arachnist/repost-plugin-server/types"
)

type server struct {
	Config    *Config
	plugins   map[string]types.Plugin
	Mux       *ServeMux
	lock      sync.RWMutex
	apikey    string
	pluginDir string
}

func New(basedir, plugins, apikey string) (s *server) {
	s = &server{
		Config:    NewConfig(basedir),
		plugins:   make(map[string]types.Plugin),
		Mux:       NewServeMux(),
		apikey:    apikey,
		pluginDir: plugins,
	}
	for _, plug := range []types.Plugin{
		{"mgmt/list", s.list, []string{}, "list"},
		{"mgmt/load", s.load, []string{}, "load"},
		{"mgmt/unload", s.unload, []string{}, "unload"},
	} {
		s.Register(plug)
	}
	return s
}

func (r *server) Load(ctx context.Context, pl string) error {
	tr, _ := trace.FromContext(ctx)
	pl = filepath.Clean(path.Join(r.pluginDir, pl+".so"))
	if !strings.HasPrefix(pl, r.pluginDir) {
		tr.SetError()
		tr.LazyPrintf("invalid plugin name")
		return errors.New("Invalid plugin name")
	}
	p, err := plugin.Open(pl)
	if err != nil {
		tr.SetError()
		tr.LazyPrintf("Error loading plugin %s: %s", pl, err.Error())
		return err
	}
	listSym, err := p.Lookup("List")
	if err != nil {
		tr.SetError()
		tr.LazyPrintf("Plugin %s List function lookup failed: %s", pl, err.Error())
		return err
	}
	for _, plug := range listSym.(func() []types.Plugin)() {
		r.Register(plug)
	}
	return nil
}

func (r *server) Register(p types.Plugin) {
	r.lock.Lock()
	defer r.lock.Unlock()

	if r.plugins == nil {
		r.plugins = make(map[string]types.Plugin)
	}

	r.wrapAPI(p)
	r.plugins[p.Name] = p
}

func (r *server) Deregister(name string) {
	r.lock.Lock()
	defer r.lock.Unlock()

	delete(r.plugins, name)
	r.Mux.Deregister("/api/v1/rps/" + name)
}
