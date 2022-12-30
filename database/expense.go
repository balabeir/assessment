package database

import (
	"database/sql"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type Expense struct {
	ID     int      `json:"id"`
	Title  string   `json:"title"`
	Amount float64  `json:"amount"`
	Note   string   `json:"note"`
	Tags   []string `json:"tags"`
}

func (e *Expense) Create(db *sql.DB) error {
	row := db.QueryRow(`
		INSERT INTO expenses (id, title, amount, note, tags)
		VALUES (DEFAULT, $1, $2, $3, $4) RETURNING id`,
		e.Title, e.Amount, e.Note, pq.Array(e.Tags))
	err := row.Scan(&e.ID)
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
