package main

import (
	"fmt"

	"github.com/arachnist/repost-plugin-server"
	_ "github.com/arachnist/repost-plugin-server/plugins"
)

func main() {
	for k, v := range rps.Plugins() {
		fmt.Println(k)
		v("", "", []string{""})
	}
}
