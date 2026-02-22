package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"regexp"
	"strings"

	"github.com/joledev/api-scheduler/models"
	"github.com/joledev/api-scheduler/services"
)

var slotDateRegex = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

type SlotHandler struct {
	db *sql.DB
}

func NewSlotHandler(db *sql.DB) *SlotHandler {
	return &SlotHandler{db: db}
}

// GetAvailableSlots returns computed available slots for a date range (public)
func (h *SlotHandler) GetAvailableSlots(w http.ResponseWriter, r *http.Request) {
	ip := r.RemoteAddr
	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		ip = strings.Split(fwd, ",")[0]
	}
	if !limiter.allow(strings.TrimSpace(ip), 60) {
		http.Error(w, `{"success":false,"message":"Too many requests"}`, http.StatusTooManyRequests)
		return
	}

	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	if !slotDateRegex.MatchString(from) || !slotDateRegex.MatchString(to) {
		http.Error(w, `{"success":false,"message":"from and to must be YYYY-MM-DD format"}`, http.StatusBadRequest)
		return
	}

	slots, err := services.GetAvailableSlots(h.db, from, to)
	if err != nil {
		http.Error(w, `{"success":false,"message":"Internal error"}`, http.StatusInternalServerError)
		return
	}

	if slots == nil {
		slots = []models.AvailableSlot{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.AvailableSlotsResponse{Slots: slots})
}
