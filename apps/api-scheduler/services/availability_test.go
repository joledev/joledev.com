package services

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestTimeToMinutes(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"00:00", 0},
		{"09:00", 540},
		{"09:30", 570},
		{"15:30", 930},
		{"16:00", 960},
		{"23:59", 1439},
		{"bad", -1},
		{"", -1},
	}

	for _, tt := range tests {
		got := timeToMinutes(tt.input)
		if got != tt.expected {
			t.Errorf("timeToMinutes(%q) = %d, want %d", tt.input, got, tt.expected)
		}
	}
}

func TestMinutesToTime(t *testing.T) {
	tests := []struct {
		input    int
		expected string
	}{
		{0, "00:00"},
		{540, "09:00"},
		{570, "09:30"},
		{930, "15:30"},
		{960, "16:00"},
	}

	for _, tt := range tests {
		got := minutesToTime(tt.input)
		if got != tt.expected {
			t.Errorf("minutesToTime(%d) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()
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
		t.Fatalf("Failed to create table: %v", err)
	}

	return db
}

func TestGetAvailableSlots_WeekdaysOnly(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// 2027-01-09 = Saturday? No, let me pick a known date range
	// 2026-06-13 = Saturday, 2026-06-14 = Sunday, 2026-06-15 = Monday
	slots, err := GetAvailableSlots(db, "2026-06-13", "2026-06-14")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(slots) != 0 {
		t.Errorf("Expected 0 slots on weekend, got %d", len(slots))
	}
}

func TestGetAvailableSlots_MondaySlotCount(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// 2026-06-15 = Monday, far future so no "past" filtering
	slots, err := GetAvailableSlots(db, "2026-06-15", "2026-06-15")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Slots from 09:00 to 15:30 every 30min = 14 slots
	if len(slots) != 14 {
		t.Errorf("Expected 14 slots on Monday, got %d", len(slots))
	}

	// First slot should be 09:00, last should be 15:30
	if len(slots) > 0 {
		if slots[0].StartTime != "09:00" {
			t.Errorf("Expected first slot at 09:00, got %s", slots[0].StartTime)
		}
		if slots[len(slots)-1].StartTime != "15:30" {
			t.Errorf("Expected last slot at 15:30, got %s", slots[len(slots)-1].StartTime)
		}
	}
}

func TestGetAvailableSlots_BufferBlocking(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Insert a confirmed booking at 12:00 on Monday 2026-06-15
	_, err := db.Exec(
		`INSERT INTO bookings (booking_id, date, start_time, end_time, meeting_type,
		 client_name, client_email, status, confirm_token, reject_token)
		 VALUES ('BK-TEST', '2026-06-15', '12:00', '12:30', 'videollamada',
		 'Test', 'test@test.com', 'confirmed', 'ct', 'rt')`)
	if err != nil {
		t.Fatal(err)
	}

	slots, err := GetAvailableSlots(db, "2026-06-15", "2026-06-15")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Blocked slots: anything within 2h of 12:00
	// 10:00 (diff=120, NOT < 120, so AVAILABLE)
	// 10:30 (diff=90, < 120, BLOCKED)
	// 11:00 (diff=60, BLOCKED)
	// 11:30 (diff=30, BLOCKED)
	// 12:00 (diff=0, BLOCKED)
	// 12:30 (diff=30, BLOCKED)
	// 13:00 (diff=60, BLOCKED)
	// 13:30 (diff=90, BLOCKED)
	// 14:00 (diff=120, NOT < 120, so AVAILABLE)

	blockedTimes := map[string]bool{
		"10:30": true, "11:00": true, "11:30": true,
		"12:00": true, "12:30": true, "13:00": true, "13:30": true,
	}

	for _, slot := range slots {
		if blockedTimes[slot.StartTime] {
			t.Errorf("Expected %s to be blocked by 2h buffer, but it was available", slot.StartTime)
		}
	}

	// 10:00 and 14:00 should be available (exactly at boundary)
	slotMap := make(map[string]bool)
	for _, slot := range slots {
		slotMap[slot.StartTime] = true
	}
	if !slotMap["10:00"] {
		t.Error("Expected 10:00 to be available (exactly 2h before 12:00)")
	}
	if !slotMap["14:00"] {
		t.Error("Expected 14:00 to be available (exactly 2h after 12:00)")
	}
}
