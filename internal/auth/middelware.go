package auth

import (
	"net/http"
)

func CheckAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value("user").(*User)
		if !ok || user == nil {
			http.Error(w, "Login required", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
