package handler

import (
	"encoding/json"
	"net/http"

	"github.com/windbnb/reservation-service/model"
	"github.com/windbnb/reservation-service/service"
)

type Handler struct {
	Service *service.ReservationRequestService
}

func (h *Handler) CreateReservationRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var createReservationRequest *model.CreateReservationRequest
	err := json.NewDecoder(r.Body).Decode(&createReservationRequest)
	if err != nil {
		return
	}

	reservationRequest, err := h.Service.SaveReservationRequest(createReservationRequest)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: err.Error(), StatusCode: http.StatusBadRequest})
		return
	}

	json.NewEncoder(w).Encode(reservationRequest)
}
