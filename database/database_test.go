package database

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func setupDB() *sql.DB {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Connect database failed:", err)
	}

	return db
}

func TestCreateExpense(t *testing.T) {
	db := setupDB()

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
