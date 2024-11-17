package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/alexstory/otters"
)

func DefaultLogger() otters.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := &responseWriter{w, http.StatusOK}
			next.ServeHTTP(ww, r)
			duration := time.Since(start)
			fmt.Printf("%s %s %v %d\n", r.Method, r.URL.Path, duration, ww.statusCode)
		})
	}

}

// responseWriter is a custom http.ResponseWriter that captures the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
