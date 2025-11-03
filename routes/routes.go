package routes

import (
	"net/http"

	"github.com/Mkhan2217/blocklist_app/handlers"
)

// Register all routes in one place
func RegisterRoutes() {
	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/add", handlers.AddNumberHandler)
	http.HandleFunc("/search", handlers.SearchHandler)
	http.HandleFunc("/unblock", handlers.UnblockNumberHandler)
}
