package models

type Booking struct {
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
	ConfirmToken   string `json:"-"`
	RejectToken    string `json:"-"`
	CreatedAt      string `json:"createdAt,omitempty"`
}

type BookingRequest struct {
	Date           string `json:"date"`
	StartTime      string `json:"startTime"`
	MeetingType    string `json:"meetingType"`
	ClientName     string `json:"clientName"`
	ClientEmail    string `json:"clientEmail"`
	ClientPhone    string `json:"clientPhone"`
	ClientCompany  string `json:"clientCompany"`
	ClientAddress  string `json:"clientAddress"`
	ClientTimezone string `json:"clientTimezone"`
	Notes          string `json:"notes"`
	Lang           string `json:"lang"`
}

type BookingResponse struct {
	Success   bool   `json:"success"`
	BookingID string `json:"bookingId,omitempty"`
	Message   string `json:"message"`
}
