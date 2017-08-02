package minify

import (
	"bufio"
	"fmt"
	"io"
	"net"
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

// Hijack implements http.Hijacker. It simply wraps the underlying
// ResponseWriter's Hijack method if there is one, or returns an error.
func (w *minifyResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hj, ok := w.ResponseWriter.(http.Hijacker); ok {
		return hj.Hijack()
	}
	return nil, nil, fmt.Errorf("not a Hijacker")
}

// Flush implements http.Flusher. It simply wraps the underlying
// ResponseWriter's Flush method if there is one, or panics.
func (w *minifyResponseWriter) Flush() {
	if f, ok := w.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	} else {
		panic("not a Flusher") // should be recovered at the beginning of middleware stack
	}
}

// CloseNotify implements http.CloseNotifier.
// It just inherits the underlying ResponseWriter's CloseNotify method.
func (w *minifyResponseWriter) CloseNotify() <-chan bool {
	if cn, ok := w.ResponseWriter.(http.CloseNotifier); ok {
		return cn.CloseNotify()
	}
	panic("not a CloseNotifier")
}
