// source - database query
package source

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/Ekvo/golang-chi-postgres-api/internal/model"
)

var (
	ErrSourceNotFound = errors.New("not found")

	// ErrSourceIncorrectData - incorrect (data any) passed to the function
	ErrSourceIncorrectData = errors.New("invalid data")
)

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
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer func() {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			log.Printf("query: insert task tx.Rollback error - %v", err)
		}
	}()
	err = d.db.QueryRowContext(ctx, `
INSERT INTO tasks(description,note,created_at)
VALUES($1,$2,$3)
RETURNING id;`,
		newTask.Description,
		emptyStringWriteNULL(newTask.Note),
		newTask.CreatedAt,
	).Scan(&newTask.ID)
	return newTask.ID, err
}

func (d *Dbinstance) UpdateTask(ctx context.Context, data any) error {
	updateTask := data.(model.Task)
	taskID := 0
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			log.Printf("query: update task tx.Rollback error - %v", err)
		}
	}()
	err = tx.QueryRowContext(ctx, `
UPDATE tasks
SET description = $2,
    note = $3,
    updated_at = $4
WHERE id = $1
RETURNING id;`,
		updateTask.ID,
		updateTask.Description,
		emptyStringWriteNULL(updateTask.Note),
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
	defer func() {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			log.Printf("query: delete task tx.Rollback error - %v", err)
		}
	}()
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
	limit, err := strconv.Atoi(taskList[1])
	if err != nil {
		return nil, ErrSourceIncorrectData
	}
	offset, err := strconv.Atoi(taskList[2])
	if err != nil {
		return nil, ErrSourceIncorrectData
	}
	query := strings.Builder{}
	args := make([]any, 0, 2)
	query.WriteString(`
SELECT * 
FROM tasks
ORDER BY id`)
	if taskList[0] == "desc" {
		query.WriteString(" DESC")
	}
	query.WriteString("\nLIMIT $1")
	args = append(args, limit)
	if offset > 0 {
		query.WriteString(" OFFSET $2")
		args = append(args, offset)
	}
	query.WriteByte(';')

	rows, err := d.db.QueryContext(ctx, query.String(), args...)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("query: rows.Close error - %v", err)
		}
	}()
	return scanTakList(rows)
}

func emptyStringWriteNULL(line string) *string {
	if line == "" {
		return nil
	}
	return &line
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
