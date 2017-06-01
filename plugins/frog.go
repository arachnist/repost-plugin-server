package plugins

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/arachnist/repost-plugin-server"
	"github.com/arachnist/repost-plugin-server/util"
)

type tip struct {
	Number int    `json:"number"`
	Tip    string `json:"tip"`
}

type tips struct {
	Tips []tip `json:"tips"`
	lock sync.RWMutex
}

func (tips *tips) fetchTips(url string) error {
	tips.lock.Lock()
	defer tips.lock.Unlock()

	data, err := util.HTTPGet(url)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, tips)
	if err != nil {
		return err
	}

	return nil
}

func (tips *tips) popTip(url string) (string, error) {
	if len(tips.Tips) == 0 {
		if err := tips.fetchTips(url); err != nil {
			return "", err
		}
	}

	tips.lock.Lock()
	defer tips.lock.Unlock()

	rmsg := tips.Tips[len(tips.Tips)-1].Tip
	tips.Tips = tips.Tips[:len(tips.Tips)-1]

	return rmsg, nil
}

var t tips

func frog(ctx context.Context, config map[string][]string, request rps.Request) (response rps.Response) {
	tip, err := t.popTip(config["URL"][0])
	if err != nil {
		response.Ok = false
		response.Err = err.Error()
	} else {
		response.Ok = true
		response.Message = []string{tip}
	}
	return
}

func init() {
	rps.Register(rps.Plugin{"frog", frog, []string{"URL"}})
}
