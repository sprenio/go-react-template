package middleware

import (
	"net/http"

)
type statusResponseWriter struct {
	http.ResponseWriter
	status int
	wroteHeader bool
}

func (w *statusResponseWriter) WriteHeader(code int) {
	if w.wroteHeader {
		return
	}
	w.status = code
	w.wroteHeader = true
	w.ResponseWriter.WriteHeader(code)
}
func (w *statusResponseWriter) Write(b []byte) (int, error) {
	if !w.wroteHeader {
		// Jeśli handler nie wywołał WriteHeader explicite,
		// status = 200
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(b)
}
