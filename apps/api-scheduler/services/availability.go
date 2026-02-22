package services

import (
	"database/sql"
	"fmt"
	"math"
	"time"

	"github.com/joledev/api-scheduler/models"
)

var tijuanaTZ *time.Location

func init() {
	var err error
	tijuanaTZ, err = time.LoadLocation("America/Tijuana")
	if err != nil {
		// Fallback: America/Tijuana is UTC-8 (PST) / UTC-7 (PDT), same as LA
		tijuanaTZ, _ = time.LoadLocation("America/Los_Angeles")
	}
}

// GetAvailableSlots computes available 30-min slots for the given date range.
// Business hours: Mon-Fri 9:00-16:00 America/Tijuana, every 30 min.
// Slots blocked by existing bookings (pending or confirmed) with a 2h buffer.
func GetAvailableSlots(db *sql.DB, fromDate, toDate string) ([]models.AvailableSlot, error) {
	from, err := time.Parse("2006-01-02", fromDate)
	if err != nil {
		return nil, err
	}
	to, err := time.Parse("2006-01-02", toDate)
	if err != nil {
		return nil, err
	}

	// Get all bookings (pending + confirmed) in the date range
	rows, err := db.Query(
		`SELECT date, start_time FROM bookings
		 WHERE status IN ('pending', 'confirmed')
		 AND date >= ? AND date <= ?`,
		fromDate, toDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Map: date -> list of booked start times in minutes since midnight
	bookedByDate := make(map[string][]int)
	for rows.Next() {
		var date, startTime string
		if err := rows.Scan(&date, &startTime); err != nil {
			continue
		}
		mins := timeToMinutes(startTime)
		if mins >= 0 {
			bookedByDate[date] = append(bookedByDate[date], mins)
		}
	}

	now := time.Now().In(tijuanaTZ)
	todayStr := now.Format("2006-01-02")

	var result []models.AvailableSlot

	for d := from; !d.After(to); d = d.AddDate(0, 0, 1) {
		weekday := d.Weekday()
		if weekday == time.Saturday || weekday == time.Sunday {
			continue
		}

		dateStr := d.Format("2006-01-02")
		booked := bookedByDate[dateStr]

		// Generate slots from 9:00 to 15:30 (last slot starts at 15:30, ends at 16:00)
		for mins := 540; mins <= 930; mins += 30 { // 540 = 9*60, 930 = 15*60+30
			// If today, skip slots whose start time <= current time
			if dateStr == todayStr {
				currentMins := now.Hour()*60 + now.Minute()
				if mins <= currentMins {
					continue
				}
			}

			// Check 2h buffer against all booked slots
			blocked := false
			for _, bMins := range booked {
				if math.Abs(float64(mins-bMins)) < 120 {
					blocked = true
					break
				}
			}
			if blocked {
				continue
			}

			endMins := mins + 30
			result = append(result, models.AvailableSlot{
				Date:      dateStr,
				StartTime: minutesToTime(mins),
				EndTime:   minutesToTime(endMins),
			})
		}
	}

	return result, nil
}

// IsSlotAvailable checks if a specific slot on a specific date is available.
// Used inside transactions to re-verify before inserting.
func IsSlotAvailable(tx *sql.Tx, date, startTime string) (bool, error) {
	slotMins := timeToMinutes(startTime)
	if slotMins < 0 {
		return false, nil
	}

	// Check it's a valid business hours slot
	if slotMins < 540 || slotMins > 930 || slotMins%30 != 0 {
		return false, nil
	}

	// Check weekday
	d, err := time.Parse("2006-01-02", date)
	if err != nil {
		return false, err
	}
	if d.Weekday() == time.Saturday || d.Weekday() == time.Sunday {
		return false, nil
	}

	// Check not in the past
	now := time.Now().In(tijuanaTZ)
	todayStr := now.Format("2006-01-02")
	if date == todayStr {
		currentMins := now.Hour()*60 + now.Minute()
		if slotMins <= currentMins {
			return false, nil
		}
	} else if date < todayStr {
		return false, nil
	}

	// Check 2h buffer against existing bookings
	rows, err := tx.Query(
		`SELECT start_time FROM bookings
		 WHERE status IN ('pending', 'confirmed') AND date = ?`, date)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	for rows.Next() {
		var bookedStart string
		if err := rows.Scan(&bookedStart); err != nil {
			continue
		}
		bMins := timeToMinutes(bookedStart)
		if math.Abs(float64(slotMins-bMins)) < 120 {
			return false, nil
		}
	}

	return true, nil
}

func timeToMinutes(t string) int {
	if len(t) < 5 {
		return -1
	}
	h := int(t[0]-'0')*10 + int(t[1]-'0')
	m := int(t[3]-'0')*10 + int(t[4]-'0')
	return h*60 + m
}

func minutesToTime(mins int) string {
	h := mins / 60
	m := mins % 60
	return fmt.Sprintf("%02d:%02d", h, m)
}
