// transport - сервер и маршрутизатор
package transport

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	m "github.com/Ekvo/golang-postgres-chi-api/internal/model"
)

// Transport - управление маршрутизацией
type Transport struct {
	r *chi.Mux
}

func NewTransport(r *chi.Mux) *Transport {
	return &Transport{r: r}
}

// в паре с 'func Timeout(timeout time.Duration) func(next http.Handler) http.Handler'
const timeOut = 10 * time.Second

func (t *Transport) Routes(db m.TaskFindUpdate) {
	router := t.r
	router.Use(Timeout(timeOut))
	router.Mount("/task", taskRoutes(db))
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

// Timeout - middleware функция
// задает время выполнени запроса через 'contxt'
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
