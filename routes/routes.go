package routes

import (
	"net/http"

	"github.com/Mkhan2217/blocklist_app/handlers"
)

func RegisterRoutes(mux *http.ServeMux) {

	// UI Homepage
	mux.HandleFunc("/", handlers.HomeHandler)

	// Unified REST Endpoint for blocklist operations
	mux.HandleFunc("/api/blocklist", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {

		case http.MethodPost: // Block or Update
			handlers.CreateOrUpdateBlockedNumber(w, r)

		case http.MethodGet: // Search
			handlers.GetBlockedNumberHandler(w, r)

		case http.MethodDelete: // Unblock
			handlers.UnblockNumberHandler(w, r)

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
