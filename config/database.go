package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// DB instance
var db *sql.DB

// Initialize Database Connection
func InitDB() {
	var err error
	dbPath := ConfigValues[DB_PATH]
	if dbPath == "" {
		dbPath = "/home/shofiya/abr/abr.db" // Default for local development
	}

	// Open SQLite database
	db, err = sql.Open("sqlite3", "file:"+dbPath+"?_journal_mode=WAL&_cache=shared")

	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	createTables()
}

func GetDB() *sql.DB {
	if db == nil {
		InitDB()
	}
	return db
}

// Create tables if they don't exist
func createTables() {

	createUsersTable := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		client_id TEXT UNIQUE NOT NULL,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL
	)`

	createVideosTable := `
	CREATE TABLE IF NOT EXISTS videos (
	    video_id TEXT PRIMARY KEY,
	    client_id TEXT NOT NULL,
	    upload_time DATETIME DEFAULT CURRENT_TIMESTAMP,
	    status TEXT CHECK(status IN ('initialized', 'transcoding', 'completed', 'failed')) NOT NULL,
	    file_key TEXT NOT NULL,
	    bucket TEXT NOT NULL,
	    strategy TEXT CHECK(strategy IN ('single', 'multipart')) NOT NULL,
	    FOREIGN KEY (client_id) REFERENCES clients(client_id) ON DELETE CASCADE
	);`

	_, err := db.Exec(createUsersTable)
	if err != nil {
		log.Fatal("Error creating table:", err)
	}

	_, err = db.Exec(createVideosTable)
	if err != nil {
		log.Fatal("Error creating table:", err)
	}

	fmt.Println("Database initialized successfully.")
}
