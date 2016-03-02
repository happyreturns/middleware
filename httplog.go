package middleware

import (
	"log"
	"net/http"
	"os"
	"path"
)

var logger = log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds)

type ResponseWriter interface {
	http.ResponseWriter
	Status() int
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func NewResponseWriter(rw http.ResponseWriter) ResponseWriter {
	return &responseWriter{rw, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(s int) {
	rw.status = s
	rw.ResponseWriter.WriteHeader(s)
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func Logger(handler http.Handler) http.Handler {
	app := path.Base(os.Args[0])
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := NewResponseWriter(w)
		handler.ServeHTTP(rw, r)
		logger.Printf("%s %s %s %d", app, r.Method, r.URL.Path, rw.Status())
	})
}
