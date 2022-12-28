package store

import (
	"database/sql"

	"github.com/lib/pq"
)

type Expense struct {
	ID     int      `json:"id"`
	Title  string   `json:"title"`
	Amount float64  `json:"amount"`
	Note   string   `json:"note"`
	Tags   []string `json:"tags"`
}

func (e *Expense) Create(db *sql.DB) error {
	result, err := db.Exec(`
		INSERT INTO expenses (title, amount, note, tags)
		VALUES ($1, $2, $3, $4)`, e.Title, e.Amount, e.Note, pq.Array(e.Tags))
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	e.ID = int(id)
	return err
}

func (e *Expense) Get(db *sql.DB) error {
	stm, err := db.Prepare("SELECT id, title, amount, note, tags FROM expenses WHERE id = $1")
	if err != nil {
		return err
	}
	row := stm.QueryRow(e.ID)

	return row.Scan(&e.ID, &e.Title, &e.Amount, &e.Note, pq.Array(&e.Tags))
}
