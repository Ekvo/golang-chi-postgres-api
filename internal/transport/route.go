package transport

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/Ekvo/golang-chi-postgres-api/internal/model"
	"github.com/Ekvo/golang-chi-postgres-api/internal/servises"
	"github.com/Ekvo/golang-chi-postgres-api/internal/source"
	vr "github.com/Ekvo/golang-chi-postgres-api/internal/variables"
	c "github.com/Ekvo/golang-chi-postgres-api/pkg/common"
)

// ErrTransportParam - ivalid param from 'chi.URLParam'
var ErrTransportParam = errors.New("invalid params")

func TaskCreate(db model.TaskUpdate) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		taskValidator := servises.NewTaskValidator()
		if err := taskValidator.DecodeJSON(r); err != nil {
			c.EncodeJSON(ctx, w, http.StatusUnprocessableEntity, c.NewMessageError(vr.Validator, err))
			return
		}
		id, err := db.SaveOneTask(ctx, taskValidator.TaskModel())
		if err != nil {
			c.EncodeJSON(ctx, w, http.StatusInternalServerError, c.NewMessageError(vr.DataBase, err))
			return
		}
		c.EncodeJSON(ctx, w, http.StatusCreated, c.Message{vr.Task: id})
	}
}

func TaskUpdate(db model.TaskUpdate) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			c.EncodeJSON(ctx, w, http.StatusBadRequest, c.NewMessageError(vr.Params, ErrTransportParam))
			return
		}
		taskValidator := servises.NewTaskValidator()
		if err := taskValidator.DecodeJSON(r); err != nil {
			c.EncodeJSON(ctx, w, http.StatusUnprocessableEntity, c.NewMessageError(vr.Validator, err))
			return
		}
		task := taskValidator.TaskModel()
		task.ID = uint(id)
		task.UpdatedAt = &task.CreatedAt
		if err := db.UpdateTask(ctx, task); err != nil {
			c.EncodeJSON(ctx, w, http.StatusNotFound, c.NewMessageError(vr.Task, source.ErrSourceNotFound))
			return
		}
		c.EncodeJSON(ctx, w, http.StatusOK, c.Message{vr.Task: "updated"})
	}
}

func TaskRemove(db model.TaskUpdate) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			c.EncodeJSON(ctx, w, http.StatusBadRequest, c.NewMessageError(vr.Params, ErrTransportParam))
			return
		}
		if err := db.EndTaskLife(ctx, uint(id)); err != nil {
			c.EncodeJSON(ctx, w, http.StatusNotFound, c.NewMessageError(vr.Task, source.ErrSourceNotFound))
			return
		}
		c.EncodeJSON(ctx, w, http.StatusOK, c.Message{vr.Task: "deleted"})
	}
}

func TaskByID(db model.TaskFind) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			c.EncodeJSON(ctx, w, http.StatusBadRequest, c.NewMessageError(vr.Params, ErrTransportParam))
			return
		}
		task, err := db.FindOneTask(ctx, uint(id))
		if err != nil {
			c.EncodeJSON(ctx, w, http.StatusNotFound, c.NewMessageError(vr.Task, source.ErrSourceNotFound))
			return
		}
		serializer := servises.TaskSerializer{task}
		c.EncodeJSON(ctx, w, http.StatusOK, c.Message{vr.Task: serializer.Response()})
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

func TaskList(db model.TaskFind) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		order := chi.URLParam(r, "order")
		limit := chi.URLParam(r, "limit")
		offset := chi.URLParam(r, "offset")
		var relimit = regexp.MustCompile(`^([1-9]|[1-9][0-9])+$`)
		var reoffset = regexp.MustCompile(`^[0-9]+$`)
		if !isValidOrder(order) || !relimit.MatchString(limit) || !reoffset.MatchString(offset) {
			c.EncodeJSON(ctx, w, http.StatusBadRequest, c.NewMessageError(vr.Params, ErrTransportParam))
			return
		}
		tasks, err := db.FindTaskList(ctx, []string{order, limit, offset})
		if err != nil || len(tasks) == 0 {
			c.EncodeJSON(ctx, w, http.StatusNoContent, c.NewMessageError(vr.DataBase, source.ErrSourceNotFound))
			return
		}
		serialize := servises.TaskListSerializer{tasks}
		c.EncodeJSON(ctx, w, http.StatusOK, c.Message{vr.TaskList: serialize.Response()})
	}
}
