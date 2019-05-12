// Package minify provides middlewarefor minifying each request before being
// sent to the browser.
package minify

import (
	"net/http"

	"github.com/mholt/caddy/caddyhttp/httpserver"
	"github.com/tdewolff/minify/v2"
)

var minifier *minify.M

// Minify is an http.Handler that is able to minify the request before it's sent
// to the browser.
type Minify struct {
	Next  httpserver.Handler
	Rules []httpserver.RequestMatcher
	Paths []string
}

// ServeHTTP is the main function of the whole plugin that routes every single
// request to its function.
func (m Minify) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	// checks if the middlware should handle this request or not
	if m.shouldHandle(r) {
		mw := minifier.ResponseWriter(w, r)
		defer mw.Close()
		return m.Next.ServeHTTP(mw, r)
	}

	return m.Next.ServeHTTP(w, r)
}

// shouldHandle checks if the request should be handled with minifier
// using the BasePath and Excludes
func (m Minify) shouldHandle(r *http.Request) bool {
	included := false

	if len(m.Paths) > 0 {
		for _, path := range m.Paths {
			if httpserver.Path(r.URL.Path).Matches(path) {
				included = true
			}
		}
	} else {
		included = true
	}

	if !included {
		return false
	}

	for _, rule := range m.Rules {
		if !rule.Match(r) {
			return false
		}
	}

	return true
}
