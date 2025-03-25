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

	"github.com/Ekvo/golang-postgres-chi-api/internal/model"
	"github.com/Ekvo/golang-postgres-chi-api/internal/source"
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
	if oldTask, ex := m.tasks[newTask.ID]; !ex {
		return source.ErrSourceNotFound
	} else {
		// из- за ptr
		update := *newTask.UpdatedAt
		newTask.UpdatedAt = &update
		newTask.CreatedAt = oldTask.CreatedAt
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

// Проверка - HandlerFunc
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
		description:    "Task create",
		url:            "/task/",
		method:         http.MethodPost,
		bodyData:       `{"task_update":{"description":"Hello my friend!","note":"first task"}}`,
		expectedCode:   http.StatusCreated,
		responseRegexp: `{"task":([1-9]|[1-9][0-9])+}`,
		msg:            "valid - task is created -> get id with status 201",
	},
	{
		description:    "Wrong task create",
		url:            "/task/",
		method:         http.MethodPost,
		bodyData:       `{"task_update":{"note":"second task"}}`,
		expectedCode:   http.StatusUnprocessableEntity,
		responseRegexp: `{"errors":{"Description":"{required:Description}"}}`,
		msg:            "invalid - task not to be must created & status 422",
	},
	{
		description:    "Wrong task update 1 (method not allowed)",
		url:            "/task/1",
		method:         http.MethodPost,
		bodyData:       `{"task_update":{"description":"Hello, world!","note":"first task up"}}`,
		expectedCode:   http.StatusMethodNotAllowed,
		responseRegexp: ``,
		msg:            "invalid - task not to be mustupdated and status 405",
	},
	{
		description:    "Task update",
		url:            "/task/1",
		method:         http.MethodPut,
		bodyData:       `{"task_update":{"description":"Hello, world!","note":"first task up"}}`,
		expectedCode:   http.StatusOK,
		responseRegexp: `{"task":"updated"}`,
		msg:            "valid - task is must be updated and status 200",
	},
	{
		description:    "Wrong task update 2 (not found)",
		url:            "/task/200",
		method:         http.MethodPut,
		bodyData:       `{"task_update":{"description":"bad update","note":"zero"}}`,
		expectedCode:   http.StatusNotFound,
		responseRegexp: `{"errors":{"task":"not found"}}`,
		msg:            "invalid - task not to be must updated and status 404",
	},
	{
		description:  "Get task by id",
		url:          "/task/1",
		method:       http.MethodGet,
		bodyData:     ``,
		expectedCode: http.StatusOK,
		responseRegexp: `{
"task":{
"description":"Hello, world!",
"note":"first task up",
"created_at":"\d\d\d\d-(0[1-9]|1[0-2])-(0[1-9]|[12][0-9]|3[01])T([01][0-9]|2[0-3]):([0-5][0-9]):([0-5][0-9]).\d+(Z|[-+]([01][0-9]|2[0-3]):[0-5][0-9])",
"updated_at":"\d\d\d\d-(0[1-9]|1[0-2])-(0[1-9]|[12][0-9]|3[01])T([01][0-9]|2[0-3]):([0-5][0-9]):([0-5][0-9]).\d+(Z|[-+]([01][0-9]|2[0-3]):[0-5][0-9])"}
}`,
		msg: "valid - task must be found and status 200",
	},
	{
		description:    "Wrong get task by id",
		url:            "/task/alpha",
		method:         http.MethodGet,
		bodyData:       ``,
		expectedCode:   http.StatusBadRequest,
		responseRegexp: `{"errors":{"param":"invalid params"}}`,
		msg:            "invalid - must be contain errors and status 400",
	},
	{
		description:    "Delete task by id",
		url:            "/task/1",
		method:         http.MethodDelete,
		bodyData:       ``,
		expectedCode:   http.StatusOK,
		responseRegexp: `{"task":"deleted"}`,
		msg:            "valid - task must be deleted and status 200",
	},
	{
		description:    "Wrong delete task by id (not found)",
		url:            "/task/1",
		method:         http.MethodDelete,
		bodyData:       ``,
		expectedCode:   http.StatusNotFound,
		responseRegexp: `{"errors":{"task":"not found"}}`,
		msg:            "invalid - nothing deleted and status 404",
	},
}

func TestRoute(t *testing.T) {
	asserts := assert.New(t)
	requires := require.New(t)

	base := NewTasksMock()
	r := chi.NewRouter()
	h := NewTransport(r)
	h.Routes(base)

	for i, test := range routeTestData {
		log.Printf("\t %d test route: %s\n", i+1, test.description)
		bodyData := strings.Replace(test.bodyData, "\n", "", -1)

		method := test.method
		req, err := http.NewRequest(method, test.url, bytes.NewBuffer([]byte(bodyData)))
		requires.NoError(err, "http.NewRequest error")
		if method == http.MethodPost || method == http.MethodPut {
			req.Header.Set("Content-Type", "application/json")
		}
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		asserts.Equal(test.expectedCode, w.Code, test.msg)

		expected := strings.Replace(test.responseRegexp, "\n", "", -1)
		if expected != "" {
			requires.NotEmpty(w.Body, "response body no nil")
			asserts.Regexp(expected, w.Body.String(), test.msg)
		} else {
			asserts.Empty(w.Body, test.msg)
		}
	}
}

var orderTestData = []struct {
	order    string
	expected bool
	msg      string
}{
	{
		order:    asc,
		expected: true,
		msg:      "valid",
	},
	{
		order:    "some string",
		expected: false,
		msg:      "invalid",
	},
	{
		order:    desc,
		expected: true,
		msg:      "valid",
	},
}

func TestValidOrder(t *testing.T) {
	asserts := assert.New(t)

	for _, test := range orderTestData {
		result := isValidOrder(test.order)
		asserts.Equal(test.expected, result, test.msg)
	}
}
