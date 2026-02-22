package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/joledev/api-scheduler/models"
	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test db: %v", err)
	}

	_, err = db.Exec(`CREATE TABLE bookings (
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
		t.Fatalf("Failed to create bookings table: %v", err)
	}

	return db
}

func insertBooking(t *testing.T, db *sql.DB, date, start, end, email, status string) {
	_, err := db.Exec(
		`INSERT INTO bookings (booking_id, date, start_time, end_time, meeting_type,
		 client_name, client_email, status, confirm_token, reject_token)
		 VALUES (?, ?, ?, ?, 'videollamada', 'Test', ?, ?, 'ct-test', 'rt-test')`,
		"BK-2026-TEST", date, start, end, email, status)
	if err != nil {
		t.Fatalf("Failed to insert booking: %v", err)
	}
}

func TestCreateBookingSuccess(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	handler := NewBookingHandler(db)
	// Use a far-future date to ensure it's a weekday and available
	body, _ := json.Marshal(models.BookingRequest{
		Date:        "2026-06-15", // Monday
		StartTime:   "09:00",
		MeetingType: "videollamada",
		ClientName:  "Test User",
		ClientEmail: "test@example.com",
		Lang:        "es",
	})

	req := httptest.NewRequest("POST", "/scheduler/bookings", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateBooking(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp models.BookingResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if !resp.Success {
		t.Errorf("Expected success=true, got false: %s", resp.Message)
	}
	if resp.BookingID == "" {
		t.Error("Expected a booking ID")
	}

	// Verify booking was created with pending status
	var status string
	db.QueryRow("SELECT status FROM bookings WHERE booking_id = ?", resp.BookingID).Scan(&status)
	if status != "pending" {
		t.Errorf("Expected status 'pending', got '%s'", status)
	}

	// Verify tokens were generated
	var confirmToken, rejectToken string
	db.QueryRow("SELECT confirm_token, reject_token FROM bookings WHERE booking_id = ?", resp.BookingID).Scan(&confirmToken, &rejectToken)
	if confirmToken == "" || rejectToken == "" {
		t.Error("Expected confirm and reject tokens to be generated")
	}
}

func TestCreateBookingDuplicateEmail(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Insert an existing active booking for this email
	insertBooking(t, db, "2026-06-15", "09:00", "09:30", "test@example.com", "pending")

	handler := NewBookingHandler(db)
	body, _ := json.Marshal(models.BookingRequest{
		Date:        "2026-06-16", // Tuesday
		StartTime:   "11:00",
		MeetingType: "videollamada",
		ClientName:  "Test User",
		ClientEmail: "test@example.com",
		Lang:        "es",
	})

	req := httptest.NewRequest("POST", "/scheduler/bookings", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateBooking(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("Expected 409 for duplicate email, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCreateBookingMissingEmail(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	handler := NewBookingHandler(db)
	body, _ := json.Marshal(models.BookingRequest{
		Date:        "2026-06-15",
		StartTime:   "09:00",
		MeetingType: "videollamada",
		ClientName:  "Test User",
		ClientEmail: "",
		Lang:        "es",
	})

	req := httptest.NewRequest("POST", "/scheduler/bookings", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateBooking(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCreateBookingBufferBlocks(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Insert a booking at 09:00
	insertBooking(t, db, "2026-06-15", "09:00", "09:30", "other@example.com", "confirmed")

	handler := NewBookingHandler(db)
	// Try to book at 10:30 — should be blocked (90 min < 120 min buffer)
	body, _ := json.Marshal(models.BookingRequest{
		Date:        "2026-06-15",
		StartTime:   "10:30",
		MeetingType: "videollamada",
		ClientName:  "Test User",
		ClientEmail: "new@example.com",
		Lang:        "es",
	})

	req := httptest.NewRequest("POST", "/scheduler/bookings", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateBooking(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("Expected 409 for buffer conflict, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCreateBookingBufferAllows(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Insert a booking at 09:00
	insertBooking(t, db, "2026-06-15", "09:00", "09:30", "other@example.com", "confirmed")

	handler := NewBookingHandler(db)
	// Try to book at 11:00 — should be allowed (120 min = exactly 2h, which is NOT < 120)
	body, _ := json.Marshal(models.BookingRequest{
		Date:        "2026-06-15",
		StartTime:   "11:00",
		MeetingType: "videollamada",
		ClientName:  "Test User",
		ClientEmail: "new@example.com",
		Lang:        "es",
	})

	req := httptest.NewRequest("POST", "/scheduler/bookings", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateBooking(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200 for allowed buffer, got %d: %s", w.Code, w.Body.String())
	}
}

func TestGetAvailableSlotsWeekdaysOnly(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	handler := NewSlotHandler(db)

	// 2026-06-13 = Saturday, 2026-06-14 = Sunday, 2026-06-15 = Monday
	req := httptest.NewRequest("GET", "/scheduler/slots?from=2026-06-13&to=2026-06-15", nil)
	w := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Get("/scheduler/slots", handler.GetAvailableSlots)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp models.AvailableSlotsResponse
	json.NewDecoder(w.Body).Decode(&resp)

	// Should only have slots for Monday (no Sat/Sun)
	for _, slot := range resp.Slots {
		if slot.Date == "2026-06-13" || slot.Date == "2026-06-14" {
			t.Errorf("Expected no slots on weekend, got slot on %s", slot.Date)
		}
	}

	// Monday should have 14 slots (09:00 to 15:30)
	mondaySlots := 0
	for _, slot := range resp.Slots {
		if slot.Date == "2026-06-15" {
			mondaySlots++
		}
	}
	if mondaySlots != 14 {
		t.Errorf("Expected 14 slots on Monday, got %d", mondaySlots)
	}
}

func TestGetAvailableSlotsWithBuffer(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Insert booking at 09:00 on Monday
	insertBooking(t, db, "2026-06-15", "09:00", "09:30", "someone@example.com", "confirmed")

	handler := NewSlotHandler(db)
	req := httptest.NewRequest("GET", "/scheduler/slots?from=2026-06-15&to=2026-06-15", nil)
	w := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Get("/scheduler/slots", handler.GetAvailableSlots)
	r.ServeHTTP(w, req)

	var resp models.AvailableSlotsResponse
	json.NewDecoder(w.Body).Decode(&resp)

	// With booking at 09:00 and 2h buffer, slots 09:00-10:30 should be blocked
	// Available should be 11:00 onwards = 11:00,11:30,12:00,12:30,13:00,13:30,14:00,14:30,15:00,15:30 = 10 slots
	// But slots at 09:30, 10:00, 10:30 are within 2h of 09:00 too
	// Blocked: 09:00 (diff=0), 09:30 (diff=30), 10:00 (diff=60), 10:30 (diff=90) - all < 120
	// Available from 11:00 (diff=120, NOT < 120)
	blockedTimes := map[string]bool{
		"09:00": true, "09:30": true, "10:00": true, "10:30": true,
	}

	for _, slot := range resp.Slots {
		if blockedTimes[slot.StartTime] {
			t.Errorf("Expected slot %s to be blocked by buffer, but it was available", slot.StartTime)
		}
	}

	// Check 11:00 IS available
	found1100 := false
	for _, slot := range resp.Slots {
		if slot.StartTime == "11:00" {
			found1100 = true
			break
		}
	}
	if !found1100 {
		t.Error("Expected 11:00 to be available (exactly 2h from 09:00)")
	}
}

func TestConfirmBooking(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Insert a pending booking with known tokens
	_, err := db.Exec(
		`INSERT INTO bookings (booking_id, date, start_time, end_time, meeting_type,
		 client_name, client_email, status, confirm_token, reject_token, lang)
		 VALUES ('BK-2026-001', '2026-06-15', '09:00', '09:30', 'videollamada',
		 'Test User', 'test@example.com', 'pending', 'confirm-token-123', 'reject-token-456', 'es')`)
	if err != nil {
		t.Fatal(err)
	}

	handler := NewBookingHandler(db)
	req := httptest.NewRequest("GET", "/scheduler/bookings/confirm?token=confirm-token-123", nil)
	w := httptest.NewRecorder()

	handler.ConfirmBooking(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}

	// Verify status changed to confirmed
	var status string
	db.QueryRow("SELECT status FROM bookings WHERE booking_id = 'BK-2026-001'").Scan(&status)
	if status != "confirmed" {
		t.Errorf("Expected status 'confirmed', got '%s'", status)
	}
}

func TestRejectBooking(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	_, err := db.Exec(
		`INSERT INTO bookings (booking_id, date, start_time, end_time, meeting_type,
		 client_name, client_email, status, confirm_token, reject_token, lang)
		 VALUES ('BK-2026-002', '2026-06-15', '11:00', '11:30', 'presencial',
		 'Test User', 'test@example.com', 'pending', 'confirm-token-789', 'reject-token-012', 'es')`)
	if err != nil {
		t.Fatal(err)
	}

	handler := NewBookingHandler(db)
	req := httptest.NewRequest("GET", "/scheduler/bookings/reject?token=reject-token-012", nil)
	w := httptest.NewRecorder()

	handler.RejectBooking(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}

	var status string
	db.QueryRow("SELECT status FROM bookings WHERE booking_id = 'BK-2026-002'").Scan(&status)
	if status != "rejected" {
		t.Errorf("Expected status 'rejected', got '%s'", status)
	}
}

func TestConfirmAlreadyConfirmed(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	_, err := db.Exec(
		`INSERT INTO bookings (booking_id, date, start_time, end_time, meeting_type,
		 client_name, client_email, status, confirm_token, reject_token, lang)
		 VALUES ('BK-2026-003', '2026-06-15', '13:00', '13:30', 'videollamada',
		 'Test User', 'test@example.com', 'confirmed', 'confirm-token-aaa', 'reject-token-bbb', 'es')`)
	if err != nil {
		t.Fatal(err)
	}

	handler := NewBookingHandler(db)
	req := httptest.NewRequest("GET", "/scheduler/bookings/confirm?token=confirm-token-aaa", nil)
	w := httptest.NewRecorder()

	handler.ConfirmBooking(w, req)

	// Should return 200 with info page, not error
	if w.Code != http.StatusOK {
		t.Errorf("Expected 200 (info page), got %d", w.Code)
	}

	// Status should still be confirmed
	var status string
	db.QueryRow("SELECT status FROM bookings WHERE booking_id = 'BK-2026-003'").Scan(&status)
	if status != "confirmed" {
		t.Errorf("Expected status to remain 'confirmed', got '%s'", status)
	}
}

func TestCreateBooking_InvalidDateFormat(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	handler := NewBookingHandler(db)
	body, _ := json.Marshal(models.BookingRequest{
		Date:        "15-06-2026", // wrong format
		StartTime:   "09:00",
		MeetingType: "videollamada",
		ClientName:  "Test User",
		ClientEmail: "test@example.com",
		Lang:        "es",
	})

	req := httptest.NewRequest("POST", "/scheduler/bookings", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateBooking(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 for invalid date format, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCreateBooking_InvalidTimeFormat(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	handler := NewBookingHandler(db)
	body, _ := json.Marshal(models.BookingRequest{
		Date:        "2026-06-15",
		StartTime:   "9:00", // missing leading zero
		MeetingType: "videollamada",
		ClientName:  "Test User",
		ClientEmail: "test@example.com",
		Lang:        "es",
	})

	req := httptest.NewRequest("POST", "/scheduler/bookings", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateBooking(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 for invalid time format, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCreateBooking_InvalidMeetingType(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	handler := NewBookingHandler(db)
	body, _ := json.Marshal(models.BookingRequest{
		Date:        "2026-06-15",
		StartTime:   "09:00",
		MeetingType: "phone", // invalid
		ClientName:  "Test User",
		ClientEmail: "test@example.com",
		Lang:        "es",
	})

	req := httptest.NewRequest("POST", "/scheduler/bookings", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateBooking(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 for invalid meeting type, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCreateBooking_NameTooLong(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	handler := NewBookingHandler(db)
	body, _ := json.Marshal(models.BookingRequest{
		Date:        "2026-06-15",
		StartTime:   "09:00",
		MeetingType: "videollamada",
		ClientName:  strings.Repeat("x", 201),
		ClientEmail: "test@example.com",
		Lang:        "es",
	})

	req := httptest.NewRequest("POST", "/scheduler/bookings", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Forwarded-For", "10.0.0.99")
	w := httptest.NewRecorder()

	handler.CreateBooking(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 for name too long, got %d: %s", w.Code, w.Body.String())
	}
}

func TestGetAvailableSlots_InvalidDateFormat(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	handler := NewSlotHandler(db)
	req := httptest.NewRequest("GET", "/scheduler/slots?from=invalid&to=2026-06-15", nil)
	w := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Get("/scheduler/slots", handler.GetAvailableSlots)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 for invalid date format, got %d: %s", w.Code, w.Body.String())
	}
}
