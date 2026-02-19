package middleware

import (
	"bytes"
	"net/http"

	"github.com/alexstory/otters"
)

const defaultServeMuxNotFoundBody = "404 page not found\n"

type bufferedResponseWriter struct {
	dst         http.ResponseWriter
	header      http.Header
	statusCode  int
	wroteHeader bool
	body        bytes.Buffer
}

func newBufferedResponseWriter(dst http.ResponseWriter) *bufferedResponseWriter {
	return &bufferedResponseWriter{
		dst:        dst,
		header:     make(http.Header),
		statusCode: http.StatusOK,
	}
}

func (w *bufferedResponseWriter) Header() http.Header {
	return w.header
}

func (w *bufferedResponseWriter) WriteHeader(code int) {
	if w.wroteHeader {
		return
	}

	w.wroteHeader = true
	w.statusCode = code
}

func (w *bufferedResponseWriter) Write(b []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	return w.body.Write(b)
}

func (w *bufferedResponseWriter) flush() error {
	for k, vv := range w.header {
		for _, v := range vv {
			w.dst.Header().Add(k, v)
		}
	}
	w.dst.WriteHeader(w.statusCode)
	_, err := w.dst.Write(w.body.Bytes())
	return err
}

// NotFound replaces the default net/http ServeMux 404 response with otters error pages.
//
// It only intercepts "uncaught" 404s (i.e. the default ServeMux not-found body),
// and will not override handlers that intentionally wrote their own 404 body.
func NotFound(app *otters.App) otters.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			bw := newBufferedResponseWriter(w)
			next.ServeHTTP(bw, r)

			isUncaught404 := bw.statusCode == http.StatusNotFound &&
				(bw.body.Len() == 0 || bw.body.String() == defaultServeMuxNotFoundBody)

			if isUncaught404 {
				ctx := app.NewCtx(w, r)
				_ = ctx.NotFound("not found")
				return
			}

			_ = bw.flush()
		})
	}
}
