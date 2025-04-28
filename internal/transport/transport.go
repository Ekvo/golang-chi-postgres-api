// transport - describe handlers
package transport

import (
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/Ekvo/golang-chi-postgres-api/internal/model"
)

// Transport - contain HTTP route multiplexer
type Transport struct {
	*chi.Mux
}

func NewTransport(r *chi.Mux) *Transport {
	return &Transport{Mux: r}
}

// in pair with 'func Timeout(timeout time.Duration) func(next http.Handler) http.Handler'
const timeOut = 10 * time.Second

type taskFindUpdate interface {
	model.TaskFind
	model.TaskUpdate
}

func (r *Transport) Routes(db taskFindUpdate) {
	r.Use(Timeout(timeOut))
	r.Mount("/task", taskRoutes(db))
}

func taskRoutes(db taskFindUpdate) chi.Router {
	r := chi.NewRouter()
	r.Post("/", TaskHandler(db, taskCreate))
	r.Get("/{id}", TaskHandler(db, taskByID))
	r.Put("/{id}", TaskHandler(db, taskUpdate))
	r.Delete("/{id}", TaskHandler(db, taskRemove))
	r.Get("/{order}/{limit}/{offset}", TaskHandler(db, taskList))
	return r
}
