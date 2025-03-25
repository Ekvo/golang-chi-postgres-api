package servises

import (
	m "github.com/Ekvo/golang-postgres-chi-api/internal/model"
	vr "github.com/Ekvo/golang-postgres-chi-api/internal/variables"
)

// TaskSerializer - contains one "models.Task" to serialize into Response
type TaskSerializer struct {
	m.Task
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
		CreatedAt:   ts.CreatedAt.UTC().Format(vr.RFC3339Milli),
	}
	if ptrUpAt := ts.UpdatedAt; ptrUpAt != nil {
		tr.UpdatedAt = ptrUpAt.UTC().Format(vr.RFC3339Milli)
	}
	return tr
}

type TaskListSerializer struct {
	Tasks []m.Task
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
