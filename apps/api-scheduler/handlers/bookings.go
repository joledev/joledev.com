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

	"github.com/go-chi/chi/v5"
	"github.com/joledev/api-scheduler/models"
	"github.com/joledev/api-scheduler/services"
)

var emailRegex = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
var dateRegex = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
var timeRegex = regexp.MustCompile(`^\d{2}:\d{2}$`)

// Rate limiter: max 10 requests per IP per hour
type rateLimiter struct {
	mu       sync.Mutex
	requests map[string][]time.Time
}

var limiter = &rateLimiter{requests: make(map[string][]time.Time)}

func (rl *rateLimiter) allow(ip string, maxReqs int) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-1 * time.Hour)

	var valid []time.Time
	for _, t := range rl.requests[ip] {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}
	rl.requests[ip] = valid

	if len(valid) >= maxReqs {
		return false
	}

	rl.requests[ip] = append(rl.requests[ip], now)
	return true
}

func getClientIP(r *http.Request) string {
	ip := r.RemoteAddr
	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		ip = strings.Split(fwd, ",")[0]
	}
	return strings.TrimSpace(ip)
}

type BookingHandler struct {
	db *sql.DB
}

func NewBookingHandler(db *sql.DB) *BookingHandler {
	return &BookingHandler{db: db}
}

// CreateBooking creates a new booking request (public). Status starts as "pending".
func (h *BookingHandler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	ip := getClientIP(r)
	if !limiter.allow(ip, 10) {
		http.Error(w, `{"success":false,"message":"Too many requests. Please try again later."}`, http.StatusTooManyRequests)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 64*1024) // 64KB max

	var req models.BookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"success":false,"message":"Invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Validate required fields
	clientName := strings.TrimSpace(req.ClientName)
	if clientName == "" || len(clientName) > 200 {
		http.Error(w, `{"success":false,"message":"Name is required (max 200 chars)"}`, http.StatusBadRequest)
		return
	}
	if !emailRegex.MatchString(strings.TrimSpace(req.ClientEmail)) || len(req.ClientEmail) > 254 {
		http.Error(w, `{"success":false,"message":"Valid email is required"}`, http.StatusBadRequest)
		return
	}
	if req.MeetingType != "presencial" && req.MeetingType != "videollamada" {
		http.Error(w, `{"success":false,"message":"meetingType must be 'presencial' or 'videollamada'"}`, http.StatusBadRequest)
		return
	}
	if !dateRegex.MatchString(req.Date) {
		http.Error(w, `{"success":false,"message":"date must be YYYY-MM-DD format"}`, http.StatusBadRequest)
		return
	}
	if !timeRegex.MatchString(req.StartTime) {
		http.Error(w, `{"success":false,"message":"startTime must be HH:MM format"}`, http.StatusBadRequest)
		return
	}
	if len(req.ClientPhone) > 30 || len(req.ClientCompany) > 200 || len(req.ClientAddress) > 500 || len(req.Notes) > 2000 {
		http.Error(w, `{"success":false,"message":"Field too long"}`, http.StatusBadRequest)
		return
	}
	if req.Lang != "es" && req.Lang != "en" {
		req.Lang = "es"
	}

	clientEmail := strings.TrimSpace(req.ClientEmail)

	// Compute endTime = startTime + 30min
	endTime := addMinutes(req.StartTime, 30)

	// Check one active booking per email
	todayStr := time.Now().Format("2006-01-02")
	var activeCount int
	err := h.db.QueryRow(
		`SELECT COUNT(*) FROM bookings WHERE client_email = ? AND status IN ('pending', 'confirmed') AND date >= ?`,
		clientEmail, todayStr).Scan(&activeCount)
	if err != nil {
		http.Error(w, `{"success":false,"message":"Internal error"}`, http.StatusInternalServerError)
		return
	}
	if activeCount > 0 {
		msg := "Ya tienes una solicitud de reunión activa. Espera a que sea procesada o cancelada antes de agendar otra."
		if req.Lang == "en" {
			msg = "You already have an active meeting request. Wait for it to be processed or cancelled before scheduling another."
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(models.BookingResponse{Success: false, Message: msg})
		return
	}

	// Generate tokens
	confirmToken, err := services.GenerateToken()
	if err != nil {
		http.Error(w, `{"success":false,"message":"Internal error"}`, http.StatusInternalServerError)
		return
	}
	rejectToken, err := services.GenerateToken()
	if err != nil {
		http.Error(w, `{"success":false,"message":"Internal error"}`, http.StatusInternalServerError)
		return
	}

	// BEGIN IMMEDIATE transaction for atomicity
	tx, err := h.db.Begin()
	if err != nil {
		http.Error(w, `{"success":false,"message":"Internal error"}`, http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Re-verify availability inside transaction
	available, err := services.IsSlotAvailable(tx, req.Date, req.StartTime)
	if err != nil {
		http.Error(w, `{"success":false,"message":"Internal error"}`, http.StatusInternalServerError)
		return
	}
	if !available {
		msg := "Este horario ya no está disponible. Por favor selecciona otro."
		if req.Lang == "en" {
			msg = "This time slot is no longer available. Please select another."
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(models.BookingResponse{Success: false, Message: msg})
		return
	}

	// Generate booking ID
	bookingID := h.generateBookingID(tx)

	// Insert booking
	_, err = tx.Exec(
		`INSERT INTO bookings (booking_id, date, start_time, end_time, meeting_type,
		 client_name, client_email, client_phone, client_company, client_address,
		 client_timezone, notes, lang, status, confirm_token, reject_token)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 'pending', ?, ?)`,
		bookingID, req.Date, req.StartTime, endTime, req.MeetingType,
		strings.TrimSpace(req.ClientName), clientEmail,
		req.ClientPhone, req.ClientCompany, req.ClientAddress,
		req.ClientTimezone, req.Notes, req.Lang,
		confirmToken, rejectToken)
	if err != nil {
		http.Error(w, `{"success":false,"message":"Internal error"}`, http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, `{"success":false,"message":"Internal error"}`, http.StatusInternalServerError)
		return
	}

	// Build booking struct for emails
	booking := &models.Booking{
		BookingID:      bookingID,
		Date:           req.Date,
		StartTime:      req.StartTime,
		EndTime:        endTime,
		MeetingType:    req.MeetingType,
		ClientName:     strings.TrimSpace(req.ClientName),
		ClientEmail:    clientEmail,
		ClientPhone:    req.ClientPhone,
		ClientCompany:  req.ClientCompany,
		ClientAddress:  req.ClientAddress,
		ClientTimezone: req.ClientTimezone,
		Notes:          req.Notes,
		Lang:           req.Lang,
		Status:         "pending",
		ConfirmToken:   confirmToken,
		RejectToken:    rejectToken,
	}

	// Send emails asynchronously
	go func() {
		if err := services.SendAdminPendingNotification(booking); err != nil {
			log.Printf("Error sending admin pending notification: %v", err)
		}
		if err := services.SendClientPendingNotification(booking); err != nil {
			log.Printf("Error sending client pending notification: %v", err)
		}
	}()

	msg := "Tu solicitud de reunión ha sido recibida. Te notificaremos cuando sea confirmada."
	if req.Lang == "en" {
		msg = "Your meeting request has been received. We'll notify you when it's confirmed."
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.BookingResponse{
		Success:   true,
		BookingID: bookingID,
		Message:   msg,
	})
}

// GetBooking returns booking details by public ID
func (h *BookingHandler) GetBooking(w http.ResponseWriter, r *http.Request) {
	ip := getClientIP(r)
	if !limiter.allow(ip, 10) {
		http.Error(w, `{"success":false,"message":"Too many requests"}`, http.StatusTooManyRequests)
		return
	}

	bookingID := chi.URLParam(r, "bookingId")
	if bookingID == "" {
		http.Error(w, `{"success":false,"message":"Booking ID required"}`, http.StatusBadRequest)
		return
	}

	var b models.Booking
	err := h.db.QueryRow(
		`SELECT id, booking_id, date, start_time, end_time, meeting_type,
		        client_name, client_email, client_phone, client_company, client_address,
		        COALESCE(client_timezone, ''), notes, lang, status, created_at
		 FROM bookings WHERE booking_id = ?`, bookingID).Scan(
		&b.ID, &b.BookingID, &b.Date, &b.StartTime, &b.EndTime, &b.MeetingType,
		&b.ClientName, &b.ClientEmail, &b.ClientPhone, &b.ClientCompany, &b.ClientAddress,
		&b.ClientTimezone, &b.Notes, &b.Lang, &b.Status, &b.CreatedAt)
	if err == sql.ErrNoRows {
		http.Error(w, `{"success":false,"message":"Booking not found"}`, http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, `{"success":false,"message":"Internal error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(b)
}

// ConfirmBooking confirms a booking via token (admin clicks link in email)
func (h *BookingHandler) ConfirmBooking(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		h.renderTokenPage(w, "error", "Token is required", "")
		return
	}

	var b models.Booking
	err := h.db.QueryRow(
		`SELECT id, booking_id, date, start_time, end_time, meeting_type,
		        client_name, client_email, COALESCE(client_phone, ''), COALESCE(client_company, ''),
		        COALESCE(client_address, ''), COALESCE(client_timezone, ''), COALESCE(notes, ''),
		        COALESCE(lang, 'es'), status
		 FROM bookings WHERE confirm_token = ?`, token).Scan(
		&b.ID, &b.BookingID, &b.Date, &b.StartTime, &b.EndTime, &b.MeetingType,
		&b.ClientName, &b.ClientEmail, &b.ClientPhone, &b.ClientCompany, &b.ClientAddress,
		&b.ClientTimezone, &b.Notes, &b.Lang, &b.Status)
	if err == sql.ErrNoRows {
		h.renderTokenPage(w, "error", "Invalid or expired token", "")
		return
	}
	if err != nil {
		h.renderTokenPage(w, "error", "Internal error", "")
		return
	}

	if b.Status != "pending" {
		h.renderTokenPage(w, "info", fmt.Sprintf("This booking is already %s", b.Status), b.BookingID)
		return
	}

	_, err = h.db.Exec(`UPDATE bookings SET status = 'confirmed' WHERE id = ?`, b.ID)
	if err != nil {
		h.renderTokenPage(w, "error", "Failed to confirm booking", "")
		return
	}

	// Send confirmation email to client
	go func() {
		if err := services.SendBookingConfirmation(&b); err != nil {
			log.Printf("Error sending confirmation email: %v", err)
		}
	}()

	h.renderTokenPage(w, "confirmed", fmt.Sprintf("Booking %s confirmed!", b.BookingID),
		fmt.Sprintf("%s — %s %s", b.ClientName, b.Date, b.StartTime))
}

// RejectBooking rejects a booking via token (admin clicks link in email)
func (h *BookingHandler) RejectBooking(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		h.renderTokenPage(w, "error", "Token is required", "")
		return
	}

	var b models.Booking
	err := h.db.QueryRow(
		`SELECT id, booking_id, date, start_time, end_time, meeting_type,
		        client_name, client_email, COALESCE(client_phone, ''), COALESCE(client_company, ''),
		        COALESCE(client_address, ''), COALESCE(client_timezone, ''), COALESCE(notes, ''),
		        COALESCE(lang, 'es'), status
		 FROM bookings WHERE reject_token = ?`, token).Scan(
		&b.ID, &b.BookingID, &b.Date, &b.StartTime, &b.EndTime, &b.MeetingType,
		&b.ClientName, &b.ClientEmail, &b.ClientPhone, &b.ClientCompany, &b.ClientAddress,
		&b.ClientTimezone, &b.Notes, &b.Lang, &b.Status)
	if err == sql.ErrNoRows {
		h.renderTokenPage(w, "error", "Invalid or expired token", "")
		return
	}
	if err != nil {
		h.renderTokenPage(w, "error", "Internal error", "")
		return
	}

	if b.Status != "pending" {
		h.renderTokenPage(w, "info", fmt.Sprintf("This booking is already %s", b.Status), b.BookingID)
		return
	}

	_, err = h.db.Exec(`UPDATE bookings SET status = 'rejected' WHERE id = ?`, b.ID)
	if err != nil {
		h.renderTokenPage(w, "error", "Failed to reject booking", "")
		return
	}

	// Send rejection email to client
	go func() {
		if err := services.SendBookingRejection(&b); err != nil {
			log.Printf("Error sending rejection email: %v", err)
		}
	}()

	h.renderTokenPage(w, "rejected", fmt.Sprintf("Booking %s rejected.", b.BookingID),
		fmt.Sprintf("%s — %s %s", b.ClientName, b.Date, b.StartTime))
}

// GetAdminBookings returns all bookings in a date range (admin)
func (h *BookingHandler) GetAdminBookings(w http.ResponseWriter, r *http.Request) {
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	if !dateRegex.MatchString(from) || !dateRegex.MatchString(to) {
		http.Error(w, `{"success":false,"message":"from and to must be YYYY-MM-DD format"}`, http.StatusBadRequest)
		return
	}

	rows, err := h.db.Query(
		`SELECT id, booking_id, date, start_time, end_time, meeting_type,
		        client_name, client_email, client_phone, client_company, client_address,
		        COALESCE(client_timezone, ''), notes, lang, status, created_at
		 FROM bookings WHERE date >= ? AND date <= ?
		 ORDER BY date, start_time`, from, to)
	if err != nil {
		http.Error(w, `{"success":false,"message":"Internal error"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	bookings := []models.AdminBooking{}
	for rows.Next() {
		var b models.AdminBooking
		if err := rows.Scan(
			&b.ID, &b.BookingID, &b.Date, &b.StartTime, &b.EndTime, &b.MeetingType,
			&b.ClientName, &b.ClientEmail, &b.ClientPhone, &b.ClientCompany, &b.ClientAddress,
			&b.ClientTimezone, &b.Notes, &b.Lang, &b.Status, &b.CreatedAt,
		); err != nil {
			continue
		}
		bookings = append(bookings, b)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.AdminBookingsResponse{Bookings: bookings})
}

// CancelBooking cancels a booking (admin)
func (h *BookingHandler) CancelBooking(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	r.Body = http.MaxBytesReader(w, r.Body, 4*1024) // 4KB max

	var status struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&status); err != nil || status.Status != "cancelled" {
		http.Error(w, `{"success":false,"message":"Status must be 'cancelled'"}`, http.StatusBadRequest)
		return
	}

	var b models.Booking
	err := h.db.QueryRow(
		`SELECT id, booking_id, date, start_time, end_time, meeting_type,
		        client_name, client_email, COALESCE(client_phone, ''), COALESCE(client_company, ''),
		        COALESCE(client_address, ''), COALESCE(client_timezone, ''), COALESCE(notes, ''),
		        COALESCE(lang, 'es'), status
		 FROM bookings WHERE id = ?`, idStr).Scan(
		&b.ID, &b.BookingID, &b.Date, &b.StartTime, &b.EndTime, &b.MeetingType,
		&b.ClientName, &b.ClientEmail, &b.ClientPhone, &b.ClientCompany, &b.ClientAddress,
		&b.ClientTimezone, &b.Notes, &b.Lang, &b.Status)
	if err == sql.ErrNoRows {
		http.Error(w, `{"success":false,"message":"Booking not found"}`, http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, `{"success":false,"message":"Internal error"}`, http.StatusInternalServerError)
		return
	}

	if b.Status == "cancelled" {
		http.Error(w, `{"success":false,"message":"Booking already cancelled"}`, http.StatusConflict)
		return
	}

	_, err = h.db.Exec(`UPDATE bookings SET status = 'cancelled' WHERE id = ?`, idStr)
	if err != nil {
		http.Error(w, `{"success":false,"message":"Internal error"}`, http.StatusInternalServerError)
		return
	}

	// Send cancellation email
	go func() {
		if err := services.SendBookingCancellation(&b); err != nil {
			log.Printf("Error sending cancellation email: %v", err)
		}
	}()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Booking cancelled",
	})
}

func (h *BookingHandler) generateBookingID(tx *sql.Tx) string {
	year := time.Now().Year()
	var count int
	tx.QueryRow("SELECT COUNT(*) FROM bookings WHERE booking_id LIKE ?", fmt.Sprintf("BK-%d-%%", year)).Scan(&count)
	return fmt.Sprintf("BK-%d-%03d", year, count+1)
}

func addMinutes(timeStr string, mins int) string {
	if len(timeStr) < 5 {
		return timeStr
	}
	h := int(timeStr[0]-'0')*10 + int(timeStr[1]-'0')
	m := int(timeStr[3]-'0')*10 + int(timeStr[4]-'0')
	total := h*60 + m + mins
	return fmt.Sprintf("%02d:%02d", total/60, total%60)
}

// renderTokenPage renders a simple HTML page for confirm/reject token responses
func (h *BookingHandler) renderTokenPage(w http.ResponseWriter, status, message, detail string) {
	var bgColor, icon string
	switch status {
	case "confirmed":
		bgColor = "#22c55e"
		icon = "&#10004;"
	case "rejected":
		bgColor = "#ef4444"
		icon = "&#10006;"
	case "info":
		bgColor = "#3b82f6"
		icon = "&#8505;"
	default:
		bgColor = "#ef4444"
		icon = "&#9888;"
	}

	detailHTML := ""
	if detail != "" {
		detailHTML = fmt.Sprintf(`<p style="color:#6b7280;margin-top:0.5rem;font-size:0.875rem">%s</p>`, detail)
	}

	html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1">
<title>JoleDev Scheduler</title></head>
<body style="margin:0;min-height:100vh;display:flex;align-items:center;justify-content:center;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,sans-serif;background:#f9fafb">
<div style="text-align:center;padding:2rem;max-width:400px">
<div style="width:64px;height:64px;border-radius:50%%;background:%s;color:#fff;display:inline-flex;align-items:center;justify-content:center;font-size:2rem;margin-bottom:1rem">%s</div>
<h1 style="font-size:1.25rem;margin:0 0 0.5rem">%s</h1>
%s
</div>
</body>
</html>`, bgColor, icon, message, detailHTML)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}
