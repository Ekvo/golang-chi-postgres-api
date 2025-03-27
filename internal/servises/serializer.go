package servises

import (
	"github.com/Ekvo/golang-chi-postgres-api/internal/model"
	"github.com/Ekvo/golang-chi-postgres-api/internal/variables"
)

// TaskSerializer - contains one "models.Task" to serialize into Response
type TaskSerializer struct {
	model.Task
}

// TaskResponse - format object 'Task' for 'Response'
type TaskResponse struct {
	Description string `json:"description"`
	Note        string `json:"note,omitempty"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at,omitempty"`
}

// (ts *TaskSerializer) Response() - returns an object to write to 'ResponseWriter'
func (ts *TaskSerializer) Response() TaskResponse {
	tr := TaskResponse{
		Description: ts.Description,
		Note:        ts.Note,
		CreatedAt:   ts.CreatedAt.UTC().Format(variables.RFC3339Milli),
	}
	if ptrUpAt := ts.UpdatedAt; ptrUpAt != nil {
		tr.UpdatedAt = ptrUpAt.UTC().Format(variables.RFC3339Milli)
	}
	return tr
}

type TaskListSerializer struct {
	Tasks []model.Task
}

func (tls *TaskListSerializer) Response() []TaskResponse {
	aliasTask := tls.Tasks
	n := len(aliasTask)
	tasksResponse := make([]TaskResponse, n)
	for i := 0; i < n; i++ {
		serialize := TaskSerializer{aliasTask[i]}
		tasksResponse[i] = serialize.Response()
	}
	return tasksResponse
}
