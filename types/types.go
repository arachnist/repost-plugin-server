package types

import (
	"context"
)

type Plugin struct {
	Name      string
	Call      func(context.Context, map[string][]string, Request) Response
	Variables []string
	Trigger   string
}

type Request struct {
	Network   string   `json:"network"`
	Sender    string   `json:"sender"`
	Recipient string   `json:"recipient"`
	Message   []string `json:"message"`
}

type Response struct {
	Ok      bool        `json:"ok"`
	Err     string      `json:"err",omitempty`
	Message interface{} `json:"message",omitempty`
}
