package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/alexstory/otters"
)

// Recover catches panics from downstream handlers, logs a stack trace,
// and replies with HTTP 500
//
// If dev is true, it also writes the panic and stack to the response body.
func Recover(dev bool) otters.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				rec := recover()
				if rec == nil {
					return
				}

				stack := debug.Stack()
				fmt.Printf("panic recovered: %v\n%s\n", rec, stack)

				w.Header().Set("Content-Type", "text/plain; charset=utf-8")
				w.WriteHeader(http.StatusInternalServerError)

				if dev {
					_, _ = fmt.Fprintf(w, "panic: %v\n%s\n", rec, stack)
					return
				}

				_, _ = fmt.Fprintf(w, "internal server error")
			}()

			next.ServeHTTP(w, r)
		})
	}
}
