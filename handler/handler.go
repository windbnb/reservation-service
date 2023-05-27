package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/windbnb/reservation-service/client"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: err.Error(), StatusCode: http.StatusBadRequest})
		return
	}

	reservationRequest, err := h.Service.SaveReservationRequest(createReservationRequest)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: err.Error(), StatusCode: http.StatusBadRequest})
		return
	}

	var reservationRequestDto = model.ReservationRequestDto{
		ID:              reservationRequest.ID.Hex(),
		Status:          reservationRequest.Status,
		GuestNumber:     reservationRequest.GuestNumber,
		GuestID:         reservationRequest.GuestID,
		AccommodationID: reservationRequest.AccommodationID,
		StartDate:       reservationRequest.StartDate,
		EndDate:         reservationRequest.EndDate}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(reservationRequestDto)
}

func (h *Handler) GetGuestsActive(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	guestID, _ := strconv.Atoi(params["id"])

	activeReservations := h.Service.GetGuestActiveReservations(uint(guestID))

	reservationRequestsDto := []model.ReservationRequestDto{}

	for _, reservationRequest := range *activeReservations {
		reservationRequestsDto = append(reservationRequestsDto, model.ReservationRequestDto{
			ID:              reservationRequest.ID.Hex(),
			Status:          reservationRequest.Status,
			GuestNumber:     reservationRequest.GuestNumber,
			GuestID:         reservationRequest.GuestID,
			AccommodationID: reservationRequest.AccommodationID,
			StartDate:       reservationRequest.StartDate,
			EndDate:         reservationRequest.EndDate})
	}

	w.WriteHeader(http.StatusOK)
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
			ID:              reservationRequest.ID.Hex(),
			Status:          reservationRequest.Status,
			GuestNumber:     reservationRequest.GuestNumber,
			GuestID:         reservationRequest.GuestID,
			AccommodationID: reservationRequest.AccommodationID,
			StartDate:       reservationRequest.StartDate,
			EndDate:         reservationRequest.EndDate})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(reservationRequestsDto)
}

func (h *Handler) DeleteReservationRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id, _ := params["id"]

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: err.Error(), StatusCode: http.StatusBadRequest})
		return
	}

	err = h.Service.DeleteReservationRequest(objectId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: err.Error(), StatusCode: http.StatusBadRequest})
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) AcceptReservationRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id, _ := params["id"]

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: err.Error(), StatusCode: http.StatusBadRequest})
		return
	}

	tokenString := r.Header.Get("Authorization")
	userResponse, err := client.AuthorizeHost(tokenString)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: err.Error(), StatusCode: http.StatusUnauthorized})
		return
	}
	if userResponse.Role != "HOST" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: "user is not a host", StatusCode: http.StatusUnauthorized})
		return
	}

	reservation, err := h.Service.AcceptReservationRequest(objectId, userResponse.Id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: err.Error(), StatusCode: http.StatusBadRequest})
		return
	}

	reservationRequestsDto := model.ReservationRequestDto{
		ID:              reservation.ID.Hex(),
		Status:          reservation.Status,
		GuestNumber:     reservation.GuestNumber,
		GuestID:         reservation.GuestID,
		AccommodationID: reservation.AccommodationID,
		StartDate:       reservation.StartDate,
		EndDate:         reservation.EndDate}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(reservationRequestsDto)
}
