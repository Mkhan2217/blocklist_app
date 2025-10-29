package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	_ "github.com/lib/pq" // Re-added import for PostgreSQL driver
)

//  Database connection settings

const (
	DB_USER     = "postgres"
	DB_PASSWORD = "rizzu"
	DB_NAME     = "blocklistdb"
	DB_HOST     = "localhost"
	DB_PORT     = "5432"
)

// ------------------------------
//
//	Struct to represent a phone record
//
// ------------------------------
type BlockedNumber struct {
	ID            int     `json:"id"`
	PhoneNumber   string  `json:"phone_number"`
	Reason        string  `json:"reason"`
	StoreLocation string  `json:"store_location"`
	IncidentDate  string  `json:"incident_date"`
	CreatedAt     string  `json:"created_at"`
	CheckAmount   float64 `json:"check_amount,omitempty"`
	Notes         string  `json:"notes,omitempty"`
}

// ------------------------------
//
//	Global DB connection variable
//
// ------------------------------
var db *sql.DB

// Update phone validation regex to match database constraint
var phoneRegex = regexp.MustCompile(`^\+[1-9][0-9]{9,14}$`)

func validatePhoneNumber(phone string) bool {
	return phoneRegex.MatchString(phone)
}

func formatPhoneNumber(phone string) string {
	// URL decode the phone number
	decoded, err := url.QueryUnescape(phone)
	if err != nil {
		log.Printf("URL decode error: %v, using original: %s", err, phone)
		decoded = phone
	}

	// Remove all non-digit characters except '+'
	reg := regexp.MustCompile(`[^\d+]`)
	cleaned := reg.ReplaceAllString(decoded, "")

	// Only add + if missing, don't modify the number otherwise
	if !strings.HasPrefix(cleaned, "+") {
		cleaned = "+" + cleaned
	}

	log.Printf("Formatted phone: %s", cleaned)
	return cleaned
}

// ------------------------------
//
//	Home handler: display form + list
//
// ------------------------------
func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	rows, err := db.Query("SELECT id, phone_number, created_at FROM blocked_numbers ORDER BY id DESC")
	if err != nil {
		http.Error(w, "DB query error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var numbers []BlockedNumber
	for rows.Next() {
		var n BlockedNumber
		rows.Scan(&n.ID, &n.PhoneNumber, &n.CreatedAt)
		numbers = append(numbers, n)
	}

	tmpl.Execute(w, numbers)
}

// ------------------------------
//
//	Add phone number handler
//
// ------------------------------
func addNumberHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	phone := formatPhoneNumber(r.FormValue("phone_number"))
	reason := strings.TrimSpace(r.FormValue("reason"))
	location := strings.TrimSpace(r.FormValue("store_location"))
	checkAmount := strings.TrimSpace(r.FormValue("check_amount"))
	notes := r.FormValue("notes")

	// Validate required fields
	if phone == "" || !validatePhoneNumber(phone) ||
		reason == "" || location == "" || checkAmount == "" {
		log.Printf("Validation failed - Phone: %s, Reason: %s, Location: %s, Amount: %s",
			phone, reason, location, checkAmount)
		http.Error(w, "Required fields are missing or invalid", http.StatusBadRequest)
		return
	}

	var checkAmountVal sql.NullFloat64
	amount, err := strconv.ParseFloat(checkAmount, 64)
	if err != nil || amount <= 0 {
		http.Error(w, "Invalid check amount", http.StatusBadRequest)
		return
	}
	checkAmountVal = sql.NullFloat64{Float64: amount, Valid: true}

	// Updated SQL statement
	stmt, err := db.Prepare(`
        INSERT INTO blocked_numbers 
        (phone_number, reason, store_location, incident_date, check_amount, notes, updated_at) 
        VALUES ($1, $2, $3, CURRENT_DATE, $4, $5, CURRENT_TIMESTAMP) 
        ON CONFLICT (phone_number) 
        DO UPDATE SET 
            reason = EXCLUDED.reason,
            store_location = EXCLUDED.store_location,
            incident_date = CURRENT_DATE,
            check_amount = EXCLUDED.check_amount,
            notes = CASE 
                WHEN blocked_numbers.notes IS NULL THEN EXCLUDED.notes
                WHEN EXCLUDED.notes IS NULL THEN blocked_numbers.notes
                ELSE EXCLUDED.notes || E'\n' || blocked_numbers.notes
            END,
            updated_at = CURRENT_TIMESTAMP`)
	if err != nil {
		log.Printf("Prepare error: %v", err)
		http.Error(w, "DB prepare error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(phone, reason, location, checkAmountVal, notes)
	if err != nil {
		log.Printf("Execute error: %v", err)
		http.Error(w, "DB insert error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ------------------------------
//
//	Search handler: find a specific blocked number
//
// ------------------------------
func searchHandler(w http.ResponseWriter, r *http.Request) {
	phone := r.URL.Query().Get("phone")
	log.Printf("Received search request for phone: %s", phone)

	if phone == "" {
		log.Printf("Empty phone number received")
		http.Error(w, "Phone number required", http.StatusBadRequest)
		return
	}

	searchPhone := formatPhoneNumber(phone)
	log.Printf("Search query - Original: %s, Formatted: %s", phone, searchPhone)

	var number BlockedNumber
	var checkAmount sql.NullFloat64
	query := `
        SELECT 
            id, 
            phone_number, 
            COALESCE(reason, '') as reason, 
            COALESCE(store_location, '') as store_location, 
            TO_CHAR(incident_date, 'YYYY-MM-DD') as incident_date, 
            TO_CHAR(created_at, 'YYYY-MM-DD HH24:MI:SS') as created_at, 
            check_amount, 
            COALESCE(notes, '') as notes
        FROM blocked_numbers 
        WHERE phone_number = $1`

	err := db.QueryRow(query, searchPhone).Scan(
		&number.ID,
		&number.PhoneNumber,
		&number.Reason,
		&number.StoreLocation,
		&number.IncidentDate,
		&number.CreatedAt,
		&checkAmount,
		&number.Notes)

	if checkAmount.Valid {
		number.CheckAmount = checkAmount.Float64
	} else {
		number.CheckAmount = 0
	}

	if err == sql.ErrNoRows {
		log.Printf("No match found for: %s", searchPhone)
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("Search error for %s: %v", searchPhone, err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	log.Printf("Found match: %+v", number)
	w.Header().Set("Content-Type", "application/json")

	response, err := json.Marshal(number)
	if err != nil {
		log.Printf("JSON marshal error: %v", err)
		http.Error(w, "Response encoding error", http.StatusInternalServerError)
		return
	}

	w.Write(response)
}

// ------------------------------
//
//	Initialize database schema
//
// ------------------------------
func initDatabase() error {
	createTableSQL := `
    CREATE TABLE IF NOT EXISTS blocked_numbers (
        id SERIAL PRIMARY KEY,
        phone_number VARCHAR(15) UNIQUE NOT NULL CHECK (phone_number ~ '^\+[1-9][0-9]{9,14}$'),
        reason VARCHAR(100) NOT NULL,
        store_location VARCHAR(100) NOT NULL,
        incident_date DATE NOT NULL,
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

// ------------------------------
//
//	Main function (entry point)
//
// ------------------------------
func main() {
	// Connect to PostgreSQL
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)

	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Error connecting to DB:", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Cannot reach database:", err)
	}

	fmt.Println("✅ Connected to database successfully!")

	// Initialize database schema
	if err := initDatabase(); err != nil {
		log.Fatal("Error initializing database:", err)
	}
	fmt.Println("✅ Database schema initialized!")

	// Route setup
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/add", addNumberHandler)
	http.HandleFunc("/search", searchHandler)

	// Serve static files (CSS/JS)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println("🚀 Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
