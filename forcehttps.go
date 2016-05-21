package middleware

import (
	"net/http"
)

func ForceHTTPS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// normally the test is `!= "https"`, but we're permitting
		// Heroku-free dev-mode.
		if r.Header.Get("X-Forwarded-Proto") == "http" {
			u := *r.URL
			u.Scheme = "https"
			http.Redirect(w, r, u.String(), http.StatusMovedPermanently)
		}
	})
}
