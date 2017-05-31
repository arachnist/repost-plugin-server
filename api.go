package rps

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/net/trace"
)

type Request struct {
	network   string   `json:"network"`
	sender    string   `json:"sender"`
	recipient string   `json:"recipient"`
	message   []string `json:"message"`
}

type Response struct {
	Ok      bool     `json:"ok"`
	Err     string   `json:"err",omitempty`
	Message []string `json:"message",omitempty`
}

func WrapAPI(name string, fun Plugin) {
	http.HandleFunc(fmt.Sprintf("/api/v1/rps/%s", name), func(w http.ResponseWriter, r *http.Request) {
		var q Request
		var res Response

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

		err := json.NewDecoder(r.Body).Decode(&q)
		if err != nil {
			tr.LazyPrintf("Error decoding request body: %+v", err)
			tr.SetError()
			w.WriteHeader(503)
			res.Ok = false
			res.Err = err.Error()
			return
		}

		res.Ok = true
		res = fun(ctx, q)
		if res.Ok != true {
			tr.LazyPrintf("Plugin %s error: %s", name, res.Err)
			tr.SetError()
		}
	})
}
