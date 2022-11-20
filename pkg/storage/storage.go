package storage

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

var ctx context.Context = context.Background()

type Storage struct {
	db *pgxpool.Pool
}

type Tasks struct {
	id       int
	opened   int64
	closed   int64
	author   int
	assigned int
	title    string
	content  string
}

func New() (*Storage, error) {
	db, err := pgxpool.Connect(ctx, "postgres://postgres:password@192.168.92.128:5432/tasksBase")
	if err != nil {
		return nil, err
	}
	s := Storage{
		db: db,
	}
	defer db.Close()
	return &s, nil
}

func (s *Storage) AddTask(task Tasks) (int, error) {
	var id int
	err := s.db.QueryRow(
		ctx,
		`INSERT INTO tasks
			(opened, closed, author, assigned, title, content) 
		VALUES
			($1,$2,$3,$4,$5,$6)
		RETURNING ID;`,
		task.opened, task.closed, task.author, task.assigned, task.title, task.content,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *Storage) GetTasks(id, author_id, assigned_id int) ([]Tasks, error) {
	rows, err := s.db.Query(
		ctx,
		`SELECT
			id,
			opened,
			closed,
			author_id,
			assigned_id,
			title,
			content
		FROM tasks
		WHERE
			( id = $1 OR $1 = 0 ) AND ( author_id = $2 OR $2 = 0 ) AND ( assigned_id = $3 OR $3 = 0 ) 
		ORDER BY id;`,
		id,
		author_id,
		assigned_id,
	)
	if err != nil {
		return nil, err
	}
	var tasks []Tasks
	for rows.Next() {
		var task Tasks
		err = rows.Scan(
			&task.id,
			&task.opened,
			&task.closed,
			&task.author,
			&task.assigned,
			&task.title,
			&task.content,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, rows.Err()
}

func (s *Storage) UpdateTask(id int, task Tasks) error {
	_, err := s.db.Exec(
		ctx,
		`UPDATE tasks
		SET
			closed = $1,
			author_id = $2,
			assigned_id = $3,
			title = $4,
			content = $5
		WHERE
			id = $6;`,
		task.closed,
		task.author,
		task.assigned,
		task.title,
		task.content,
		id,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) DeleteTask(id int) error {
	_, err := s.db.Exec(
		ctx,
		`DELETE FROM tasks
		WHERE
			id = $1;`,
		id,
	)
	if err != nil {
		return err
	}
	return nil
}
