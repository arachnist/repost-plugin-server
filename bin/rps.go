package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/golang/glog"
	"github.com/gorilla/handlers"

	"github.com/arachnist/repost-plugin-server"
	_ "github.com/arachnist/repost-plugin-server/plugins"
)

var (
	bindAddress string
)

func main() {
	glog.Info("Starting repost plugin server...")
	for name, plugin := range rps.Plugins() {
		glog.Info("Registering plugin", name)
		rps.WrapAPI(name, plugin)
	}

	glog.Error(http.ListenAndServe(bindAddress, handlers.CombinedLoggingHandler(os.Stdout, http.DefaultServeMux)))
}

func init() {
	flag.StringVar(&bindAddress, "bind_address", ":8081", "Address to bind the web api")
	flag.Parse()
}
