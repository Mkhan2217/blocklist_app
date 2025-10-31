package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// package-level DB used by handlers
var db *sql.DB

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// ConnectDB opens the global DB connection using environment variables or defaults.
func ConnectDB() error {
	user := envOr("CHECKGUARD_DB_USER", "postgres")
	pass := envOr("CHECKGUARD_DB_PASSWORD", "rizzu")
	name := envOr("CHECKGUARD_DB_NAME", "blocklistdb")
	host := envOr("CHECKGUARD_DB_HOST", "localhost")
	port := envOr("CHECKGUARD_DB_PORT", "5432")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, pass, name)

	var err error
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		return err
	}

	// sensible pool defaults
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err = db.Ping(); err != nil {
		log.Printf("DB ping error: %v", err)
		return err
	}

	log.Println("Connected to DB")
	return nil
}

// CloseDB closes the global DB connection
func CloseDB() {
	if db != nil {
		_ = db.Close()
	}
}

// InitSchema ensures the required table and indexes exist
func InitSchema() error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS blocked_numbers (
		id SERIAL PRIMARY KEY,
		phone_number VARCHAR(15) UNIQUE NOT NULL CHECK (phone_number ~ '^\+[1-9][0-9]{9,14}$'),
		reason VARCHAR(100) NOT NULL,
		store_location VARCHAR(100) NOT NULL,
		incident_date DATE NOT NULL DEFAULT CURRENT_DATE,
		check_amount NUMERIC(12,2),
		notes VARCHAR(255),
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_phone_number ON blocked_numbers(phone_number);
	CREATE INDEX IF NOT EXISTS idx_store_location ON blocked_numbers(store_location);`

	_, err := db.Exec(createTableSQL)
	return err
}
