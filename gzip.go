// per https://gist.github.com/bryfry/09a650eb8aac0fb76c24
// and https://github.com/NYTimes/gziphandler/blob/master/gzip.go

package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

// Gzip Compression
type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func (w gzipResponseWriter) Flush() {
	if gz, ok := w.Writer.(http.Flusher); ok {
		gz.Flush()
	}
	if fw, ok := w.ResponseWriter.(http.Flusher); ok {
		fw.Flush()
	}
}

func Gzip(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			handler.ServeHTTP(w, r)
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		gzw := gzipResponseWriter{gz, w}
		handler.ServeHTTP(gzw, r)
	})
}
