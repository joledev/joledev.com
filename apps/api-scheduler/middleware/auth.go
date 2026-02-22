package middleware

import (
	"crypto/subtle"
	"net/http"
	"os"
)

func AdminAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		password := os.Getenv("SCHEDULER_ADMIN_PASSWORD")
		if password == "" {
			http.Error(w, `{"success":false,"message":"Admin not configured"}`, http.StatusInternalServerError)
			return
		}

		user, pass, ok := r.BasicAuth()
		if !ok || user != "admin" || subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="admin"`)
			http.Error(w, `{"success":false,"message":"Unauthorized"}`, http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
