package model

import (
	"time"
)

type CreateReservationRequest struct {
	StartDate       time.Time
	NumberOfDays    uint
	AccommodationID uint
	GuestID         uint
	GuestNumber     uint
}
