package minify

import (
	"strings"

	"github.com/mholt/caddy"
	"github.com/mholt/caddy/caddyhttp/httpserver"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/html"
	"github.com/tdewolff/minify/js"
	"github.com/tdewolff/minify/json"
	"github.com/tdewolff/minify/svg"
	"github.com/tdewolff/minify/xml"
)

// init initializes the minifier variable and configures the plugin on
// Caddy webserver.
func init() {
	minifier = minify.New()
	minifier.AddFunc("css", css.Minify)
	minifier.AddFunc("html", html.Minify)
	minifier.AddFunc("javascript", js.Minify)
	minifier.AddFunc("svg", svg.Minify)
	minifier.AddFunc("json", json.Minify)
	minifier.AddFunc("xml", xml.Minify)

	caddy.RegisterPlugin("minify", caddy.Plugin{
		ServerType: "http",
		Action:     setup,
	})
}

// setup configures the middlware.
func setup(c *caddy.Controller) error {
	cnf := httpserver.GetConfig(c)
	excludes, includes, err := parse(c)

	if err != nil {
		return err
	}

	mid := func(next httpserver.Handler) httpserver.Handler {
		return Minify{
			Next:     next,
			Excludes: excludes,
			Includes: includes,
		}
	}

	cnf.AddMiddleware(mid)
	return nil
}

// parse parses the configuration of the plugin using caddy.Controller.
func parse(c *caddy.Controller) ([]string, []string, error) {
	excludes := []string{}
	includes := []string{"/"}

	for c.Next() {
		for c.NextBlock() {
			switch c.Val() {
			case "exclude":
				if !c.NextArg() {
					return []string{}, []string{}, c.ArgErr()
				}
				excludes = strings.Split(c.Val(), " ")
			case "only":
				if !c.NextArg() {
					return []string{}, []string{}, c.ArgErr()
				}
				includes = strings.Split(c.Val(), " ")
			}
		}
	}

	return excludes, includes, nil
}
