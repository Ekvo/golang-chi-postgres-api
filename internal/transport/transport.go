// transport - describe handlers
package transport

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	m "github.com/Ekvo/golang-postgres-chi-api/internal/model"
)

// Transport - contain HTTP route multiplexer
type Transport struct {
	r *chi.Mux
}

func NewTransport(r *chi.Mux) *Transport {
	return &Transport{r: r}
}

// in pair with 'func Timeout(timeout time.Duration) func(next http.Handler) http.Handler'
const timeOut = 10 * time.Second

func (t *Transport) Routes(db m.TaskFindUpdate) {
	r := t.r
	r.Use(Timeout(timeOut))
	r.Mount("/task", taskRoutes(db))
}

func taskRoutes(db m.TaskFindUpdate) chi.Router {
	r := chi.NewRouter()
	r.Post("/", TaskCreate(db))
	r.Get("/{id}", TaskByID(db))
	r.Put("/{id}", TaskUpdate(db))
	r.Delete("/{id}", TaskRemove(db))
	r.Get("/{order}/{limit}/{offset}", TaskList(db))
	return r
}

// Timeout - middleware
// sets the query execution time use 'context'
func Timeout(timeout time.Duration) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer func() {
				if ctx.Err() == context.DeadlineExceeded {
					w.WriteHeader(http.StatusGatewayTimeout)
				}
				cancel()
			}()

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
