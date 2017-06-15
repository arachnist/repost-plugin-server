package rps

import (
	"context"

	"github.com/arachnist/repost-plugin-server/types"
)

func (s *server) list(context.Context, map[string][]string, types.Request) (r types.Response) {
	s.lock.Lock()
	defer s.lock.Unlock()
	var msg map[string]string
	msg = make(map[string]string)
	for _, p := range s.plugins {
		msg[p.Name] = p.Trigger
	}
	r.Message = msg
	r.Ok = true
	return
}

func (s *server) load(ctx context.Context, env map[string][]string, q types.Request) (r types.Response) {
	r.Ok = true
	for _, name := range q.Message {
		if err := s.Load(ctx, name); err != nil {
			r.Ok = false
			r.Err = "Plugin load failed"
			r.Message = append(r.Message.([]string), err.Error())
		}
	}
	return
}

func (s *server) unload(ctx context.Context, env map[string][]string, q types.Request) (r types.Response) {
	r.Ok = true
	for _, name := range q.Message {
		s.Deregister(name)
	}
	return
}
