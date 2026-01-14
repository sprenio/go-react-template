package middleware

import (
	"backend/internal/contexthelper"
	"net"
	"net/http"
	"strings"
)


func IP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		ip := getClientIP(r)
		ctx = contexthelper.SetClientIp(ctx, ip)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}


func getClientIP(r *http.Request) string {
	// 1. X-Forwarded-For (może zawierać kilka IP)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// pierwszy = prawdziwy klient
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}

	// 2. X-Real-IP
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}

	// 3. RemoteAddr (fallback)
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return ip
	}

	return r.RemoteAddr
}
