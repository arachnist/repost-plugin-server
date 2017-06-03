package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"path"
	"plugin"

	"golang.org/x/net/trace"

	"github.com/arachnist/repost-plugin-server"
	"github.com/arachnist/repost-plugin-server/types"
)

var (
	bindAddress string
	baseConfig  string
	plugins     string
)

func main() {
	RPS := rps.New(baseConfig)
	tr := trace.New("rps.init", "")
	ctx := trace.NewContext(context.Background(), tr)

	for _, name := range RPS.Config.Lookup(ctx, nil, "plugins") {
		tr.LazyPrintf("Loading plugin %s", name)
		p, err := plugin.Open(path.Join(plugins, name+".so"))
		if err != nil {
			tr.SetError()
			tr.LazyPrintf("Error loading plugin %s: %s", name, err.Error())
			continue
		}
		listSym, err := p.Lookup("List")
		if err != nil {
			tr.SetError()
			tr.LazyPrintf("Plugin %s List function lookup failed: %s", name, err.Error())
			continue
		}
		for _, plug := range listSym.(func() []types.Plugin)() {
			RPS.WrapAPI(plug)
		}
	}

	tr.Finish()

	http.ListenAndServe(bindAddress, nil)
}

func init() {
	flag.StringVar(&bindAddress, "bind_address", ":8081", "Address to bind the web api")
	flag.StringVar(&baseConfig, "base_config", path.Join(os.Getenv("HOME"), ".repost", "config"), "Base configuration directory")
	flag.StringVar(&plugins, "plugins", path.Join(os.Getenv("HOME"), ".repost", "plugins"), "Base configuration directory")
	flag.Parse()
}
