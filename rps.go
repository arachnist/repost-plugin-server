package rps

type server struct {
	*config
}

func New(basedir string) *server {
	// rps.cache = make(map[string]cacheEntry)
	return &server{
		&config{
			basedir: basedir,
			cache:   make(map[string]cacheEntry),
		},
	}
}
