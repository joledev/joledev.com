package models

// AvailableSlot represents a computed available time slot (no DB table)
type AvailableSlot struct {
	Date      string `json:"date"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

type AvailableSlotsResponse struct {
	Slots []AvailableSlot `json:"slots"`
}

// AdminBooking is a booking with full details for the admin panel
type AdminBooking struct {
	ID             int    `json:"id"`
	BookingID      string `json:"bookingId"`
	Date           string `json:"date"`
	StartTime      string `json:"startTime"`
	EndTime        string `json:"endTime"`
	MeetingType    string `json:"meetingType"`
	ClientName     string `json:"clientName"`
	ClientEmail    string `json:"clientEmail"`
	ClientPhone    string `json:"clientPhone,omitempty"`
	ClientCompany  string `json:"clientCompany,omitempty"`
	ClientAddress  string `json:"clientAddress,omitempty"`
	ClientTimezone string `json:"clientTimezone,omitempty"`
	Notes          string `json:"notes,omitempty"`
	Lang           string `json:"lang"`
	Status         string `json:"status"`
	CreatedAt      string `json:"createdAt,omitempty"`
}

type AdminBookingsResponse struct {
	Bookings []AdminBooking `json:"bookings"`
}
