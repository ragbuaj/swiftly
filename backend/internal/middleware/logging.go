package middleware

import (
	"log"
	"net/http"
	"time"
)

// LoggingMiddleware logs details of each HTTP request
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Create a wrapper for ResponseWriter to capture status code
		wrappedWriter := &responseWriter{ResponseWriter: w, status: http.StatusOK}
		
		next.ServeHTTP(wrappedWriter, r)
		
		log.Printf(
			"[%s] %s %s %d %s",
			r.Method,
			r.RequestURI,
			r.RemoteAddr,
			wrappedWriter.status,
			time.Since(start),
		)
	})
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
