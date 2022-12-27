package store

import (
	"database/sql"
	"log"
)

func Initial(db *sql.DB) {
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
