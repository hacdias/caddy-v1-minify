// Package minify provides middlewarefor minifying each request before being
// sent to the browser.
package minify

import (
	"bytes"
	"mime"
	"net/http"
	"regexp"
	"strconv"

	"github.com/mholt/caddy/caddyhttp/httpserver"
	"github.com/tdewolff/minify"
)

var (
	minifier  *minify.M
	jsonRegex = regexp.MustCompile("[/+]json$")
	xmlRegex  = regexp.MustCompile("[/+]xml$")
)

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
		code, err := m.Next.ServeHTTP(rw, r)

		// only handle if the status code is 200
		if code != http.StatusOK {
			return code, err
		}

		// gets the short version of Content-Type
		contentType, _, err := mime.ParseMediaType(w.Header().Get("Content-Type"))
		if err != nil {
			// this musn't happen
			return http.StatusInternalServerError, err
		}

		// if the contentType is blank, try getting it by the extension
		if contentType == "" {
			contentType = mime.TypeByExtension(r.URL.Path)
		}

		contentType = sanitizeContentType(contentType)

		if contentType != "" {
			var data []byte
			data, err = minifier.Bytes(contentType, b.Bytes())
			if err != nil {
				return http.StatusInternalServerError, err
			}
			rw.Header().Set("Content-Length", strconv.Itoa(len(data)))
			w.Write(data)
			return code, err
		}

		w.Header().Set("Content-Length", strconv.Itoa(len(b.Bytes())))
		w.Write(b.Bytes())
		return code, err
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

// sanitizeContentType simplifies the content type of the file to be used
// with the minifier
func sanitizeContentType(mimetype string) string {
	switch mimetype {
	case "text/css":
		return "css"
	case "​text/javascript",
		"text/x-javascript",
		"application/x-javascript",
		"application/javascript":
		return "javascript"
	case "image/svg+xml":
		return "svg"
	case "text/html":
		return "html"
	}

	if jsonRegex.FindString(mimetype) != "" {
		return "json"
	}

	if xmlRegex.FindString(mimetype) != "" {
		return "xml"
	}

	return ""
}
