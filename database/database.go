package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func Initial() *sql.DB {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Connect database failed:", err)
	}

	_, err = db.Exec(`
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

	return db
}
