package authcheck

import (
	"github.com/happyreturns/getenv"
	"io/ioutil"
	"net/http"
)

var authURL = getenv.Get("AUTH_URL", "http://localhost:8001")

func Handler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			auth = r.FormValue("Authorization")
		}
		if auth == "" {
			code := http.StatusUnauthorized
			http.Error(w, http.StatusText(code), code)
			return
		}

		// Check with the auth server if this is ok
		req, err := http.NewRequest("GET", authURL+"/", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		req.Header.Set("Authorization", auth)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			http.Error(w, string(body), resp.StatusCode)
			return
		}

		handler.ServeHTTP(w, r)
	})
}
