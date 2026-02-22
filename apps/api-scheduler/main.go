package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/joledev/api-scheduler/handlers"
	"github.com/joledev/api-scheduler/middleware"
	_ "github.com/mattn/go-sqlite3"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := os.Getenv("CORS_ORIGIN")
		if origin == "" {
			origin = "https://joledev.com"
		}

		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

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
		dbPath = "/data/scheduler.db"
	}

	os.MkdirAll("/data", 0755)

	db, err := sql.Open("sqlite3", dbPath+"?_journal_mode=WAL")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Migration: drop old tables if they exist
	db.Exec(`DROP TABLE IF EXISTS availability_slots`)

	// Drop old bookings table if it has the old schema (slot_id column)
	var hasSlotID bool
	row := db.QueryRow(`SELECT COUNT(*) FROM pragma_table_info('bookings') WHERE name='slot_id'`)
	row.Scan(&hasSlotID)
	if hasSlotID {
		db.Exec(`DROP TABLE IF EXISTS bookings`)
	}

	// Create new bookings table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS bookings (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		booking_id TEXT UNIQUE NOT NULL,
		date TEXT NOT NULL,
		start_time TEXT NOT NULL,
		end_time TEXT NOT NULL,
		meeting_type TEXT NOT NULL,
		client_name TEXT NOT NULL,
		client_email TEXT NOT NULL,
		client_phone TEXT,
		client_company TEXT,
		client_address TEXT,
		client_timezone TEXT,
		notes TEXT,
		lang TEXT DEFAULT 'es',
		status TEXT DEFAULT 'pending',
		confirm_token TEXT UNIQUE,
		reject_token TEXT UNIQUE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		log.Fatalf("Failed to create bookings table: %v", err)
	}

	// Indexes
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_bookings_status ON bookings(status)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_bookings_date ON bookings(date)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_bookings_email_status ON bookings(client_email, status)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_bookings_confirm_token ON bookings(confirm_token)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_bookings_reject_token ON bookings(reject_token)`)

	// Handlers
	slotHandler := handlers.NewSlotHandler(db)
	bookingHandler := handlers.NewBookingHandler(db)

	// Router
	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(corsMiddleware)

	// Health check
	r.Get("/scheduler/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	// Public routes
	r.Get("/scheduler/slots", slotHandler.GetAvailableSlots)
	r.Post("/scheduler/bookings", bookingHandler.CreateBooking)
	r.Get("/scheduler/bookings/{bookingId}", bookingHandler.GetBooking)

	// Token-based confirm/reject (public, no auth — links sent in admin email)
	r.Get("/scheduler/bookings/confirm", bookingHandler.ConfirmBooking)
	r.Get("/scheduler/bookings/reject", bookingHandler.RejectBooking)

	// Admin routes (Basic Auth protected)
	r.Route("/scheduler/admin", func(r chi.Router) {
		r.Use(middleware.AdminAuth)
		r.Get("/bookings", bookingHandler.GetAdminBookings)
		r.Patch("/bookings/{id}", bookingHandler.CancelBooking)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	log.Printf("api-scheduler listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
