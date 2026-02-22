package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/joledev/api-scheduler/models"
	"github.com/joledev/api-scheduler/services"
)

type SlotHandler struct {
	db *sql.DB
}

func NewSlotHandler(db *sql.DB) *SlotHandler {
	return &SlotHandler{db: db}
}

// GetAvailableSlots returns computed available slots for a date range (public)
func (h *SlotHandler) GetAvailableSlots(w http.ResponseWriter, r *http.Request) {
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	if from == "" || to == "" {
		http.Error(w, `{"success":false,"message":"from and to query params required"}`, http.StatusBadRequest)
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
