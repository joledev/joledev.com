package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/joledev/api-quoter/models"
	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test db: %v", err)
	}

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
		payment_plan TEXT DEFAULT '',
		include_source_code INTEGER DEFAULT 0,
		contact_name TEXT NOT NULL,
		contact_email TEXT NOT NULL,
		contact_phone TEXT,
		contact_company TEXT,
		contact_notes TEXT,
		lang TEXT DEFAULT 'es',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	return db
}

func TestCreateQuote_ValidRequest(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	handler := NewQuoteHandler(db)

	req := models.QuoteRequest{
		ProjectTypes: []string{"web"},
		Features:     []string{"auth"},
		BusinessSize: "small",
		CurrentState: "fromScratch",
		Timeline:     "1-3months",
		Currency:     "MXN",
		EstimatedMin: 25000,
		EstimatedMax: 40000,
		Contact: models.QuoteContact{
			Name:  "Test User",
			Email: "test@example.com",
		},
		Lang: "es",
	}

	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/quotes", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateQuote(w, httpReq)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp models.QuoteResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if !resp.Success {
		t.Errorf("Expected success=true, got false: %s", resp.Message)
	}
	if resp.QuoteID == "" {
		t.Error("Expected non-empty quoteId")
	}
}

func TestCreateQuote_MissingEmail(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	handler := NewQuoteHandler(db)

	req := models.QuoteRequest{
		ProjectTypes: []string{"web"},
		Features:     []string{"auth"},
		BusinessSize: "small",
		CurrentState: "fromScratch",
		Timeline:     "1-3months",
		Currency:     "MXN",
		EstimatedMin: 25000,
		EstimatedMax: 40000,
		Contact: models.QuoteContact{
			Name:  "Test User",
			Email: "",
		},
		Lang: "es",
	}

	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/quotes", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateQuote(w, httpReq)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCreateQuote_MissingName(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	handler := NewQuoteHandler(db)
	req := models.QuoteRequest{
		ProjectTypes: []string{"web"},
		Features:     []string{"auth"},
		BusinessSize: "small",
		CurrentState: "fromScratch",
		Timeline:     "1-3months",
		Currency:     "MXN",
		EstimatedMin: 25000,
		EstimatedMax: 40000,
		Contact: models.QuoteContact{
			Name:  "",
			Email: "test@example.com",
		},
	}

	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/quotes", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateQuote(w, httpReq)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 for missing name, got %d", w.Code)
	}
}

func TestCreateQuote_EmptyProjectTypes(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	handler := NewQuoteHandler(db)
	req := models.QuoteRequest{
		ProjectTypes: []string{},
		Features:     []string{"auth"},
		BusinessSize: "small",
		CurrentState: "fromScratch",
		Timeline:     "1-3months",
		Currency:     "MXN",
		EstimatedMin: 25000,
		EstimatedMax: 40000,
		Contact: models.QuoteContact{
			Name:  "Test User",
			Email: "test@example.com",
		},
	}

	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/quotes", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateQuote(w, httpReq)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 for empty project types, got %d", w.Code)
	}
}

func TestCreateQuote_NotesTooLong(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	handler := NewQuoteHandler(db)
	req := models.QuoteRequest{
		ProjectTypes: []string{"web"},
		Features:     []string{"auth"},
		BusinessSize: "small",
		CurrentState: "fromScratch",
		Timeline:     "1-3months",
		Currency:     "MXN",
		EstimatedMin: 25000,
		EstimatedMax: 40000,
		Contact: models.QuoteContact{
			Name:  "Test User",
			Email: "test@example.com",
			Notes: strings.Repeat("x", 2001),
		},
	}

	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/quotes", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateQuote(w, httpReq)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 for notes too long, got %d", w.Code)
	}
}

func TestCreateQuote_NameTooLong(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	handler := NewQuoteHandler(db)
	req := models.QuoteRequest{
		ProjectTypes: []string{"web"},
		Features:     []string{"auth"},
		BusinessSize: "small",
		CurrentState: "fromScratch",
		Timeline:     "1-3months",
		Currency:     "MXN",
		EstimatedMin: 25000,
		EstimatedMax: 40000,
		Contact: models.QuoteContact{
			Name:  strings.Repeat("x", 201),
			Email: "test@example.com",
		},
	}

	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/quotes", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-Forwarded-For", "10.0.0.99")
	w := httptest.NewRecorder()

	handler.CreateQuote(w, httpReq)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 for name too long, got %d", w.Code)
	}
}
