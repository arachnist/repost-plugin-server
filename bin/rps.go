package main

import (
	"flag"
	"net/http"
	"os"
	"path"

	"github.com/arachnist/repost-plugin-server"
	_ "github.com/arachnist/repost-plugin-server/plugins"
)

var (
	bindAddress string
	baseConfig  string
)

func main() {
	RPS := rps.New(baseConfig)
	for _, plugin := range rps.Plugins() {
		RPS.WrapAPI(plugin)
	}

	http.ListenAndServe(bindAddress, nil)
}

func init() {
	flag.StringVar(&bindAddress, "bind_address", ":8081", "Address to bind the web api")
	flag.StringVar(&baseConfig, "base_config", path.Join(os.Getenv("HOME"), ".repost"), "Base configuration directory")
	flag.Parse()
}
