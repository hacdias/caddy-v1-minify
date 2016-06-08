package minify

import (
	"net/http"

	"github.com/mholt/caddy"
	"github.com/mholt/caddy/caddyhttp/httpserver"
)

func init() {
	caddy.RegisterPlugin("minify", caddy.Plugin{
		ServerType: "http",
		Action:     setup,
	})
}

// Setup is the init function of Caddy plugins and it configures the whole
// middleware thing.
func setup(c *caddy.Controller) error {
	cnf := httpserver.GetConfig(c.Key)

	mid := func(next httpserver.Handler) httpserver.Handler {
		return CaddyMinify{Next: next}
	}

	cnf.AddMiddleware(mid)
	return nil
}

type CaddyMinify struct {
	Next httpserver.Handler
}

// ServeHTTP is the main function of the whole plugin that routes every single
// request to its function.
func (h CaddyMinify) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {

	return h.Next.ServeHTTP(w, r)
}
