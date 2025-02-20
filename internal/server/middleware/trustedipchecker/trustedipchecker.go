package trustedipchecker

import "net/http"

func TrustedIPChecker(trustedSubnet string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		senderIP := r.Header.Get("X-Real-IP")

		if trustedSubnet == "" {
			next.ServeHTTP(w, r)
		}

		if trustedSubnet != senderIP {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
