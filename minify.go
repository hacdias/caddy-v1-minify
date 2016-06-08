package minify

import (
	"bytes"
	"mime"
	"net/http"
	"regexp"
	"strconv"
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

// Minify is [finish this]
type Minify struct {
	Next     httpserver.Handler
	Excludes []string
	BasePath string
}

// ServeHTTP is the main function of the whole plugin that routes every single
// request to its function.
func (h Minify) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	if httpserver.Path(r.URL.Path).Matches(h.BasePath) {
		if isExcluded(strings.Replace(r.URL.Path, h.BasePath, "/", 1), h.Excludes) {
			return h.Next.ServeHTTP(w, r)
		}

		b := bytes.NewBuffer(nil)
		rw := &minifyResponseWriter{b, w}
		h.Next.ServeHTTP(rw, r)

		contentType, _, _ := mime.ParseMediaType(w.Header().Get("Content-Type"))

		if canBeMinified(contentType) {
			data, err := m.Bytes(contentType, b.Bytes())
			if err != nil {
				return 500, err
			}
			rw.Header().Set("Content-Length", strconv.Itoa(len(data)))
			w.Write(data)
			return 0, nil
		}

		rw.Header().Set("Content-Length", strconv.Itoa(len(b.Bytes())))
		w.Write(b.Bytes())
		return 0, nil

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
