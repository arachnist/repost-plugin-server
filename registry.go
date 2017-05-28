package rps

import (
	"sync"
)

type plugin func(string, string, []string) []string

var registry struct {
	lock    sync.RWMutex
	plugins map[string]plugin
}

func Register(name string, function plugin) {
	registry.lock.Lock()
	defer registry.lock.Unlock()

	if registry.plugins == nil {
		registry.plugins = make(map[string]plugin)
	}

	registry.plugins[name] = function
}

func Plugins() (rplugins map[string]plugin) {
	registry.lock.RLock()
	defer registry.lock.RUnlock()

	rplugins = make(map[string]plugin)
	for k, v := range registry.plugins {
		rplugins[k] = v
	}

	return rplugins
}
