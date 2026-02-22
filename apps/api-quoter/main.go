package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joledev/api-quoter/handlers"
	_ "github.com/mattn/go-sqlite3"
)

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		next.ServeHTTP(w, r)
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := os.Getenv("CORS_ORIGIN")
		if origin == "" {
			origin = "https://joledev.com"
		}

		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	// Database setup
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./data/quotes.db"
	}

	// Ensure data directory exists
	os.MkdirAll("./data", 0755)

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Create table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS quotes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		quote_id TEXT UNIQUE NOT NULL,
		project_types TEXT NOT NULL,
		features TEXT NOT NULL,
		business_size TEXT NOT NULL,
		current_state TEXT NOT NULL,
		timeline TEXT NOT NULL,
		currency TEXT NOT NULL,
		estimated_min INTEGER NOT NULL,
		estimated_max INTEGER NOT NULL,
		contact_name TEXT NOT NULL,
		contact_email TEXT NOT NULL,
		contact_phone TEXT,
		contact_company TEXT,
		contact_notes TEXT,
		lang TEXT DEFAULT 'es',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	// Migrations: add new columns (ignore errors if columns already exist)
	db.Exec(`ALTER TABLE quotes ADD COLUMN payment_plan TEXT DEFAULT ''`)
	db.Exec(`ALTER TABLE quotes ADD COLUMN include_source_code INTEGER DEFAULT 0`)

	// Router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(securityHeaders)
	r.Use(corsMiddleware)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	quoteHandler := handlers.NewQuoteHandler(db)
	r.Post("/quotes", quoteHandler.CreateQuote)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("api-quoter listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
