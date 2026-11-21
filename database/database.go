package database

import (
    "database/sql"
    "log"

    _ "modernc.org/sqlite"
)

func ConnectDB() *sql.DB {
    db, err := sql.Open("sqlite", "./together.db")
    if err != nil {
        log.Fatal("Failed to open database:", err)
    }

    if err := db.Ping(); err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    query := `
    CREATE TABLE IF NOT EXISTS spaces (
        id TEXT PRIMARY KEY,
        code TEXT UNIQUE NOT NULL,
        is_active BOOLEAN NOT NULL DEFAULT TRUE,
        created_at DATETIME NOT NULL
    );

	CREATE TABLE IF NOT EXISTS members (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		space_id TEXT NOT NULL,
		is_host BOOLEAN NOT NULL DEFAULT FALSE,
		FOREIGN KEY (space_id) REFERENCES spaces(id)
	);
    `

    _, err = db.Exec(query)
    if err != nil {
        log.Fatal("Failed to create spaces table:", err)
    }

    log.Println("SQLite database connected successfully")

    return db
}