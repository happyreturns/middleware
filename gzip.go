package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func (w gzipResponseWriter) WriteHeader(code int) {
	w.Header().Del("Content-Length")
	w.ResponseWriter.WriteHeader(code)
}

func (w gzipResponseWriter) Flush() {
	if gz, ok := w.Writer.(http.Flusher); ok {
		gz.Flush()
	}
	if fw, ok := w.ResponseWriter.(http.Flusher); ok {
		fw.Flush()
	}
}

/*
Wraps an http.Handler and gzips the response if the request Accept-Encoding header contains gzip.

Example usage:
    http.ListenAndServe(":80", middleware.Gzip(router))
*/
func Gzip(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			handler.ServeHTTP(w, r)
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		gzw := gzipResponseWriter{gz, w}
		defer gz.Close()
		handler.ServeHTTP(gzw, r)
	})
}
