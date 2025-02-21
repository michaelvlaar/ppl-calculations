package middleware

import (
	"github.com/gorilla/csrf"
	"net/http"
	"os"
)

func CSRF(next http.Handler) http.Handler {
	cookieName := "csrf"
	if os.Getenv("SECURE_COOKIE") == "true" {
		cookieName = "__Secure-" + cookieName
	}

	return csrf.Protect([]byte(os.Getenv("CSRF_KEY")), csrf.CookieName(cookieName), csrf.SameSite(csrf.SameSiteStrictMode), csrf.Path("/"), csrf.HttpOnly(true))(next)
}
