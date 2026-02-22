package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/joledev/api-quoter/models"
	"github.com/joledev/api-quoter/services"
)

var emailRegex = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

// Rate limiter: max 5 requests per IP per hour
type rateLimiter struct {
	mu       sync.Mutex
	requests map[string][]time.Time
}

var limiter = &rateLimiter{requests: make(map[string][]time.Time)}

func (rl *rateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-1 * time.Hour)

	// Clean old entries
	var valid []time.Time
	for _, t := range rl.requests[ip] {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}
	rl.requests[ip] = valid

	if len(valid) >= 5 {
		return false
	}

	rl.requests[ip] = append(rl.requests[ip], now)
	return true
}

type QuoteHandler struct {
	db *sql.DB
}

func NewQuoteHandler(db *sql.DB) *QuoteHandler {
	return &QuoteHandler{db: db}
}

func (h *QuoteHandler) CreateQuote(w http.ResponseWriter, r *http.Request) {
	// Rate limit
	ip := r.RemoteAddr
	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		ip = strings.Split(fwd, ",")[0]
	}
	if !limiter.allow(strings.TrimSpace(ip)) {
		http.Error(w, `{"success":false,"message":"Too many requests. Please try again later."}`, http.StatusTooManyRequests)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 64*1024) // 64KB max

	var req models.QuoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"success":false,"message":"Invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Validate required fields
	name := strings.TrimSpace(req.Contact.Name)
	email := strings.TrimSpace(req.Contact.Email)
	if name == "" || len(name) > 200 {
		http.Error(w, `{"success":false,"message":"Name is required (max 200 chars)"}`, http.StatusBadRequest)
		return
	}
	if !emailRegex.MatchString(email) || len(email) > 254 {
		http.Error(w, `{"success":false,"message":"Valid email is required"}`, http.StatusBadRequest)
		return
	}
	if len(req.ProjectTypes) == 0 || len(req.ProjectTypes) > 20 {
		http.Error(w, `{"success":false,"message":"At least one project type is required"}`, http.StatusBadRequest)
		return
	}
	if len(req.Contact.Phone) > 30 || len(req.Contact.Company) > 200 || len(req.Contact.Notes) > 2000 {
		http.Error(w, `{"success":false,"message":"Field too long"}`, http.StatusBadRequest)
		return
	}

	// Generate quote ID
	quoteID := h.generateQuoteID()

	// Save to DB
	projectTypesJSON, _ := json.Marshal(req.ProjectTypes)
	featuresJSON, _ := json.Marshal(req.Features)

	includeSourceCodeInt := 0
	if req.IncludeSourceCode {
		includeSourceCodeInt = 1
	}

	_, err := h.db.Exec(`INSERT INTO quotes (quote_id, project_types, features, business_size, current_state, timeline, currency, estimated_min, estimated_max, payment_plan, include_source_code, contact_name, contact_email, contact_phone, contact_company, contact_notes, lang) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		quoteID, string(projectTypesJSON), string(featuresJSON),
		req.BusinessSize, req.CurrentState, req.Timeline, req.Currency,
		req.EstimatedMin, req.EstimatedMax,
		req.PaymentPlan, includeSourceCodeInt,
		strings.TrimSpace(req.Contact.Name), strings.TrimSpace(req.Contact.Email),
		req.Contact.Phone, req.Contact.Company, req.Contact.Notes, req.Lang)
	if err != nil {
		log.Printf("Error saving quote: %v", err)
		http.Error(w, `{"success":false,"message":"Internal error"}`, http.StatusInternalServerError)
		return
	}

	// Send emails (non-blocking, log errors)
	go func() {
		if err := services.SendQuoteNotification(&req, quoteID); err != nil {
			log.Printf("Error sending notification email: %v", err)
		}
		if err := services.SendQuoteConfirmation(&req, quoteID); err != nil {
			log.Printf("Error sending confirmation email: %v", err)
		}
	}()

	msg := "Cotización enviada correctamente"
	if req.Lang == "en" {
		msg = "Quote sent successfully"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.QuoteResponse{
		Success: true,
		Message: msg,
		QuoteID: quoteID,
	})
}

func (h *QuoteHandler) generateQuoteID() string {
	year := time.Now().Year()
	var count int
	h.db.QueryRow("SELECT COUNT(*) FROM quotes WHERE quote_id LIKE ?", fmt.Sprintf("QT-%d-%%", year)).Scan(&count)
	return fmt.Sprintf("QT-%d-%03d", year, count+1)
}
