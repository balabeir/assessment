package database

import (
	"database/sql"
	"log"
	"os"

	"github.com/lib/pq"
)

type DB struct {
	db *sql.DB
}

func NewDB() *DB {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Connect database failed:", err)
	}

	return &DB{
		db: db,
	}
}

func InitialDB(db *sql.DB) {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS expenses (
		id SERIAL PRIMARY KEY,
		title TEXT,
		amount FLOAT,
		note TEXT,
		tags TEXT[]
	);`)
	if err != nil {
		log.Fatal(err)
	}
}

func (db *DB) CreateExpense(e *Expense) error {
	row := db.db.QueryRow(
		`INSERT INTO expenses (title, amount, note, tags) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id`,
		e.Title, e.Amount, e.Note, pq.Array(e.Tags),
	)
	return row.Scan(&e.ID)
}
