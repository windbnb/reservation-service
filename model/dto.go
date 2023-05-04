package model

import "time"

type ReservationRequestDTO struct {
}

type ErrorResponse struct {
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}

type AccommodationInfo struct {
	Id             uint            `json:"id"`
	MinimimGuests  uint            `json:"minimimGuests"`
	MaximumGuests  uint            `json:"maximumGuests"`
	AvailableTerms []AvailableTerm `json:"availableTerms"`
	UserID         uint            `json:"userID"`
}

type AvailableTerm struct {
	StartDate time.Time
	EndDate   time.Time
}
