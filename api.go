package rps

import (
	"context"
	"encoding/json"
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"golang.org/x/net/trace"

	"github.com/arachnist/repost-plugin-server/types"
)

func (srv *server) wrapAPI(p types.Plugin) {
	srv.Mux.HandleFunc("/api/v1/rps/"+p.Name, func(w http.ResponseWriter, r *http.Request) {
		var q types.Request
		var res types.Response
		var config = make(map[string][]string)
		var environment map[string]string
		var whitelisted bool

		tr := trace.New("rps.api", r.URL.String())
		tr.LazyPrintf("API request from %s", r.RemoteAddr)
		ctx := trace.NewContext(context.Background(), tr)

		defer func() {
			err := json.NewEncoder(w).Encode(res)
			if err != nil {
				tr.LazyPrintf("Error encoding response: %+v", err)
				tr.SetError()
			}
			tr.Finish()
		}()

		val, ok := r.Header["X-Api-Key"]
		if !ok || (ok && val[0] != srv.apikey) {
			tr.LazyPrintf("Invalid API key")
			tr.SetError()
			w.WriteHeader(403)
			res.Ok = false
			res.Err = "Invalid API key"
			return
		}

		err := json.NewDecoder(r.Body).Decode(&q)
		if err != nil {
			tr.LazyPrintf("Error decoding request body: %+v", err)
			tr.SetError()
			w.WriteHeader(503)
			res.Ok = false
			res.Err = err.Error()
			return
		}

		environment = map[string]string{
			"network":   q.Network,
			"sender":    q.Sender,
			"recipient": q.Recipient,
			"plugin":    p.Name,
		}

		// try to avoid directory traversal in configuration lookup
		for key, value := range environment {
			if !strings.HasPrefix(filepath.Clean(path.Join(srv.Config.basedir, value)), srv.Config.basedir) {
				res.Ok = false
				res.Err = "invalid env value"
				tr.LazyPrintf("Key %s invalid value %s", key, value)
				tr.SetError()
				return
			}
		}

		whitelist := srv.Config.Lookup(ctx, environment, "whitelist")
		if whitelist != nil {
			whitelisted = false
			for _, plugin := range whitelist {
				if plugin == p.Name {
					whitelisted = true
				}
			}

			if whitelisted == false {
				res.Ok = false
				res.Err = "plugin not whitelisted"
				tr.LazyPrintf("Plugin %s not whitelisted in env %+v", p.Name, environment)
				tr.SetError()
				return
			}
		}

		blacklist := srv.Config.Lookup(ctx, environment, "blacklist")
		if blacklist != nil {
			for _, plugin := range blacklist {
				if plugin == p.Name {
					res.Ok = false
					res.Err = "plugin blacklisted"
					tr.LazyPrintf("Plugin %s blacklisted in env %+v", p.Name, environment)
					tr.SetError()
					return
				}
			}
		}

		// populate plugin config
		for _, vName := range p.Variables {
			if tVal := srv.Config.Lookup(ctx, environment, vName); tVal == nil {
				res.Ok = false
				res.Err = "plugin configuration error"
				tr.LazyPrintf("Plugin %s not configured. Missing configuration key: %s", p.Name, vName)
				tr.SetError()
				return
			} else {
				config[vName] = tVal
			}
		}

		res.Ok = true
		res = p.Call(ctx, config, q)
		if res.Ok != true {
			tr.LazyPrintf("Plugin %s error: %s", p.Name, res.Err)
			tr.SetError()
		}
	})
}
