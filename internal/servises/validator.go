package servises

import (
	"errors"
	"net/http"
	"time"

	"github.com/Ekvo/golang-chi-postgres-api/internal/model"
	"github.com/Ekvo/golang-chi-postgres-api/pkg/common"
)

var ErrservisesValidatorInvalidTask = errors.New("invalid task update")

// TaskValidator - describe property of getting and creating 'Task' object from a Request
type TaskValidator struct {
	Data struct {
		Description string `json:"description"`
		Note        string `json:"note"`
	} `json:"task_update"`
	task model.Task `json:"-"`
}

// NewTaskValidator - if need add 'Default' params
func NewTaskValidator() *TaskValidator {
	return &TaskValidator{}
}

func (tv *TaskValidator) TaskModel() model.Task {
	return tv.task
}

// DecodeJSON - get 'Data' and create 'Task'
func (tv *TaskValidator) DecodeJSON(r *http.Request) error {
	if err := common.DecodeJSON(r, tv); err != nil {
		return err
	}
	if tv.Data.Description == "" {
		return ErrservisesValidatorInvalidTask
	}
	tv.task.Description = tv.Data.Description
	tv.task.Note = tv.Data.Note
	tv.task.CreatedAt = time.Now().UTC()
	return nil
}
