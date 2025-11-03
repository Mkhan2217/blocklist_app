package main

import (
	"fmt"
	"github.com/Mkhan2217/blocklist_app/db"
	"github.com/Mkhan2217/blocklist_app/routes"
	"log"
	"net/http"
)

func main() {
	// Connect to DB
	if err := db.ConnectDB(); err != nil {
		log.Fatal("DB connection error:", err)
	}
	defer db.CloseDB()

	if err := db.InitSchema(); err != nil {
		log.Fatal("DB schema initialization error:", err)
	}
	fmt.Println("✅ Database initialized")

	// Register routes
	routes.RegisterRoutes()

	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println("🚀 Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
