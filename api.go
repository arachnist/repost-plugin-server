package rps

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/net/trace"
)

type request struct {
	sender    string   `json:"sender"`
	recipient string   `json:"recipient"`
	message   []string `json:"message"`
}

type response struct {
	Ok      bool     `json:"ok"`
	Err     string   `json:"err",omitempty`
	Message []string `json:"message",omitempty`
}

func WrapAPI(name string, fun func(context.Context, string, string, []string) []string) {
	http.HandleFunc(fmt.Sprintf("/api/v1/rps/%s", name), func(w http.ResponseWriter, r *http.Request) {
		var q request
		var res response

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
		res.Message = fun(ctx, q.sender, q.recipient, q.message)
	})
}
