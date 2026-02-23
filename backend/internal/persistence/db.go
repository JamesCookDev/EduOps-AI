package persistence

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	createTable := `
	CREATE TABLE IF NOT EXISTS container_metrics (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		container_id TEXT,
		name TEXT,
		cpu_percent REAL,
		memory_usage INTEGER,
		memory_limit INTEGER,
		memory_percent REAL,
		restart_count INTEGER,
		timestamp DATETIME
	);
	`

	_, err = db.Exec(createTable)
	if err != nil {
		return nil, err
	}

	log.Println("Banco inicializado")

	return db, nil
}