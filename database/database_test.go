//go:build unit
// +build unit

package database

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func setup(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	return db, mock
}

func TestCreateExpense(t *testing.T) {
	db, mock := setup(t)
	defer db.Close()

	expense := Expense{
		Title:  "john",
		Amount: 20,
		Note:   "test",
		Tags:   []string{"foo", "bar"},
	}

	mock.ExpectQuery("INSERT INTO expenses").
		WithArgs(expense.Title, expense.Amount, expense.Note, pq.Array(expense.Tags)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	err := expense.Create(db)
	if assert.NoError(t, err) {
		assert.NoError(t, mock.ExpectationsWereMet())
	}
}

func TestGetExpense(t *testing.T) {
	db, mock := setup(t)
	defer db.Close()

	mockRows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
		AddRow(1, "one", 10, "test1", pq.Array([]string{"foo", "bar"}))

	mock.ExpectPrepare("SELECT id, title, amount, note, tags FROM expenses").
		ExpectQuery().
		WithArgs(1).
		WillReturnRows(mockRows)

	expense := Expense{ID: 1}
	err := expense.Get(db)

	if assert.NoError(t, err) {
		assert.NoError(t, mock.ExpectationsWereMet())
	}
}

func TestUpdateExpense(t *testing.T) {
	db, mock := setup(t)
	defer db.Close()

	expense := Expense{
		Title:  "john",
		Amount: 20,
		Note:   "test",
		Tags:   []string{"foo", "bar"},
	}

	mock.ExpectExec("UPDATE expenses").
		WithArgs(expense.ID, expense.Title, expense.Amount, expense.Note, pq.Array(expense.Tags)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := expense.Update(db)
	if assert.NoError(t, err) {
		assert.NoError(t, mock.ExpectationsWereMet())
	}
}
