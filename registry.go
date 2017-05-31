package rps

import (
	"context"
	"sync"
)

type Plugin func(context.Context, Request) Response

var registry struct {
	lock    sync.RWMutex
	plugins map[string]Plugin
}

func Register(name string, function Plugin) {
	registry.lock.Lock()
	defer registry.lock.Unlock()

	if registry.plugins == nil {
		registry.plugins = make(map[string]Plugin)
	}

	registry.plugins[name] = function
}

func Plugins() (rplugins map[string]Plugin) {
	registry.lock.RLock()
	defer registry.lock.RUnlock()

	rplugins = make(map[string]Plugin)
	for k, v := range registry.plugins {
		rplugins[k] = v
	}

	return rplugins
}
