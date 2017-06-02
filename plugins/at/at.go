package main

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/arachnist/repost-plugin-server"
	"github.com/arachnist/repost-plugin-server/util"
)

type user struct {
	Timestamp  float64
	Login      string
	PrettyTime string `json:"pretty_time"`
}

type checkinator struct {
	Kektops int
	Esps    int
	Unknown int
	Users   []user
}

func at(ctx context.Context, config map[string][]string, request rps.Request) (response rps.Response) {
	var values checkinator
	var recently []string
	var now []string

	response.Ok = true

	data, err := util.HTTPGet(config["URL"][0])
	if err != nil {
		response.Ok = false
		response.Err = err.Error()
		return response
	}

	err = json.Unmarshal(data, &values)
	if err != nil {
		response.Ok = false
		response.Err = err.Error()
		return response
	}

	response.Message = []string{"at:"}

	for _, u := range values.Users {
		t := time.Unix(int64(u.Timestamp), 0)
		if t.Add(time.Minute * 10).After(time.Now()) {
			now = append(now, u.Login)
		} else {
			recently = append(recently, u.Login)
		}
	}

	if len(now) > 0 {
		response.Message = append(append(response.Message, "now:"), now...)
	}
	if len(recently) > 0 {
		response.Message = append(append(response.Message, "recently:"), recently...)
	}
	if len(now) == 0 && len(recently) == 0 {
		response.Message = append(response.Message, config["empty"]...)
	}

	if values.Kektops > 0 {
		response.Message = append(response.Message, []string{"kektops:", strconv.Itoa(values.Kektops)}...)
	}

	if values.Esps > 0 {
		response.Message = append(response.Message, []string{"esps:", strconv.Itoa(values.Esps)}...)
	}

	if values.Unknown > 0 {
		response.Message = append(response.Message, []string{"kektops:", strconv.Itoa(values.Unknown)}...)
	}

	return response
}

func List() []rps.Plugin {
	return []rps.Plugin{
		{"at", at, []string{"URL", "empty"}},
	}
}
