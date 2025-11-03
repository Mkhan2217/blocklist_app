package handlers

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

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

/* ----------------------- ADD BLOCKED NUMBER ----------------------- */

func AddNumberHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	number := models.BlockedNumber{
		PhoneNumber:   utils.FormatPhoneNumber(r.FormValue("phone_number")),
		Reason:        strings.TrimSpace(r.FormValue("reason")),
		StoreLocation: strings.TrimSpace(r.FormValue("store_location")),
		Notes:         r.FormValue("notes"),
	}

	checkAmount := strings.TrimSpace(r.FormValue("check_amount"))

	if number.PhoneNumber == "" || !utils.ValidatePhoneNumber(number.PhoneNumber) ||
		number.Reason == "" || number.StoreLocation == "" || checkAmount == "" {
		writeError(w, http.StatusBadRequest, "Required fields missing or invalid")
		return
	}

	amount, err := strconv.ParseFloat(checkAmount, 64)
	if err != nil || amount <= 0 {
		writeError(w, http.StatusBadRequest, "Invalid check amount")
		return
	}

	number.CheckAmount = sql.NullFloat64{Float64: amount, Valid: true}

	if err := models.UpsertBlockedNumber(number); err != nil {
		log.Println("DB insert error:", err)
		writeError(w, http.StatusInternalServerError, "Database error")
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

/* ----------------------- SEARCH BLOCKED NUMBER API ----------------------- */

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	phone := utils.FormatPhoneNumber(r.URL.Query().Get("phone"))
	if phone == "" {
		writeError(w, http.StatusBadRequest, "Phone number required")
		return
	}

	record, err := models.GetBlockedNumber(phone)
	if err == sql.ErrNoRows {
		writeError(w, http.StatusNotFound, "Number not found")
		return
	} else if err != nil {
		writeError(w, http.StatusInternalServerError, "Database error")
		return
	}

	writeJSON(w, http.StatusOK, record.ToResponse())
}

/* ----------------------- UNBLOCK API ----------------------- */

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
		writeError(w, http.StatusInternalServerError, "Database error")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Unblocked"})
}
