package main

import (
	"context"
	"flag"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/net/trace"

	"github.com/arachnist/repost-plugin-server"
)

var (
	bindAddress string
	baseConfig  string
	plugins     string
	apikey      string
)

func main() {
	RPS := rps.New(baseConfig, plugins, apikey)
	tr := trace.New("rps.init", "")
	ctx := trace.NewContext(context.Background(), tr)

	RPS.Mux.Handle("/", http.DefaultServeMux)
	RPS.Mux.Handle("/metrics", promhttp.Handler())

	for _, name := range RPS.Config.Lookup(ctx, nil, "plugins") {
		tr.LazyPrintf("Loading plugin %s", name)
		err := RPS.Load(ctx, name)
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
	flag.StringVar(&apikey, "api_key", "f33dc0d3", "API key")
	flag.Parse()
}
