package middleware

import (
	"net/http"
)

func ForceHTTPS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// normally the test is `!= "https"`, but we're permitting
		// Heroku-free dev-mode.
		if r.Header.Get("X-Forwarded-Proto") == "http" {
			u := "https:" + r.URL.Path
			if r.URL.RawQuery != "" {
				u = u + "?" + r.URL.RawQuery
			}
			http.Redirect(w, r, u, http.StatusMovedPermanently)
		} else {
			h.ServeHTTP(w, r)
		}
	})
}
