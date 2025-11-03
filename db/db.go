package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func ConnectDB() error {
	user := envOr("CHECKGUARD_DB_USER", "postgres")
	pass := envOr("CHECKGUARD_DB_PASSWORD", "rizzu")
	name := envOr("CHECKGUARD_DB_NAME", "blocklistdb")
	host := envOr("CHECKGUARD_DB_HOST", "localhost")
	port := envOr("CHECKGUARD_DB_PORT", "5432")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, pass, name)

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		return err
	}

	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)
	DB.SetConnMaxLifetime(5 * time.Minute)

	if err = DB.Ping(); err != nil {
		log.Println("DB ping error:", err)
		return err
	}

	log.Println("Connected to DB")
	return nil
}

func CloseDB() {
	if DB != nil {
		_ = DB.Close()
	}
}

// Schema initialization
func InitSchema() error {
	sql := `
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
	_, err := DB.Exec(sql)
	return err
}
