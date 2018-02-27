// Package minify provides middlewarefor minifying each request before being
// sent to the browser.
package minify

import (
	"bytes"
	"log"
	"mime"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mholt/caddy/caddyhttp/httpserver"
	"github.com/tdewolff/minify"
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
		b := bytes.NewBuffer(nil)
		rw := &minifyResponseWriter{Writer: b, ResponseWriter: w}
		code, middlewareErr := m.Next.ServeHTTP(rw, r)

		// Bypass informational and redirection codes.
		if (code >= 100 && code < 200) || (code >= 300 && code < 400) {
			if length := b.Len(); length != 0 {
				rw.Header().Set("Content-Length", strconv.Itoa(length))
				w.Write(b.Bytes())
			}

			return code, middlewareErr
		}

		// gets the short version of Content-Type
		contentType, _, err := mime.ParseMediaType(w.Header().Get("Content-Type"))

		// If there is an error, log it
		// NOTE: this does not prevent the execution
		if err != nil {
			log.Println(err)
		}

		// If the content type is still blank, try getting it by the extension
		if contentType == "" {
			contentType = mime.TypeByExtension(filepath.Ext(r.URL.Path))
		}

		// If the content type is still blank and the Path ends with a /,
		// use the .html content type
		if contentType == "" && strings.HasSuffix(r.URL.Path, "/") {
			contentType = mime.TypeByExtension(".html")
		}

		var data []byte
		data, err = minifier.Bytes(contentType, b.Bytes())

		// Logs the error if it's not nil and different from Not Exist
		// NOTE: this does not prevent the execution
		if err != nil && err != minify.ErrNotExist {
			log.Println(err)
		}

		// Only send this header if the length is different from 0. It avoids
		// errors with Basic Auth and more stuff
		if length := len(data); length != 0 {
			rw.Header().Set("Content-Length", strconv.Itoa(length))
			w.Write(data)
		}

		return code, middlewareErr
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
