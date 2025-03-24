package transport

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Ekvo/golang-postgres-chi-docker-api/internal/model"
	"github.com/Ekvo/golang-postgres-chi-docker-api/internal/source"
)

type TasksMock struct {
	nextID uint
	tasks  map[uint]model.Task
}

func NewTasksMock() *TasksMock {
	return &TasksMock{
		nextID: 1,
		tasks:  map[uint]model.Task{},
	}
}

func (m *TasksMock) SaveOneTask(ctx context.Context, data any) (uint, error) {
	newTask := data.(model.Task)
	newTask.ID = m.nextID
	m.tasks[newTask.ID] = newTask
	m.nextID++
	return newTask.ID, ctx.Err()
}

func (m *TasksMock) UpdateTask(ctx context.Context, data any) error {
	newTask := data.(model.Task)
	if _, ex := m.tasks[newTask.ID]; !ex {
		return source.ErrSourceNotFound
	}
	m.tasks[newTask.ID] = newTask
	return ctx.Err()
}

func (m *TasksMock) EndTaskLife(ctx context.Context, data any) error {
	taskId := data.(uint)
	if _, ex := m.tasks[taskId]; !ex {
		return source.ErrSourceNotFound
	}
	delete(m.tasks, taskId)
	return ctx.Err()
}

func (m *TasksMock) FindOneTask(ctx context.Context, data any) (model.Task, error) {
	taskId := data.(uint)
	if task, ex := m.tasks[taskId]; !ex {
		return model.Task{}, source.ErrSourceNotFound
	} else {
		return task, ctx.Err()
	}
}

func (m *TasksMock) FindTaskList(ctx context.Context, data any) ([]model.Task, error) {
	arrID := make([]uint, 0, len(m.tasks))
	for id, _ := range m.tasks {
		arrID = append(arrID, id)
	}
	taskList := data.([]string)

	//ORDER BY (ASC OR DESC)
	var fn func(i, j int) bool
	if taskList[0] == asc {
		fn = func(i, j int) bool { return arrID[i] < arrID[j] }
	} else {
		fn = func(i, j int) bool { return arrID[i] > arrID[j] }
	}
	sort.Slice(arrID, fn)
	offset, _ := strconv.Atoi(taskList[2])
	if offset >= len(arrID) {
		return nil, ctx.Err()
	}
	arrID = arrID[offset:]
	limit, _ := strconv.Atoi(taskList[1])
	if limit < len(arrID) {
		arrID = arrID[:limit]
	}
	tasks := make([]model.Task, 0, len(arrID))
	for _, id := range arrID {
		tasks = append(tasks, m.tasks[id])
	}
	return tasks, ctx.Err()
}

var routeTestData = []struct {
	description    string
	url            string
	method         string
	bodyData       string
	expectedCode   int
	responseRegexp string
	msg            string
}{
	{
		description:    "Task Create",
		url:            "/task/",
		method:         http.MethodPost,
		bodyData:       `{"task_update":{"description":"Hello ny friend!","note":"second task"}}`,
		expectedCode:   http.StatusCreated,
		responseRegexp: `{"task":([1-9]|[1-9][0-9])+}`,
		msg:            "valid - task is created -> get id",
	},
}

func TestRoute(t *testing.T) {
	asserts := assert.New(t)

	base := NewTasksMock()
	r := chi.NewRouter()
	h := NewTransport(r)
	h.Routes(base)

	for i, test := range routeTestData {
		log.Printf("\t %d test route: %s\n", i+1, test.description)
		bodyData := strings.Replace(test.bodyData, "\n", "", -1)

		method := test.method
		req, err := http.NewRequest(method, test.url, bytes.NewBuffer([]byte(bodyData)))
		require.NoError(t, err, "http.NewRequest error")
		if method == http.MethodPost || method == http.MethodPut {
			req.Header.Set("Content-Type", "application/json")
		}

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		asserts.Equal(test.expectedCode, w.Code, test.msg)
		require.NotNil(t, w.Body, "response body no nil")
		asserts.Regexp(test.responseRegexp, w.Body.String(), test.msg)
	}
}
