package model

import "time"

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

type ReservationRequestDto struct {
	Status               ReservationRequestStatus `json:"status"`
	GuestID              uint                     `json:"guestID"`
	AccommodationID      uint                     `json:"accommodationID"`
	ReservationRequestID uint                     `json:"reservationRequestID"`
	StartDate            time.Time                `json:"startDate"`
	EndDate              time.Time                `json:"endDate"`
	GuestNumber          uint                     `json:"guestNumber"`
}
