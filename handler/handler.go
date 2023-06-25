package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"github.com/windbnb/reservation-service/client"
	"github.com/windbnb/reservation-service/tracer"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/windbnb/reservation-service/model"
	"github.com/windbnb/reservation-service/service"
)

type Handler struct {
	Service *service.ReservationRequestService
	Tracer  opentracing.Tracer
	Closer  io.Closer
}

func (handler *Handler) Healthcheck(w http.ResponseWriter, _ *http.Request) {
    _, _ = fmt.Fprintln(w, "Healthy!")
}

func (handler *Handler) Ready(w http.ResponseWriter, _ *http.Request) {
    _, _ = fmt.Fprintln(w, "Ready!")
}

func (h *Handler) CreateReservationRequest(w http.ResponseWriter, r *http.Request) {
	span := tracer.StartSpanFromRequest("createReservationRequestHandler", h.Tracer, r)
	defer span.Finish()
	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling create reservation request at %s\n", r.URL.Path)),
	)
	ctx := tracer.ContextWithSpan(context.Background(), span)

	userResponse := h.authorizeGuest(r)
	if userResponse == nil || userResponse.Role != "GUEST" {
		tracer.LogError(span, errors.New("Unauthorized"))
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: "user is not a guest", StatusCode: http.StatusUnauthorized})
		return
	}
	w.Header().Set("Content-Type", "application/json")

	var createReservationRequest *model.CreateReservationRequest
	err := json.NewDecoder(r.Body).Decode(&createReservationRequest)
	if err != nil {
		tracer.LogError(span, err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: err.Error(), StatusCode: http.StatusBadRequest})
		return
	}

	createReservationRequest.GuestID = userResponse.Id
	reservationRequest, err := h.Service.SaveReservationRequest(createReservationRequest, ctx)

	if err != nil {
		tracer.LogError(span, err)
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
	span := tracer.StartSpanFromRequest("getGuestsActiveHandler", h.Tracer, r)
	defer span.Finish()
	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling get guests active reservations at %s\n", r.URL.Path)),
	)
	ctx := tracer.ContextWithSpan(context.Background(), span)

	userResponse := h.authorizeGuest(r)
	if userResponse.Role != "GUEST" {
		tracer.LogError(span, errors.New("Unauthorized"))
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: "user is not a guest", StatusCode: http.StatusUnauthorized})
		return
	}

	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	guestID, _ := strconv.Atoi(params["id"])

	activeReservations := h.Service.GetGuestActiveReservations(uint(guestID), ctx)

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
	span := tracer.StartSpanFromRequest("getOwnersActiveHandler", h.Tracer, r)
	defer span.Finish()
	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling get owners active reservations at %s\n", r.URL.Path)),
	)
	ctx := tracer.ContextWithSpan(context.Background(), span)

	userResponse := h.authorizeHost(r)
	if userResponse.Role != "HOST" {
		tracer.LogError(span, errors.New("Unauthorized"))
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: "user is not a host", StatusCode: http.StatusUnauthorized})
		return
	}

	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	guestID, _ := strconv.Atoi(params["id"])

	activeReservations := h.Service.GetOwnersActiveReservations(uint(guestID), ctx)

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
	span := tracer.StartSpanFromRequest("deleteReservationRequestHandler", h.Tracer, r)
	defer span.Finish()
	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling delete reservation request at %s\n", r.URL.Path)),
	)
	ctx := tracer.ContextWithSpan(context.Background(), span)

	userResponse := h.authorizeGuest(r)
	if userResponse.Role != "GUEST" {
		tracer.LogError(span, errors.New("Unauthorized"))
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: "user is not a guest", StatusCode: http.StatusUnauthorized})
		return
	}

	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id, _ := params["id"]

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		tracer.LogError(span, err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: err.Error(), StatusCode: http.StatusBadRequest})
		return
	}

	err = h.Service.DeleteReservationRequest(objectId, userResponse.Id, ctx)
	if err != nil {
		tracer.LogError(span, err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: err.Error(), StatusCode: http.StatusBadRequest})
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) AcceptReservationRequest(w http.ResponseWriter, r *http.Request) {
	span := tracer.StartSpanFromRequest("acceptReservationRequestHandler", h.Tracer, r)
	defer span.Finish()
	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling accept reservation request at %s\n", r.URL.Path)),
	)
	ctx := tracer.ContextWithSpan(context.Background(), span)

	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id, _ := params["id"]

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		tracer.LogError(span, err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: err.Error(), StatusCode: http.StatusBadRequest})
		return
	}

	userResponse := h.authorizeHost(r)
	if userResponse.Role != "HOST" {
		tracer.LogError(span, errors.New("Unauthorized"))
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: "user is not a host", StatusCode: http.StatusUnauthorized})
		return
	}

	reservation, err := h.Service.AcceptReservationRequest(objectId, userResponse.Id, ctx)

	if err != nil {
		tracer.LogError(span, err)
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

func (h *Handler) CancelReservationRequest(w http.ResponseWriter, r *http.Request) {
	span := tracer.StartSpanFromRequest("cancelReservationRequestHandler", h.Tracer, r)
	defer span.Finish()
	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling cancel reservation request at %s\n", r.URL.Path)),
	)
	ctx := tracer.ContextWithSpan(context.Background(), span)

	userResponse := h.authorizeGuest(r)
	if userResponse.Role != "GUEST" {
		tracer.LogError(span, errors.New("Unauthorized"))
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: "user is not a guest", StatusCode: http.StatusUnauthorized})
		return
	}

	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id, _ := params["id"]

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		tracer.LogError(span, err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: err.Error(), StatusCode: http.StatusBadRequest})
		return
	}

	_, err = h.Service.CancelReservationRequest(objectId, userResponse.Id, ctx)
	if err != nil {
		tracer.LogError(span, err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: err.Error(), StatusCode: http.StatusBadRequest})
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) CountGuestsCancelledReservations(w http.ResponseWriter, r *http.Request) {
	span := tracer.StartSpanFromRequest("getCountGuestCancelledReservationsHandler", h.Tracer, r)
	defer span.Finish()
	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling count guests cancelled reservations at %s\n", r.URL.Path)),
	)
	ctx := tracer.ContextWithSpan(context.Background(), span)

	userResponse := h.authorizeHost(r)
	if userResponse.Role != "HOST" {
		tracer.LogError(span, errors.New("Unauthorized"))
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: "user is not a host", StatusCode: http.StatusUnauthorized})
		return
	}

	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	guestId, err := strconv.Atoi(params["guestId"])
	if err != nil {
		tracer.LogError(span, err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: err.Error(), StatusCode: http.StatusBadRequest})
		return
	}

	count := h.Service.CountCancelledReservations(uint(guestId), ctx)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model.CancelledReservations{
		Count: count})
}

func (h *Handler) authorizeHost(r *http.Request) *model.UserResponseDTO {
	tokenString := r.Header.Get("Authorization")
	userResponse, err := client.AuthorizeHost(tokenString)
	if err != nil {
		return nil
	}

	return &userResponse
}

func (h *Handler) authorizeGuest(r *http.Request) *model.UserResponseDTO {
	tokenString := r.Header.Get("Authorization")
	userResponse, err := client.AuthorizeGueest(tokenString)
	if err != nil {
		return nil
	}

	return &userResponse
}
