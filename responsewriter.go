package minify

import (
	"io"
	"net/http"
)

// minifyResponseWriter wraps the ResponseWriter so it can save record
// the response to minify.
type minifyResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

// WriteHeader implements http.WriteHeader.
func (w *minifyResponseWriter) WriteHeader(code int) {
	w.Header().Del("Content-Length")
	w.ResponseWriter.WriteHeader(code)
}

// Write implements http.Write.
func (w *minifyResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}
