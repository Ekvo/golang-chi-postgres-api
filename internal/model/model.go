// model - описывает объект Task и интерфейс
package model

import (
	"context"
	"time"
)

// объект описывающий таблицу в базе данных
type Task struct {
	ID          uint
	Description string
	Note        string
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}

// TaskTables - создание таблицы хранязей объект 'Task'
type TaskTables interface {
	CreateTables(ctx context.Context) error
}

// TaskUpdate - созданиеб обновлениеб удаление 'Task'
type TaskUpdate interface {
	SaveOneTask(ctx context.Context, data any) (uint, error)
	UpdateTask(ctx context.Context, data any) error
	EndTaskLife(ctx context.Context, data any) error
}

// TaskFind - поиск статьи или набора статей по переданному параметру 'data'
type TaskFind interface {
	FindOneTask(ctx context.Context, data any) (Task, error)
	FindTaskList(ctx context.Context, data any) ([]Task, error)
}

type TaskFindUpdate interface {
	TaskFind
	TaskUpdate
}
