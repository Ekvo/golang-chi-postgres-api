// source - database query
package source

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Ekvo/golang-chi-postgres-api/internal/model"
)

var ErrSourceNotFound = errors.New("not found")

// ErrSourceIncorrectData - incorrect (data any) passed to the function
var ErrSourceIncorrectData = errors.New("invalid data")

func (d *Dbinstance) CreateTables(ctx context.Context) error {
	_, err := d.db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS tasks
(
    id SERIAL PRIMARY KEY,
    description VARCHAR(2048) NOT NULL,
    note VARCHAR(2048) NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NULL
);`)
	return err
}

func (d *Dbinstance) SaveOneTask(ctx context.Context, data any) (uint, error) {
	newTask := data.(model.Task)
	taskNote := newTask.Note
	err := d.db.QueryRowContext(ctx, `
WITH new_task AS(
     INSERT INTO tasks(description,note,created_at)
     VALUES($1,$2,$3)
     RETURNING id)
SELECT id 
FROM new_task;`,
		newTask.Description,
		sql.NullString{taskNote, len(taskNote) != 0},
		newTask.CreatedAt,
	).Scan(&newTask.ID)
	return newTask.ID, err
}

func (d *Dbinstance) UpdateTask(ctx context.Context, data any) error {
	updateTask := data.(model.Task)
	taskNote := updateTask.Note
	taskID := 0
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	err = tx.QueryRowContext(ctx, `
UPDATE tasks
SET description = $2,
    note = $3,
    updated_at = $4
WHERE id = $1
RETURNING id;`,
		updateTask.ID,
		updateTask.Description,
		sql.NullString{taskNote, len(taskNote) != 0},
		updateTask.UpdatedAt,
	).Scan(&taskID)
	if err != nil || uint(taskID) != updateTask.ID {
		return ErrSourceNotFound
	}
	return tx.Commit()
}

func (d *Dbinstance) EndTaskLife(ctx context.Context, data any) error {
	taskID := data.(uint)
	delTaskID := 0
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	err = tx.QueryRowContext(ctx, `
DELETE 
FROM tasks
WHERE id = $1
RETURNING id;`, taskID).Scan(&delTaskID)
	if err != nil || taskID != uint(delTaskID) {
		return ErrSourceNotFound
	}
	return tx.Commit()
}

func (d *Dbinstance) FindOneTask(ctx context.Context, data any) (model.Task, error) {
	taskID := data.(uint)
	row := d.db.QueryRowContext(ctx, `
SELECT *
FROM tasks
WHERE id = $1
LIMIT 1;`, taskID)
	return scanOneTask[*sql.Row](row)
}

func (d *Dbinstance) FindTaskList(ctx context.Context, data any) ([]model.Task, error) {
	taskList := data.([]string)
	if len(taskList) != 3 {
		return nil, ErrSourceIncorrectData
	}
	query := fmt.Sprintf(`
SELECT * 
FROM tasks
ORDER BY id %s
LIMIT %s OFFSET %s;`,
		taskList[0], // desc or asc
		taskList[1],
		taskList[2],
	)
	rows, err := d.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	return scanTakList(rows)
}

// TaskScaner - for generic function 'scannerTask'
type RowScaner interface {
	*sql.Row | *sql.Rows
	Scan(dest ...any) error
}

func scanOneTask[S RowScaner](r S) (model.Task, error) {
	task := model.Task{}
	updatedAt := sql.NullTime{}
	note := sql.NullString{}
	if err := r.Scan(
		&task.ID,
		&task.Description,
		&note,
		&task.CreatedAt,
		&updatedAt,
	); err != nil {
		return task, ErrSourceNotFound
	}
	if note.Valid {
		task.Note = note.String
	}
	if updatedAt.Valid {
		task.UpdatedAt = &updatedAt.Time
	}
	return task, nil
}

func scanTakList(rows *sql.Rows) ([]model.Task, error) {
	var tasks []model.Task
	for rows.Next() {
		task, err := scanOneTask[*sql.Rows](rows)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, rows.Err()
}
