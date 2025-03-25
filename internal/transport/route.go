package transport

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/go-chi/chi/v5"

	m "github.com/Ekvo/golang-postgres-chi-api/internal/model"
	"github.com/Ekvo/golang-postgres-chi-api/internal/servises"
	s "github.com/Ekvo/golang-postgres-chi-api/internal/source"
	vr "github.com/Ekvo/golang-postgres-chi-api/internal/variables"
	"github.com/Ekvo/golang-postgres-chi-api/pkg/common"
)

// ErrTransportParam - ivalid param from 'chi.URLParam'
var ErrTransportParam = errors.New("invalid params")

func TaskCreate(db m.TaskUpdate) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		taskValidator := servises.NewTaskValidator()
		if err := taskValidator.Bind(r); err != nil {
			encodeJSON(ctx, w, http.StatusUnprocessableEntity, common.NewMessageErrorFromValidator(err))
			return
		}
		id, err := db.SaveOneTask(ctx, taskValidator.TaskModel())
		if err != nil {
			encodeJSON(ctx, w, http.StatusInternalServerError, common.NewMessageError(vr.DataBase, err))
			return
		}
		encodeJSON(ctx, w, http.StatusCreated, common.Message{vr.Task: id})
	}
}

func TaskUpdate(db m.TaskUpdate) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			encodeJSON(ctx, w, http.StatusBadRequest, common.NewMessageError(vr.Params, ErrTransportParam))
			return
		}
		taskValidator := servises.NewTaskValidator()
		if err := taskValidator.Bind(r); err != nil {
			encodeJSON(ctx, w, http.StatusUnprocessableEntity, common.NewMessageErrorFromValidator(err))
			return
		}
		task := taskValidator.TaskModel()
		task.ID = uint(id)
		task.UpdatedAt = &task.CreatedAt
		if err := db.UpdateTask(ctx, task); err != nil {
			encodeJSON(ctx, w, http.StatusNotFound, common.NewMessageError(vr.Task, s.ErrSourceNotFound))
			return
		}
		encodeJSON(ctx, w, http.StatusOK, common.Message{vr.Task: "updated"})
	}
}

func TaskRemove(db m.TaskUpdate) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			encodeJSON(ctx, w, http.StatusBadRequest, common.NewMessageError(vr.Params, ErrTransportParam))
			return
		}
		if err := db.EndTaskLife(ctx, uint(id)); err != nil {
			encodeJSON(ctx, w, http.StatusNotFound, common.NewMessageError(vr.Task, s.ErrSourceNotFound))
			return
		}
		encodeJSON(ctx, w, http.StatusOK, common.Message{vr.Task: "deleted"})
	}
}

func TaskByID(db m.TaskFind) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			encodeJSON(ctx, w, http.StatusBadRequest, common.NewMessageError(vr.Params, ErrTransportParam))
			return
		}
		task, err := db.FindOneTask(ctx, uint(id))
		if err != nil {
			encodeJSON(ctx, w, http.StatusNotFound, common.NewMessageError(vr.Task, s.ErrSourceNotFound))
			return
		}
		serializer := servises.TaskSerializer{task}
		encodeJSON(ctx, w, http.StatusOK, common.Message{vr.Task: serializer.Response()})
	}
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

func TaskList(db m.TaskFind) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		order := chi.URLParam(r, "order")
		limit := chi.URLParam(r, "limit")
		offset := chi.URLParam(r, "offset")
		var relimit = regexp.MustCompile(`^([1-9]|[1-9][0-9])+$`)
		var reoffset = regexp.MustCompile(`^[0-9]+$`)
		if !isValidOrder(order) || !relimit.MatchString(limit) || !reoffset.MatchString(offset) {
			encodeJSON(ctx, w, http.StatusBadRequest, common.NewMessageError(vr.Params, ErrTransportParam))
			return
		}
		tasks, err := db.FindTaskList(ctx, []string{order, limit, offset})
		if err != nil || len(tasks) == 0 {
			encodeJSON(ctx, w, http.StatusNoContent, common.NewMessageError(vr.DataBase, s.ErrSourceNotFound))
			return
		}
		serialize := servises.TaskListSerializer{tasks}
		encodeJSON(ctx, w, http.StatusOK, common.Message{vr.TaskList: serialize.Response()})
	}
}

// encodeJSON - we write the status and the object type of 'json' to 'ResponseWriter'
//
// if ctx.Err() == context.DeadlineExceeded - return us to 'func Timeout(timeout time.Duration)'
func encodeJSON(ctx context.Context, w http.ResponseWriter, status int, obj any) {
	if ctx.Err() == context.DeadlineExceeded {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(obj); err != nil {
		log.Printf("json.Encode error - %v", err)
	}
}
