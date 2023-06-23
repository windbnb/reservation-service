package model

import "time"

type ErrorResponse struct {
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}

type AcceptReservationType string

const (
	MANUAL        AcceptReservationType = "MANUAL"
	AUTOMATICALLY AcceptReservationType = "AUTOMATICALLY"
)

type AccommodationInfo struct {
	Id                    uint                  `json:"id"`
	MinimimGuests         uint                  `json:"minimimGuests"`
	MaximumGuests         uint                  `json:"maximumGuests"`
	AvailableTerms        []AvailableTerm       `json:"availableTerms"`
	UserID                uint                  `json:"userID"`
	AcceptReservationType AcceptReservationType `json:"acceptReservationType"`
	Name                  string                `json:"name"`
}

type AvailableTerm struct {
	StartDate time.Time
	EndDate   time.Time
}

type ReservationRequestDto struct {
	Status            ReservationRequestStatus `json:"status"`
	GuestID           uint                     `json:"guestID"`
	AccommodationID   uint                     `json:"accommodationID"`
	StartDate         time.Time                `json:"startDate"`
	EndDate           time.Time                `json:"endDate"`
	GuestNumber       uint                     `json:"guestNumber"`
	ID                string                   `json:"id"`
	AccommodationName string                   `json:"accommodationName"`
}

type UserRole string

const (
	HOST  UserRole = "HOST"
	GUEST UserRole = "GUEST"
)

type UserResponseDTO struct {
	Id       uint     `json:"id"`
	Email    string   `json:"email"`
	Name     string   `json:"name"`
	Surname  string   `json:"surname"`
	Address  string   `json:"address"`
	Username string   `json:"username"`
	Role     UserRole `json:"role"`
}

type ReservedTermRequest struct {
	StartDate      time.Time `json:"startDate"`
	EndDate        time.Time `json:"endDate"`
	AccomodationID uint      `json:"accomodationId"`
}

type ReservedTermResponse struct {
	StartDate      time.Time `json:"startDate"`
	EndDate        time.Time `json:"endDate"`
	AccomodationID uint      `json:"accomodationId"`
	Id             uint      `json:"id"`
}

type CancelledReservations struct {
	Count int `json:"count"`
}
