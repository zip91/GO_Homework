package storage

import (
	"database/sql"
	"go_course/Homework_5/internal/model"

	_ "github.com/lib/pq"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(connStr string) (*PostgresStore, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS tasks (
		id SERIAL PRIMARY KEY,
		uid TEXT NOT NULL,
		title TEXT NOT NULL,
		is_done BOOLEAN DEFAULT false
	)`)
	if err != nil {
		return nil, err
	}

	return &PostgresStore{db: db}, nil
}

func (s *PostgresStore) Create(task model.Task) error {
	_, err := s.db.Exec(`INSERT INTO tasks (uid, title, is_done) VALUES ($1, $2, $3)`, task.UID, task.Title, task.IsDone)
	return err
}

func (s *PostgresStore) GetByUID(uid string) ([]model.Task, error) {
	rows, err := s.db.Query(`SELECT id, uid, title, is_done FROM tasks WHERE uid = $1`, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []model.Task
	for rows.Next() {
		var t model.Task
		if err := rows.Scan(&t.ID, &t.UID, &t.Title, &t.IsDone); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}
