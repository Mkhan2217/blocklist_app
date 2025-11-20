package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Mkhan2217/blocklist_app/internal/api/routes"
	"github.com/Mkhan2217/blocklist_app/internal/db"
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

	// Create HTTP request multiplexer
	mux := http.NewServeMux()

	// Register all API routes
	routes.RegisterRoutes(mux)
	// Serve static files through mux
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Create HTTP server
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Channel to listen for OS interrupt signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// Start server in a separate goroutine
	go func() {
		log.Println("🚀 Server running on http://localhost:8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-stop
	log.Println("⚠️  Shutting down server...")

	// Create context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("❌ Server forced to shutdown: %v", err)
	}

	log.Println("✅ Server stopped gracefully")
}
