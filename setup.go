package minify

import (
	"strings"

	"github.com/mholt/caddy"
	"github.com/mholt/caddy/caddyhttp/httpserver"
)

// Setup is the init function of Caddy plugins and it configures the whole
// middleware thing.
func setup(c *caddy.Controller) error {
	cnf := httpserver.GetConfig(c.Key)
	excludes := []string{}
	basePath := "/"

	for c.Next() {
		args := c.RemainingArgs()

		switch len(args) {
		case 1:
			basePath = args[0]
			basePath = strings.TrimSuffix(basePath, "/")
			basePath += "/"
		}

		for c.NextBlock() {
			switch c.Val() {
			case "exclude":
				if !c.NextArg() {
					return c.ArgErr()
				}
				excludes = strings.Split(c.Val(), " ")
			}
		}
	}

	mid := func(next httpserver.Handler) httpserver.Handler {
		return Minify{
			Next:     next,
			Excludes: excludes,
			BasePath: basePath,
		}
	}

	cnf.AddMiddleware(mid)
	return nil
}
