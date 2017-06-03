package main

import (
	"context"
	"flag"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path"

	"golang.org/x/net/trace"

	"github.com/arachnist/repost-plugin-server"
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

	RPS.Mux.Handle("/", http.DefaultServeMux)

	for _, name := range RPS.Config.Lookup(ctx, nil, "plugins") {
		tr.LazyPrintf("Loading plugin %s", name)
		err := RPS.Load(ctx, path.Join(plugins, name+".so"))
		if err != nil {
			panic(err)
		}
	}

	tr.Finish()

	http.ListenAndServe(bindAddress, RPS.Mux)
}

func init() {
	flag.StringVar(&bindAddress, "bind_address", ":8081", "Address to bind the web api")
	flag.StringVar(&baseConfig, "base_config", path.Join(os.Getenv("HOME"), ".repost", "config"), "Base configuration directory")
	flag.StringVar(&plugins, "plugins", path.Join(os.Getenv("HOME"), ".repost", "plugins"), "Base plugin directory")
	flag.Parse()
}
