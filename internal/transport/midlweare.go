// midlweare - contains middleware functions
package transport

import (
	"context"
	"net/http"
	"time"
)

// Timeout - middleware
// sets the query execution time use 'context'
func Timeout(timeout time.Duration) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer func() {
				if ctx.Err() == context.DeadlineExceeded {
					w.WriteHeader(http.StatusRequestTimeout)
				}
				cancel()
			}()

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
