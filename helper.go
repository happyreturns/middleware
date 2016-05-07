package middleware

import (
	"net/http"
)

func Use(h http.Handler, wares ...func(http.Handler) http.Handler) http.Handler {
	for _, ware := range wares {
		h = ware(h)
	}
	return h
}
