package store

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

	mock.ExpectExec("INSERT INTO expenses").
		WithArgs(expense.Title, expense.Amount, expense.Note, pq.Array(expense.Tags)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := expense.Create(db)
	assert.NoError(t, err)
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

	assert.NoError(t, err)
}
