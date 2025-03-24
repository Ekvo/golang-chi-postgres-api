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

// ErrTransportParam - некоректный параметр 'chi.URLParam'
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
		if err := db.UpdateTask(ctx, task); err != nil {
			encodeJSON(ctx, w, http.StatusNotFound, common.NewMessageError(vr.Task, s.ErrSourceNotFound))
			return
		}
		encodeJSON(ctx, w, http.StatusOK, common.Message{vr.Task: "update"})
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
			encodeJSON(ctx, w, http.StatusBadRequest, ErrTransportParam)
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

// формирования запроса 't.r.Post("/tasks/{order}/{limit}/{offset}")'
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
			encodeJSON(ctx, w, http.StatusNoContent, common.NewMessageError(vr.DataBase, err))
			return
		}
		serialize := servises.TaskListSerializer{tasks}
		encodeJSON(ctx, w, http.StatusOK, common.Message{vr.TaskList: serialize.Response()})
	}
}

// encodeJSON - записываем статус и объект типа 'json' в 'ResponseWriter'
//
// if ctx.Err() == context.DeadlineExceeded - возвращаемся в  'func Timeout(timeout time.Duration)'
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
