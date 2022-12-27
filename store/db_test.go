package store

import (
	"database/sql"
	"os"
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func setup(t *testing.T) *sql.DB {
	db, _ := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	return db
}

func TestCreateExpense(t *testing.T) {
	db := setup(t)
	defer db.Close()

	expect := Expense{
		Title:  "john",
		Amount: 20,
		Note:   "test",
		Tags:   []string{"foo", "bar"},
	}
	expect.Create(db)

	got := Expense{}
	stm, _ := db.Prepare("SELECT * FROM expenses WHERE id = $1")
	err := stm.QueryRow(expect.ID).Scan(&got.ID, &got.Title, &got.Amount, &got.Note, pq.Array(&got.Tags))

	assert := assert.New(t)
	assert.NotEqual(sql.ErrNoRows, err)
	assert.Nil(err)
	assert.Equal(expect.ID, got.ID)
	assert.Equal(expect.Title, got.Title)
	assert.Equal(expect.Amount, got.Amount)
	assert.Equal(expect.Note, got.Note)
	assert.Equal(len(expect.Tags), len(got.Tags))
}
