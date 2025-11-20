package models

import (
	"fmt"

	"github.com/Mkhan2217/blocklist_app/internal/db"
)

/* --------------------------  DB Model (internal)  -------------------------- */

type BlockedNumber struct {
	ID            int
	PhoneNumber   string
	Reason        string
	StoreLocation string
	IncidentDate  string
	CreatedAt     string
	CheckAmount   float64
	Notes         string
}

/* --------------------------  API Response DTO  --------------------------- */

type BlockedNumberResponse struct {
	ID            int     `json:"id"`
	PhoneNumber   string  `json:"phoneNumber"`
	Reason        string  `json:"reason"`
	StoreLocation string  `json:"storeLocation"`
	IncidentDate  string  `json:"incidentDate"`
	CreatedAt     string  `json:"createdAt"`
	CheckAmount   float64 `json:"checkAmount"`
	Notes         string  `json:"notes"`
}

/* ------------- Convert DB Model → API Response DTO  -------- */
func (b BlockedNumber) ToResponse() BlockedNumberResponse {
return BlockedNumberResponse(b)
}

/* --------------------------  DB Queries  ---------------------------------- */

// Fetch all blocked numbers for home page
func GetAllBlockedNumbers() ([]BlockedNumber, error) {
	rows, err := db.DB.Query(
		"SELECT id, phone_number, created_at FROM blocked_numbers ORDER BY id DESC",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var numbers []BlockedNumber
	for rows.Next() {
		var n BlockedNumber
		if err := rows.Scan(&n.ID, &n.PhoneNumber, &n.CreatedAt); err != nil {
			return nil, err
		}
		numbers = append(numbers, n)
	}
	return numbers, nil
}

// GetBlockedNumber fetches one record
func GetBlockedNumber(phone string) (BlockedNumber, error) {
	var number BlockedNumber
	query := `SELECT id, phone_number, COALESCE(reason,''), COALESCE(store_location,''), 
			  TO_CHAR(incident_date,'YYYY-MM-DD'), TO_CHAR(created_at,'YYYY-MM-DD HH24:MI:SS'), 
			  check_amount, COALESCE(notes,'') 
			  FROM blocked_numbers WHERE phone_number=$1`

	err := db.DB.QueryRow(query, phone).Scan(
		&number.ID,
		&number.PhoneNumber,
		&number.Reason,
		&number.StoreLocation,
		&number.IncidentDate,
		&number.CreatedAt,
		&number.CheckAmount,
		&number.Notes,
	)

	return number, err
}

// Insert/update record
func UpsertBlockedNumber(n BlockedNumber) error {
	stmt := `
	INSERT INTO blocked_numbers
	(phone_number, reason, store_location, incident_date, check_amount, notes, updated_at)
	VALUES ($1,$2,$3,CURRENT_DATE,$4,$5,CURRENT_TIMESTAMP)
	ON CONFLICT (phone_number)
	DO UPDATE SET
		reason=EXCLUDED.reason,
		store_location=EXCLUDED.store_location,
		incident_date=CURRENT_DATE,
		check_amount=EXCLUDED.check_amount,
		notes=CASE
			WHEN blocked_numbers.notes IS NULL THEN EXCLUDED.notes
			WHEN EXCLUDED.notes IS NULL THEN blocked_numbers.notes
			ELSE EXCLUDED.notes || E'\n' || blocked_numbers.notes
		END,
		updated_at=CURRENT_TIMESTAMP
	`

	_, err := db.DB.Exec(stmt, n.PhoneNumber, n.Reason, n.StoreLocation, n.CheckAmount, n.Notes)
	return err
}

// Delete record
func DeleteBlockedNumber(phone string) error {
	result, err := db.DB.Exec("DELETE FROM blocked_numbers WHERE phone_number=$1", phone)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("number not found")
	}
	return nil
}
