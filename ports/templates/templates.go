package templates

import (
	"context"
	"github.com/gorilla/csrf"
	"net/http"
)

const TemplateVersionKey = "TEMPLATES_VERSION"
const TemplateCSRFKey = "TEMPLATES_CSRF"

func GetVersion(ctx context.Context) string {
	if version, ok := ctx.Value(TemplateVersionKey).(string); ok {
		return version
	}
	return ""
}

func GetCSRF(ctx context.Context) string {
	if csrfValue, ok := ctx.Value(TemplateCSRFKey).(string); ok {
		return csrfValue
	}

	return ""
}

func equalsPointer(a *string, b string) bool {
	if a == nil {
		return false
	}

	if *a == b {
		return true
	}

	return false
}

func HttpMiddleware(next http.Handler, version string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(context.WithValue(r.Context(), TemplateVersionKey, version))
		r = r.WithContext(context.WithValue(r.Context(), TemplateCSRFKey, csrf.Token(r)))
		next.ServeHTTP(w, r)
	})
}
