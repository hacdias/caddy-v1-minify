package minify

import (
	"errors"
	"regexp"
	"strconv"

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

	minifiers := map[string]minify.Minifier{
		"css":  &css.Minifier{},
		"html": &html.Minifier{},
		"svg":  &svg.Minifier{},
		"json": &json.Minifier{},
		"xml":  &xml.Minifier{},
		"js":   &js.Minifier{},
	}

	for c.Next() {
		paths = append(paths, c.RemainingArgs()...)

		matcher, err := httpserver.SetupIfMatcher(c)
		if err != nil {
			return err
		}

		rules = append(rules, matcher)

		for c.NextBlock() {
			if httpserver.IfMatcherKeyword(c) {
				continue
			}

			switch opt := c.Val(); opt {
			case "disable":
				for c.NextArg() {
					delete(minifiers, c.Val())
				}
			case "css", "svg":
				if !c.NextArg() {
					return c.ArgErr()
				}

				if _, ok := minifiers[opt]; !ok {
					return errors.New("You have disabled " + opt + " minifier")
				}

				option := c.Val()
				if option != "decimals" {
					return c.ArgErr()
				}

				if !c.NextArg() {
					return c.ArgErr()
				}

				decimals, err := strconv.Atoi(c.Val())
				if err != nil {
					return err
				}

				if opt == "css" {
					minifiers["css"].(*css.Minifier).Decimals = decimals
					continue
				}

				minifiers["svg"].(*svg.Minifier).Decimals = decimals
			case "xml":
				if !c.NextArg() {
					return c.ArgErr()
				}

				if _, ok := minifiers["xml"]; !ok {
					return errors.New("You have disabled xml minifier")
				}

				option := c.Val()
				if option != "keep_whitespace" {
					return c.ArgErr()
				}

				val := true

				if !c.NextArg() {
					minifiers["xml"].(*xml.Minifier).KeepWhitespace = val
					continue
				}

				val, err = strconv.ParseBool(c.Val())
				if err != nil {
					return err
				}

				minifiers["xml"].(*xml.Minifier).KeepWhitespace = val
			case "html":
				if !c.NextArg() {
					return c.ArgErr()
				}

				if _, ok := minifiers["html"]; !ok {
					return errors.New("You have disabled html minifier")
				}

				option := c.Val()
				val := true

				if c.NextArg() {
					val, err = strconv.ParseBool(c.Val())
					if err != nil {
						return err
					}
				}

				switch option {
				case "keep_default_attr_vals":
					minifiers["html"].(*html.Minifier).KeepDefaultAttrVals = val
				case "keep_document_tags":
					minifiers["html"].(*html.Minifier).KeepDocumentTags = val
				case "keep_end_tags":
					minifiers["html"].(*html.Minifier).KeepEndTags = val
				case "keep_whitespace":
					minifiers["html"].(*html.Minifier).KeepWhitespace = val
				default:
					return errors.New("Unknown option " + option)
				}
			}
		}
	}

	for name, fn := range minifiers {
		switch name {
		case "css":
			minifier.Add("text/css", fn)
		case "html":
			minifier.Add("text/html", fn)
		case "svg":
			minifier.Add("image/svg+xml", fn)
		case "json":
			minifier.AddRegexp(regexp.MustCompile("[/+]json$"), fn)
		case "xml":
			minifier.AddRegexp(regexp.MustCompile("[/+]xml$"), fn)
		case "js":
			minifier.AddRegexp(regexp.MustCompile("[/-]javascript$"), fn)
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
