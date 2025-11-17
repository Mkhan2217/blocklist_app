package handlers

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"github.com/Mkhan2217/blocklist_app/models"
	"github.com/Mkhan2217/blocklist_app/utils"
)

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

/* ----------------------- HOME PAGE RENDER ----------------------- */

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	numbers, err := models.GetAllBlockedNumbers()
	if err != nil {
		http.Error(w, "DB query error", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, numbers)
}

/* ------------------ POST: Block or Update Number ------------------ */

func CreateOrUpdateBlockedNumber(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var number models.BlockedNumber
	if err := json.NewDecoder(r.Body).Decode(&number); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	// Validation
	if number.PhoneNumber == "" ||
		!utils.ValidatePhoneNumber(number.PhoneNumber) ||
		number.Reason == "" ||
		number.StoreLocation == "" ||
		number.CheckAmount <= 0 { // float64 works now
		writeError(w, http.StatusBadRequest, "Required fields missing or invalid")
		return
	}

	number.PhoneNumber = utils.FormatPhoneNumber(number.PhoneNumber)

	if err := models.UpsertBlockedNumber(number); err != nil {
		log.Println("DB Upsert error:", err)
		writeError(w, http.StatusInternalServerError, "Database error")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Number blocked"})
}

/* ---------------------- GET: Search Number ------------------------ */

func GetBlockedNumberHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	phone := utils.FormatPhoneNumber(r.URL.Query().Get("phone"))
	if phone == "" {
		writeError(w, http.StatusBadRequest, "Phone number required")
		return
	}

	record, err := models.GetBlockedNumber(phone)
	if err == sql.ErrNoRows {
		writeError(w, http.StatusNotFound, "Number not found")
		return
	}
	if err != nil {
		log.Println("DB query error:", err)
		writeError(w, http.StatusInternalServerError, "Database error")
		return
	}

	writeJSON(w, http.StatusOK, record.ToResponse())
}

/* --------------------- DELETE: Unblock Number --------------------- */

func UnblockNumberHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	phone := utils.FormatPhoneNumber(r.URL.Query().Get("phone"))
	if phone == "" {
		writeError(w, http.StatusBadRequest, "Missing phone number")
		return
	}

	if err := models.DeleteBlockedNumber(phone); err != nil {
		log.Println("DB delete error:", err)
		writeError(w, http.StatusInternalServerError, "Database error")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Number unblocked"})
}
