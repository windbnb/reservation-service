package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Reservation struct {
}

type ReservationRequestStatus string

const (
	SUBMITTED ReservationRequestStatus = "SUBMITTED"
	ACCEPTED  ReservationRequestStatus = "ACCEPTED"
	DECLINED  ReservationRequestStatus = "DECLINED"
)

type ReservationRequest struct {
	gorm.Model
	StartDate       time.Time
	EndDate         time.Time
	AccommodationID uint
	GuestID         uint
	GuestNumber     uint
	Status          ReservationRequestStatus
	OwnerID         uint
}
