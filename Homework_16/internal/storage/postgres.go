package storage

import (
	"database/sql"

	_ "github.com/lib/pq"

	"go_course/Homework_5/internal/model"
)

// PostgresStore — хранилище поверх PostgreSQL
type PostgresStore struct {
	db *sql.DB
}

// NewPostgresStore открывает соединение и приводит схему в порядок
func NewPostgresStore(conn string) (*PostgresStore, error) {
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}
	p := &PostgresStore{db: db}
	if err := p.ensureSchema(); err != nil {
		_ = db.Close()
		return nil, err
	}
	return p, nil
}

func (p *PostgresStore) ensureSchema() error {
	const q = `
CREATE TABLE IF NOT EXISTS tasks (
  id BIGSERIAL PRIMARY KEY,
  uid TEXT NOT NULL,
  title TEXT NOT NULL,
  is_done BOOLEAN NOT NULL DEFAULT false
);
CREATE INDEX IF NOT EXISTS idx_tasks_uid ON tasks(uid);
`
	_, err := p.db.Exec(q)
	return err
}

// Create вставляет строку и возвращает id
func (p *PostgresStore) Create(t model.Task) (int64, error) {
	const q = `INSERT INTO tasks(uid, title, is_done) VALUES($1,$2,$3) RETURNING id`
	var id int64
	if err := p.db.QueryRow(q, t.UID, t.Title, t.IsDone).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

// GetByUID выбирает задачи пользователя
func (p *PostgresStore) GetByUID(uid string) ([]model.Task, error) {
	const q = `SELECT id, uid, title, is_done FROM tasks WHERE uid=$1 ORDER BY id`
	rows, err := p.db.Query(q, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []model.Task
	for rows.Next() {
		var t model.Task
		if err := rows.Scan(&t.ID, &t.UID, &t.Title, &t.IsDone); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}
