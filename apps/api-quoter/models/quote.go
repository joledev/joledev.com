package models

import "time"

type QuoteContact struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Company string `json:"company"`
	Notes   string `json:"notes"`
}

type QuoteRequest struct {
	ProjectTypes      []string     `json:"projectTypes"`
	Features          []string     `json:"features"`
	BusinessSize      string       `json:"businessSize"`
	CurrentState      string       `json:"currentState"`
	Timeline          string       `json:"timeline"`
	Currency          string       `json:"currency"`
	EstimatedMin      int          `json:"estimatedMin"`
	EstimatedMax      int          `json:"estimatedMax"`
	PaymentPlan       string       `json:"paymentPlan"`
	IncludeSourceCode bool         `json:"includeSourceCode"`
	Contact           QuoteContact `json:"contact"`
	Lang              string       `json:"lang"`
}

type QuoteResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	QuoteID string `json:"quoteId"`
}

type Quote struct {
	ID                int       `json:"id"`
	QuoteID           string    `json:"quoteId"`
	ProjectTypes      string    `json:"projectTypes"`
	Features          string    `json:"features"`
	BusinessSize      string    `json:"businessSize"`
	CurrentState      string    `json:"currentState"`
	Timeline          string    `json:"timeline"`
	Currency          string    `json:"currency"`
	EstimatedMin      int       `json:"estimatedMin"`
	EstimatedMax      int       `json:"estimatedMax"`
	PaymentPlan       string    `json:"paymentPlan"`
	IncludeSourceCode bool      `json:"includeSourceCode"`
	ContactName       string    `json:"contactName"`
	ContactEmail      string    `json:"contactEmail"`
	ContactPhone      string    `json:"contactPhone"`
	ContactCo         string    `json:"contactCompany"`
	ContactNotes      string    `json:"contactNotes"`
	Lang              string    `json:"lang"`
	CreatedAt         time.Time `json:"createdAt"`
}
