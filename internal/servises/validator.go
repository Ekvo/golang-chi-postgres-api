package servises

import (
	"net/http"
	"time"

	"github.com/Ekvo/golang-postgres-chi-api/internal/model"
	"github.com/Ekvo/golang-postgres-chi-api/pkg/common"
)

// TaskValidator - описывает свойсва получаемого объекта для создания 'models.Task'
type TaskValidator struct {
	Data struct {
		Description string `json:"description" binding:"required,min=1,max=2048"`
		Note        string `json:"note" binding:"omitempty,min=1,max=2048"`
	} `json:"task_update"`
	task model.Task `json:"-"`
}

// NewTaskValidator - если нужны 'Default' значения
func NewTaskValidator() *TaskValidator {
	return &TaskValidator{}
}

func (tv *TaskValidator) TaskModel() model.Task {
	return tv.task
}

// Bind - получение объекта и создание 'models.Task' на основе 'Data'
func (tv *TaskValidator) Bind(r *http.Request) error {
	if err := common.Bind(r, tv); err != nil {
		return err
	}
	tv.task.Description = tv.Data.Description
	tv.task.Note = tv.Data.Note
	tv.task.CreatedAt = time.Now().UTC()
	return nil
}
