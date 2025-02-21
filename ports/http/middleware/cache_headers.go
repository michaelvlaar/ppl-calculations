package middleware

import (
	"net/http"
	"os"
)

func CacheHeaders(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if os.Getenv("CACHE_DISABLED") != "true" {
			w.Header().Set("Cache-Control", "public, max-age=86400")
		}
		h.ServeHTTP(w, r)
	})
}

func CacheHeadersFunc(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if os.Getenv("CACHE_DISABLED") != "true" {
			w.Header().Set("Cache-Control", "public, max-age=86400")
		}
		h.ServeHTTP(w, r)
	}
}
