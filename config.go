package rps

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"sync"
	"time"

	"golang.org/x/net/trace"
)

type cacheEntry struct {
	modTime  time.Time
	contents map[string][]string
}

type config struct {
	basedir string
	cache   map[string]cacheEntry
	lock    sync.Mutex
}

func (c *config) fileList(env map[string]string) (r []string) {
	if env["network"] != "" {
		if env["recipient"] != "" {
			if env["sender"] != "" {
				r = append(r, []string{path.Join(c.basedir, env["network"], env["recipient"], env["sender"], env["plugin"]+".json"),
					path.Join(c.basedir, env["network"], env["recipient"], env["sender"]+".json"),
					path.Join(c.basedir, env["network"], env["sender"]+".json")}...)
			}

			r = append(r, []string{path.Join(c.basedir, env["network"], env["recipient"], env["plugin"]+".json"),
				path.Join(c.basedir, env["network"], env["recipient"]+".json")}...)
		}

		r = append(r, []string{path.Join(c.basedir, env["network"], env["plugin"]+".json"),
			path.Join(c.basedir, env["network"]+".json")}...)
	}

	return append(r, []string{path.Join(c.basedir, env["plugin"]+".json"),
		path.Join(c.basedir, "common.json")}...)
}

func (c *config) cacheUpdate(ctx context.Context, file string) error {
	tr, _ := trace.FromContext(ctx)
	var f map[string][]string

	i, err := os.Stat(file)
	_, ok := c.cache[file]
	if os.IsNotExist(err) && ok {
		tr.LazyPrintf("Purging cache: file %s was removed", file)
		delete(c.cache, file)
	}

	if err != nil {
		tr.LazyPrintf("Error occured: %s", err.Error())
		return err
	}

	if c.cache[file].modTime.Before(i.ModTime()) || !ok {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			tr.LazyPrintf("Purging cache: file %s is unreadable", file)
			delete(c.cache, file)
			return err
		}

		err = json.Unmarshal(data, &f)
		if err != nil {
			tr.LazyPrintf("Purging cache: file %s is unparsable", file)
			delete(c.cache, file)
			return err
		}

		for key, _ := range f {
			sort.Strings(f[key])
		}

		tr.LazyPrintf("Updating cache: file %s", file)
		c.cache[file] = cacheEntry{
			modTime:  i.ModTime(),
			contents: f,
		}
	}

	return nil
}

func (c *config) Lookup(ctx context.Context, env map[string]string, key string) []string {
	tr, _ := trace.FromContext(ctx)
	c.lock.Lock()
	defer c.lock.Unlock()

	for _, file := range c.fileList(env) {
		if c.cacheUpdate(ctx, file) != nil {
			continue
		}

		if value, exists := c.cache[file].contents[key]; exists {
			tr.LazyPrintf("Lookup ok for key %s with env %+v value %+v", key, env, value)
			return value
		}
	}

	tr.LazyPrintf("Lookup failed for key %s with env %+v", key, env)
	return nil
}
