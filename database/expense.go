package database

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
	row := db.QueryRow(
		`INSERT INTO expenses (title, amount, note, tags) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id`,
		e.Title, e.Amount, e.Note, pq.Array(e.Tags),
	)
	return row.Scan(&e.ID)
}
