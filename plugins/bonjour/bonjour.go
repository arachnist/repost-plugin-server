package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/arachnist/repost-plugin-server/types"
	"github.com/arachnist/repost-plugin-server/util"
)

func bonjour(ctx context.Context, config map[string][]string, request types.Request) (response types.Response) {
	t, _ := time.Parse("2006-01-02", config["startDate"][0])
	max := int(time.Now().Sub(t).Hours())/24 + 1

	img, err := util.HTTPGetXpath(config["URL"][0]+fmt.Sprintf("%d", rand.Intn(max)+1), config["xpath"][0])
	if err != nil {
		response.Ok = false
		response.Err = err.Error()
	} else {
		response.Ok = true
		response.Message = []string{img}
	}

	return
}

func List() []types.Plugin {
	return []types.Plugin{
		{"bonjour", bonjour, []string{"URL", "empty"}, "bonjour"},
	}
}
