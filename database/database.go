package database

import (
	"database/sql"
	"log"
)

type DB struct {
	DB *sql.DB
}

func NewDB(db *sql.DB) *DB {
	return &DB{DB: db}
}

func InitialDB(db *DB) {
	_, err := db.DB.Exec(`
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
