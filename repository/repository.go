package repository

import (
	"github.com/jinzhu/gorm"

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
