package middleware

import (
	"log"
	"net/http"
	"os"
	"path"
)

var logger = log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds)

type statusResponseWriter interface {
	http.ResponseWriter
	Status() int
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func newStatusResponseWriter(rw http.ResponseWriter) statusResponseWriter {
	return &responseWriter{rw, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(s int) {
	rw.status = s
	rw.ResponseWriter.WriteHeader(s)
}

func (rw *responseWriter) Status() int {
	return rw.status
}

/*
Logger prints to stdout the method, url, and response status of requests.

Example usage:
   http.ListenAndServe(":80", middleware.Logger(router))
*/
func Logger(handler http.Handler) http.Handler {
	app := path.Base(os.Args[0])
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := newStatusResponseWriter(w)
		handler.ServeHTTP(rw, r)
		logger.Printf("%s %s %s %s %d", app, r.Header.Get("requestID"), r.Method, r.URL.RequestURI(), rw.Status())
	})
}
