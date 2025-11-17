package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Mkhan2217/blocklist_app/config"
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
	// Build DSN from config
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost(),
		config.DBPort(),
		config.DBUser(),
		config.DBPassword(),
		config.DBName(),
	)

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

// CloseDB safely closes the database connection
func CloseDB() {
	if DB != nil {
		if err := DB.Close(); err != nil {
			log.Printf("⚠️  Error closing DB: %v", err)
		} else {
			log.Println("✅ Database connection closed")
		}
	}
}

// Schema initialization
func InitSchema() error {
	schema := `
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
`
	if _, err := DB.Exec(schema); err != nil {
		return fmt.Errorf("failed to initialize schema: %w", err)
	}

	log.Println("✅ Database schema initialized")
	return nil
}
