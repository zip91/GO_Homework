package storage

import (
	"regexp"
	"testing"

	"go_course/Homework_5/internal/model"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestPostgresStore_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	s := &PostgresStore{db: db}

	mock.ExpectExec(regexp.QuoteMeta(
		`INSERT INTO tasks (uid, title, is_done) VALUES ($1, $2, $3)`)).
		WithArgs("u1", "Title", true).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = s.Create(model.Task{UID: "u1", Title: "Title", IsDone: true})
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresStore_GetByUID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	s := &PostgresStore{db: db}

	rows := sqlmock.NewRows([]string{"id", "uid", "title", "is_done"}).
		AddRow(1, "u1", "A", false).
		AddRow(2, "u1", "B", true)

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT id, uid, title, is_done FROM tasks WHERE uid = $1`
