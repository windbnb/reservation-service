package repository

import (
	"github.com/jinzhu/gorm"
	"time"

	"github.com/windbnb/reservation-service/model"
)

type Repository struct {
	Db *gorm.DB
}

func (r *Repository) FindAcceptedReservationRequests(accomodationId uint) *[]model.ReservationRequest {
	reservationRequests := &[]model.ReservationRequest{}

	r.Db.Where("accommodation_id = ? and status = 'ACCEPTED'", accomodationId).Find(reservationRequests)

	return reservationRequests
}

func (r *Repository) SaveReservationRequest(reservationRequest *model.ReservationRequest) *model.ReservationRequest {
	r.Db.Create(&reservationRequest)
	return reservationRequest
}

func (r *Repository) FindGuestsActive(guestID uint) *[]model.ReservationRequest {
	reservationRequests := &[]model.ReservationRequest{}

	r.Db.Where("guest_id = ? and status = 'ACCEPTED' and end_date >= ?", guestID, time.Now()).Find(reservationRequests)

	return reservationRequests
}

func (r *Repository) FindOwnersActive(ownerID uint) *[]model.ReservationRequest {
	reservationRequests := &[]model.ReservationRequest{}

	r.Db.Where("owner_id = ? and status = 'ACCEPTED' and end_date >= ?", ownerID, time.Now()).Find(reservationRequests)

	return reservationRequests
}
