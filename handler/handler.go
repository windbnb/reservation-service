package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"

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

func (h *Handler) GetGuestsActive(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	guestID, _ := strconv.Atoi(params["id"])

	activeReservations := h.Service.GetGuestActiveReservations(uint(guestID))

	reservationRequestsDto := []model.ReservationRequestDto{}

	for _, reservationRequest := range *activeReservations {
		reservationRequestsDto = append(reservationRequestsDto, model.ReservationRequestDto{
			Status:          reservationRequest.Status,
			GuestNumber:     reservationRequest.GuestNumber,
			GuestID:         reservationRequest.GuestID,
			AccommodationID: reservationRequest.AccommodationID,
			StartDate:       reservationRequest.StartDate,
			EndDate:         reservationRequest.EndDate})
	}

	json.NewEncoder(w).Encode(reservationRequestsDto)
}

func (h *Handler) GetOwnersActive(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	guestID, _ := strconv.Atoi(params["id"])

	activeReservations := h.Service.GetOwnersActiveReservations(uint(guestID))

	reservationRequestsDto := []model.ReservationRequestDto{}

	for _, reservationRequest := range *activeReservations {
		reservationRequestsDto = append(reservationRequestsDto, model.ReservationRequestDto{
			Status:          reservationRequest.Status,
			GuestNumber:     reservationRequest.GuestNumber,
			GuestID:         reservationRequest.GuestID,
			AccommodationID: reservationRequest.AccommodationID,
			StartDate:       reservationRequest.StartDate,
			EndDate:         reservationRequest.EndDate})
	}

	json.NewEncoder(w).Encode(reservationRequestsDto)
}
