package minify

import (
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

// Rules is a type that stores the rules to monify
type Rules struct {
	Excludes, Includes []string
	Matches            []httpserver.RequestMatcher
}

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
	rules, err := parse(c)

	if err != nil {
		return err
	}

	mid := func(next httpserver.Handler) httpserver.Handler {
		return Minify{
			Next:  next,
			Rules: rules,
		}
	}

	cnf.AddMiddleware(mid)
	return nil
}

// parse parses the configuration of the plugin using caddy.Controller.
func parse(c *caddy.Controller) (Rules, error) {
	rules := Rules{}

	for c.Next() {
		matcher, err := httpserver.SetupIfMatcher(c)
		if err != nil {
			return rules, err
		}

		rules.Matches = append(rules.Matches, matcher)

		for c.NextBlock() {
			switch c.Val() {
			case "exclude":
				if !c.NextArg() {
					return rules, c.ArgErr()
				}
				rules.Excludes = append(rules.Excludes, c.Val())
			case "include":
				if !c.NextArg() {
					return rules, c.ArgErr()
				}
				rules.Includes = append(rules.Includes, c.Val())
			}
		}
	}

	return rules, nil
}
