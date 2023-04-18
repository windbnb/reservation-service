package handler

import (
	"github.com/windbnb/reservation-service/service"
)

type Handler struct {
	Service *service.ReservationService
}