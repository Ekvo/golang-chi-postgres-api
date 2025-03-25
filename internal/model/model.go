// model - describe the Task object and its implementing interfaces
package model

import (
	"context"
	"time"
)

type Task struct {
	ID          uint
	Description string
	Note        string
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}

// TaskTables - create table of 'Task'
type TaskTables interface {
	CreateTables(ctx context.Context) error
}

// TaskUpdate - create update, dalete 'Task'
type TaskUpdate interface {
	SaveOneTask(ctx context.Context, data any) (uint, error)
	UpdateTask(ctx context.Context, data any) error
	EndTaskLife(ctx context.Context, data any) error
}

// TaskFind - find of 'Task', 'TaskList'
type TaskFind interface {
	FindOneTask(ctx context.Context, data any) (Task, error)
	FindTaskList(ctx context.Context, data any) ([]Task, error)
}

type TaskFindUpdate interface {
	TaskFind
	TaskUpdate
}
