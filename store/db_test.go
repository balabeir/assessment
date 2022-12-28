package store

import (
	"database/sql"
	"regexp"
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

	mock.ExpectExec(regexp.QuoteMeta(`
			INSERT INTO expenses`)).
		WithArgs(expense.Title, expense.Amount, expense.Note, pq.Array(expense.Tags)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := expense.Create(db)
	assert.Nil(t, err)
}

func TestGetExpense(t *testing.T) {
	db, _ := setup(t)
	defer db.Close()

	expense := Expense{ID: 1}
	err := expense.Get(db)

	assert := assert.New(t)
	assert.Nil(err)
	assert.NotEqual(sql.ErrNoRows, err)
	assert.NotEqual(0, expense.ID)
}
