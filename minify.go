package minify

import (
	"net/http"
	"net/http/httptest"
	"regexp"
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
	if strings.HasSuffix(r.URL.Path, ".css") {

		rec := httptest.NewRecorder()

		code, err := h.Next.ServeHTTP(rec, r)

		m := minify.New()
		m.AddFunc("text/css", css.Minify)
		m.AddFunc("text/html", html.Minify)
		m.AddFunc("text/javascript", js.Minify)
		m.AddFunc("image/svg+xml", svg.Minify)
		m.AddFuncRegexp(regexp.MustCompile("[/+]json$"), json.Minify)
		m.AddFuncRegexp(regexp.MustCompile("[/+]xml$"), xml.Minify)

		b, err := m.Bytes("text/css", rec.Body.Bytes())
		if err != nil {
			panic(err)
		}

		w.Write(b)

		//	w.Write(rec.Body.Bytes())

		return code, err
	}

	return h.Next.ServeHTTP(w, r)
}
