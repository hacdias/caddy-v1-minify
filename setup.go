package minify

import (
	"regexp"

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
	minifier.AddFunc("text/css", css.Minify)
	minifier.AddFunc("text/html", html.Minify)
	minifier.AddFunc("image/svg+xml", svg.Minify)
	minifier.AddFuncRegexp(regexp.MustCompile("[/+]json$"), json.Minify)
	minifier.AddFuncRegexp(regexp.MustCompile("[/+]xml$"), xml.Minify)
	minifier.AddFuncRegexp(regexp.MustCompile("[/-]javascript$"), js.Minify)

	caddy.RegisterPlugin("minify", caddy.Plugin{
		ServerType: "http",
		Action:     setup,
	})
}

// setup configures the middlware.
func setup(c *caddy.Controller) error {
	cnf := httpserver.GetConfig(c)

	rules := []httpserver.RequestMatcher{}
	paths := []string{}

	for c.Next() {
		paths = append(paths, c.RemainingArgs()...)

		matcher, err := httpserver.SetupIfMatcher(c)
		if err != nil {
			return err
		}

		rules = append(rules, matcher)

		// TODO: Remove this to break the plugin!
		for c.NextBlock() {
		}
	}

	mid := func(next httpserver.Handler) httpserver.Handler {
		return Minify{
			Next:  next,
			Rules: rules,
			Paths: paths,
		}
	}

	cnf.AddMiddleware(mid)
	return nil
}
