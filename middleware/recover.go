package middleware

import (
	"fmt"
	"html"
	"net/http"
	"runtime/debug"
	"time"

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

				w.WriteHeader(http.StatusInternalServerError)

				if dev {
					w.Header().Set("Content-Type", "text/html; charset=utf-8")
					title := "otters: panic recovered"
					reqLine := fmt.Sprintf("%s %s", r.Method, r.URL.Path)

					panicStr := html.EscapeString(fmt.Sprintf("%v", rec))
					stackStr := html.EscapeString(string(stack))
					reqLineStr := html.EscapeString(reqLine)
					whenStr := html.EscapeString(time.Now().Format(time.RFC3339))
					_, _ = fmt.Fprintf(w, `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>%s</title>
  <style>
    :root { color-scheme: light dark; }
    body { font-family: ui-sans-serif, system-ui, -apple-system, Segoe UI, Roboto, Helvetica, Arial, "Apple Color Emoji", "Segoe UI Emoji"; margin: 0; padding: 0; }
    header { padding: 16px 20px; border-bottom: 1px solid rgba(127,127,127,.35); }
    main { padding: 20px; max-width: 1000px; margin: 0 auto; }
    .meta { opacity: .8; font-size: 14px; margin-top: 6px; }
    .card { margin-top: 16px; padding: 14px 16px; border: 1px solid rgba(127,127,127,.35); border-radius: 10px; }
    h1 { font-size: 18px; margin: 0; }
    h2 { font-size: 14px; margin: 0 0 10px 0; opacity: .9; }
    pre { white-space: pre-wrap; word-break: break-word; margin: 0; font-size: 12px; line-height: 1.45; }
    code { font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace; }
  </style>
</head>
<body>
  <header>
    <h1>%s</h1>
    <div class="meta"><code>%s</code> · %s</div>
  </header>
  <main>
    <section class="card">
      <h2>Panic</h2>
      <pre><code>%s</code></pre>
    </section>
    <section class="card">
      <h2>Stack</h2>
      <pre><code>%s</code></pre>
    </section>
  </main>
</body>
</html>`, title, title, reqLineStr, whenStr, panicStr, stackStr)
					return
				}

				w.Header().Set("Content-Type", "text/plain; charset=utf-8")
				_, _ = fmt.Fprintf(w, "internal server error")
			}()

			next.ServeHTTP(w, r)
		})
	}
}
