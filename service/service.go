package service

import (
	"github.com/windbnb/reservation-service/repository"
)

type ReservationService struct {
	Repo *repository.Repository
}