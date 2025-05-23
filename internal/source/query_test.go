package source

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Ekvo/golang-chi-postgres-api/internal/config"
	"github.com/Ekvo/golang-chi-postgres-api/internal/model"
)

// to check when receiving a task from the database
var (
	timeCreate = time.Now().UTC()
	timeUpdate = timeCreate.Add(1 * time.Second)
)

func newValidTask() model.Task {
	return model.Task{
		ID:          0,
		Description: "Task for testing.",
		Note:        "must be deleted",
		CreatedAt:   timeCreate,
		UpdatedAt:   nil,
	}
}

func updateTask(id uint) model.Task {
	return model.Task{
		ID:          id,
		Description: "Task for testing is Update.",
		Note:        "must be deleted after this test",
		CreatedAt:   timeCreate,
		UpdatedAt:   &timeUpdate,
	}
}

var qq = []struct {
	description    string
	init           func(ctx context.Context, d *Dbinstance, data any) (any, error)
	ctxTimeOut     time.Duration
	data           any
	expectedResutl any
	haveErr        bool
	err            error
	msg            string
}{
	{
		description: "create table",
		init: func(ctx context.Context, d *Dbinstance, data any) (any, error) {
			return nil, d.CreateTables(ctx)
		},
		ctxTimeOut:     1 * time.Second,
		data:           nil,
		expectedResutl: nil,
		haveErr:        false,
		msg:            "success - table must be created",
	},
	{
		description: ("save task"),
		init: func(ctx context.Context, d *Dbinstance, data any) (any, error) {
			return d.SaveOneTask(ctx, data)
		},
		ctxTimeOut:     1 * time.Second,
		data:           newValidTask(),
		expectedResutl: uint(1),
		haveErr:        false,
		msg:            "success - task must be created",
	},
	{
		description: ("update task"),
		init: func(ctx context.Context, d *Dbinstance, data any) (any, error) {
			return nil, d.UpdateTask(ctx, data)
		},
		ctxTimeOut:     1 * time.Second,
		data:           updateTask(1),
		expectedResutl: nil,
		haveErr:        false,
		msg:            "success - task must be updated (descriptin,nore and updated_at)",
	},
	{
		description: ("wrong update task"),
		init: func(ctx context.Context, d *Dbinstance, data any) (any, error) {
			return nil, d.UpdateTask(ctx, data)
		},
		ctxTimeOut:     1 * time.Second,
		data:           updateTask(200),
		expectedResutl: nil,
		haveErr:        true,
		err:            ErrSourceNotFound,
		msg:            "wrong - task not to be must updated",
	},
	{
		description: ("find task"),
		init: func(ctx context.Context, d *Dbinstance, data any) (any, error) {
			return d.FindOneTask(ctx, data)
		},
		ctxTimeOut:     1 * time.Second,
		data:           uint(1),
		expectedResutl: updateTask(1),
		haveErr:        false,
		err:            nil,
		msg:            "valid - task must be return",
	},
	{
		description: ("wrong - find task"),
		init: func(ctx context.Context, d *Dbinstance, data any) (any, error) {
			return d.FindOneTask(ctx, data)
		},
		ctxTimeOut:     1 * time.Second,
		data:           uint(200),
		expectedResutl: model.Task{},
		haveErr:        true,
		err:            ErrSourceNotFound,
		msg:            "invalid - task must br return err and empty Task",
	},
	{
		description: ("delete task"),
		init: func(ctx context.Context, d *Dbinstance, data any) (any, error) {
			return nil, d.EndTaskLife(ctx, data)
		},
		ctxTimeOut:     1 * time.Second,
		data:           uint(1),
		expectedResutl: nil,
		haveErr:        false,
		err:            ErrSourceNotFound,
		msg:            "valid task must be deleted",
	},
	{
		description: ("wrong delete task"),
		init: func(ctx context.Context, d *Dbinstance, data any) (any, error) {
			return nil, d.EndTaskLife(ctx, data)
		},
		ctxTimeOut:     1 * time.Second,
		data:           uint(1),
		expectedResutl: nil,
		haveErr:        true,
		err:            ErrSourceNotFound,
		msg:            "Invalid - task cannot be deleted task does not exist",
	},
}

// connect for other test base 'postgres'
// table in base 'prometheus' - in save
func TestQueryDbinstance(t *testing.T) {
	asserts := assert.New(t)
	requires := require.New(t)

	cfg, err := config.NewConfig("../../.env", true)
	requires.NoError(err, fmt.Sprintf("query_test: config error - %v", err))

	db, err := Init(cfg)
	requires.NoError(err, fmt.Sprintf("query_test: db error - %v", err))
	defer db.Close()

	base := NewDbinstance(db)
	// for clear test
	_, err = db.Exec(`DROP TABLE tasks;`)
	requires.NoError(err, fmt.Sprintf("query_test: drop table error -%v", err))

	for i, query := range qq {
		log.Printf("\t %d query: %s\n", i+1, query.description)
		ctx, cancel := context.WithTimeout(context.Background(), query.ctxTimeOut)
		defer cancel()

		result, err := query.init(ctx, base, query.data)

		if err != nil {
			asserts.True(query.haveErr, query.msg)
			asserts.ErrorIs(err, query.err)
		} else {
			asserts.False(query.haveErr, query.msg)
		}

		switch expectedObj := query.expectedResutl.(type) {
		case model.Task:
			resultObj, ok := result.(model.Task)
			require.True(t, ok, "result must have type model.Task!!!")
			asserts.Equal(expectedObj.ID, resultObj.ID, "ID")
			asserts.Equal(expectedObj.Description, resultObj.Description, "Description")
			asserts.Equal(expectedObj.Note, resultObj.Note, "Note")
			asserts.Equal(expectedObj.CreatedAt.UTC().Format(time.RFC3339), resultObj.CreatedAt.UTC().Format(time.RFC3339), "created_at")
			if expectedObj.UpdatedAt != nil {
				asserts.Equal(expectedObj.UpdatedAt.UTC().Format(time.RFC3339), resultObj.UpdatedAt.UTC().Format(time.RFC3339), "updated_at")
			}
		default:
			asserts.Equal(query.expectedResutl, result, query.msg)
		}
	}

}
