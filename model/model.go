package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type ReservationRequestStatus string

const (
	SUBMITTED ReservationRequestStatus = "SUBMITTED"
	ACCEPTED  ReservationRequestStatus = "ACCEPTED"
	DECLINED  ReservationRequestStatus = "DECLINED"
)

type ReservationRequest struct {
	ID              primitive.ObjectID       `bson:"_id"`
	StartDate       time.Time                `bson:"startDate"`
	EndDate         time.Time                `bson:"endDate"`
	AccommodationID uint                     `bson:"accommodationID"`
	GuestID         uint                     `bson:"guestID"`
	GuestNumber     uint                     `bson:"guestNumber"`
	Status          ReservationRequestStatus `bson:"status"`
	OwnerID         uint                     `bson:"ownerID"`
}
