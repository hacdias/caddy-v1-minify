package minify

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"strconv"
	"strings"

	"github.com/mholt/caddy/caddyhttp/httpserver"

	"github.com/mholt/caddy"

	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/html"
	"github.com/tdewolff/minify/js"
	"github.com/tdewolff/minify/json"
	"github.com/tdewolff/minify/svg"
	"github.com/tdewolff/minify/xml"
)

var (
	m         *minify.M
	jsonRegex = regexp.MustCompile("[/+]json$")
	xmlRegex  = regexp.MustCompile("[/+]xml$")
)

func init() {
	m = minify.New()
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("text/html", html.Minify)
	m.AddFunc("text/javascript", js.Minify)
	m.AddFunc("image/svg+xml", svg.Minify)
	m.AddFuncRegexp(jsonRegex, json.Minify)
	m.AddFuncRegexp(xmlRegex, xml.Minify)

	caddy.RegisterPlugin("minify", caddy.Plugin{
		ServerType: "http",
		Action:     setup,
	})
}

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
		return CaddyMinify{
			Next:     next,
			Excludes: excludes,
			BasePath: basePath,
		}
	}

	cnf.AddMiddleware(mid)
	return nil
}

// CaddyMinify is [finish this]
type CaddyMinify struct {
	Next     httpserver.Handler
	Excludes []string
	BasePath string
}

// ServeHTTP is the main function of the whole plugin that routes every single
// request to its function.
func (h CaddyMinify) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	if httpserver.Path(r.URL.Path).Matches(h.BasePath) {
		if isExcluded(strings.Replace(r.URL.Path, h.BasePath, "/", 1), h.Excludes) {
			return h.Next.ServeHTTP(w, r)
		}

		rec := httptest.NewRecorder()
		code, err := h.Next.ServeHTTP(rec, r)
		data := rec.Body.Bytes()

		if val, ok := rec.Header()["Content-Type"]; ok {
			r := regexp.MustCompile(`(\w+\/[\w-]+)`)
			matches := r.FindStringSubmatch(val[0])

			if len(matches) != 0 && canBeMinified(matches[0]) {
				data, err = m.Bytes(matches[0], data)
				if err != nil {
					return 500, err
				}
			}
		}

		// copy the original headers
		for k, v := range rec.Header() {
			if k == "Content-Length" {
				w.Header().Set("Content-Length", strconv.Itoa(len(data)))
				continue
			}

			w.Header()[k] = v
		}

		w.Write(data)
		return code, err
	}

	return h.Next.ServeHTTP(w, r)
}

func isExcluded(path string, excludes []string) bool {
	for _, el := range excludes {
		if httpserver.Path(path).Matches(el) {
			return true
		}
	}

	return false
}

func canBeMinified(mediatype string) bool {
	switch mediatype {
	case "text/css", "text/html", "text/javascript", "image/svg+xml":
		return true
	}

	if jsonRegex.FindString(mediatype) != "" {
		return true
	}

	if xmlRegex.FindString(mediatype) != "" {
		return true
	}

	return false
}
