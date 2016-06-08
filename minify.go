package minify

import (
	"net/http"
	"net/http/httptest"
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

var m *minify.M

func init() {
	m = minify.New()
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("text/html", html.Minify)
	m.AddFunc("text/javascript", js.Minify)
	m.AddFunc("image/svg+xml", svg.Minify)
	m.AddFuncRegexp(regexp.MustCompile("[/+]json$"), json.Minify)
	m.AddFuncRegexp(regexp.MustCompile("[/+]xml$"), xml.Minify)

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

// CaddyMinify is [finish this]
type CaddyMinify struct {
	Next httpserver.Handler
}

// ServeHTTP is the main function of the whole plugin that routes every single
// request to its function.
func (h CaddyMinify) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	rec := httptest.NewRecorder()
	code, err := h.Next.ServeHTTP(rec, r)
	data := rec.Body.Bytes()

	// copy the original headers
	for k, v := range rec.Header() {
		w.Header()[k] = v
	}

	if val, ok := rec.Header()["Content-Type"]; ok {
		r := regexp.MustCompile(`(\w+\/[\w-]+)`)
		matches := r.FindStringSubmatch(val[0])

		if len(matches) != 0 && canBeMinified(matches[0]) {
			data, err = m.Bytes(matches[0], rec.Body.Bytes())
			if err != nil {
				return 500, err
			}
		}
	}

	w.WriteHeader(code)
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.Write(data)
	return code, err
}

func canBeMinified(mediatype string) bool {
	switch mediatype {
	case "text/css", "text/html", "text/javascript", "image/svg+xml":
		return true
	}

	// add regext for "[/+]json$" and "[/+]xml$"

	return false
}
