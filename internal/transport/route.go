package transport

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/Ekvo/golang-chi-postgres-api/internal/servises"
	"github.com/Ekvo/golang-chi-postgres-api/internal/source"
	vr "github.com/Ekvo/golang-chi-postgres-api/internal/variables"
	c "github.com/Ekvo/golang-chi-postgres-api/pkg/common"
)

// ErrTransportParam - ivalid param from 'chi.URLParam'
var ErrTransportParam = errors.New("invalid params")

// responseData - contain data for 'http.ResponseWriter'
type responseData struct {
	status int
	body   any
}

// taskFunc - layout of function for TashHandler
type taskFunc func(db taskFindUpdate, r *http.Request) responseData

// TaskHandler - main function on route(work with Timeout see ./middlweare.go)
//
// call in goroutines 'taskFn' for get 'responseData' to chan 'response'
// in 'select' checks execution time and create body for 'http.ResponseWriter'
func TaskHandler(db taskFindUpdate, taskFn taskFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		response := make(chan responseData)

		go func() {
			defer close(response)
			select {
			case <-ctx.Done():
				return
			case response <- taskFn(db, r):
			}
		}()

		select {
		case <-ctx.Done():
			return
		case responseData := <-response:
			c.EncodeJSON(w, responseData.status, responseData.body)
		}
	}
}

func taskCreate(db taskFindUpdate, r *http.Request) responseData {
	taskValidator := servises.NewTaskValidator()
	if err := taskValidator.DecodeJSON(r); err != nil {
		return responseData{http.StatusUnprocessableEntity, c.NewMessageError(vr.Validator, err)}
	}
	id, err := db.SaveOneTask(r.Context(), taskValidator.TaskModel())
	if err != nil {
		return responseData{http.StatusInternalServerError, c.NewMessageError(vr.DataBase, err)}
	}
	return responseData{http.StatusCreated, c.Message{vr.Task: id}}
}

func taskUpdate(db taskFindUpdate, r *http.Request) responseData {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return responseData{http.StatusBadRequest, c.NewMessageError(vr.Params, ErrTransportParam)}
	}
	taskValidator := servises.NewTaskValidator()
	if err := taskValidator.DecodeJSON(r); err != nil {
		return responseData{http.StatusUnprocessableEntity, c.NewMessageError(vr.Validator, err)}
	}
	task := taskValidator.TaskModel()
	task.ID = uint(id)
	task.UpdatedAt = &task.CreatedAt
	if err := db.UpdateTask(r.Context(), task); err != nil {
		return responseData{http.StatusNotFound, c.NewMessageError(vr.Task, source.ErrSourceNotFound)}
	}
	return responseData{http.StatusOK, c.Message{vr.Task: "updated"}}
}

func taskRemove(db taskFindUpdate, r *http.Request) responseData {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return responseData{http.StatusBadRequest, c.NewMessageError(vr.Params, ErrTransportParam)}
	}
	if err := db.EndTaskLife(r.Context(), uint(id)); err != nil {
		return responseData{http.StatusNotFound, c.NewMessageError(vr.Task, source.ErrSourceNotFound)}
	}
	return responseData{http.StatusOK, c.Message{vr.Task: "deleted"}}
}

func taskByID(db taskFindUpdate, r *http.Request) responseData {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return responseData{http.StatusBadRequest, c.NewMessageError(vr.Params, ErrTransportParam)}
	}
	task, err := db.FindOneTask(r.Context(), uint(id))
	if err != nil {
		return responseData{http.StatusNotFound, c.NewMessageError(vr.Task, source.ErrSourceNotFound)}
	}
	serializer := servises.TaskSerializer{Task: task}
	return responseData{http.StatusOK, c.Message{vr.Task: serializer.Response()}}
}

// param from request 't.r.Post("/tasks/{order}/{limit}/{offset}")'
// describes 'ORDER BY in PostgresSQL'
const (
	asc  = "asc"
	desc = "desc"
)

func isValidOrder(order string) bool {
	return order == asc || order == desc
}

var (
	// relimit - rules for use LIMIT in SQL query
	relimit = regexp.MustCompile(`^([1-9]|[1-9][0-9])+$`)

	// reoffset - OFFSET in SQL query
	reoffset = regexp.MustCompile(`^[0-9]+$`)
)

// taskList - get (order, limit, offset) for use in SQL query
func taskList(db taskFindUpdate, r *http.Request) responseData {
	order := chi.URLParam(r, "order")
	limit := chi.URLParam(r, "limit")
	offset := chi.URLParam(r, "offset")
	if !isValidOrder(order) ||
		!relimit.MatchString(limit) ||
		!reoffset.MatchString(offset) {
		return responseData{http.StatusBadRequest, c.NewMessageError(vr.Params, ErrTransportParam)}
	}
	tasks, err := db.FindTaskList(r.Context(), []string{order, limit, offset})
	if err != nil || len(tasks) == 0 {
		return responseData{http.StatusNoContent, c.NewMessageError(vr.DataBase, source.ErrSourceNotFound)}
	}
	serialize := servises.TaskListSerializer{Tasks: tasks}
	return responseData{http.StatusOK, c.Message{vr.TaskList: serialize.Response()}}
}
