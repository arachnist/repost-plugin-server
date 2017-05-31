package rps

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"time"

	"golang.org/x/net/trace"
)

type cacheEntry struct {
	modTime  time.Time
	contents map[string][]string
}

type Config struct {
	basedir string
	cache   map[string]cacheEntry
}

func (c *Config) fileList(env map[string]string) (r []string) {
	if env["network"] != "" {
		if env["recipient"] != "" {
			if env["sender"] != "" {
				r = append(r, path.Join(c.basedir, env["network"], env["recipient"], env["sender"]+".json"))
				r = append(r, path.Join(c.basedir, env["network"], env["sender"]+".json"))
			}
			r = append(r, path.Join(c.basedir, env["network"], env["recipient"]+".json"))
		}
		r = append(r, path.Join(c.basedir, env["network"]+".json"))
	}

	return append(r, path.Join(c.basedir, "common.json"))
}

func (c *Config) cacheUpdate(ctx context.Context, file string) error {
	tr, _ := trace.FromContext(ctx)
	var f map[string][]string

	i, err := os.Stat(file)
	_, ok := c.cache[file]
	if os.IsNotExist(err) && ok {
		tr.LazyPrintf("Purging cache: file %s was removed", file)
		delete(c.cache, file)
	}

	if err != nil {
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

		tr.LazyPrintf("Updating cache: file %s", file)
		c.cache[file] = cacheEntry{
			modTime:  i.ModTime(),
			contents: f,
		}
	}

	return nil
}

func (c *Config) Lookup(ctx context.Context, env map[string]string, key string) []string {
	tr, _ := trace.FromContext(ctx)

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
