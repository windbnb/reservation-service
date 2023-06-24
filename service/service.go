package service

import (
	"context"
	"errors"
	"github.com/windbnb/reservation-service/client"
	"github.com/windbnb/reservation-service/model"
	"github.com/windbnb/reservation-service/repository"
	"github.com/windbnb/reservation-service/tracer"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
	"time"
)

var accommodationServiceUrl = os.Getenv("accommodationServiceUrl") + "/api/accomodation/"

type ReservationRequestService struct {
	Repo *repository.Repository
}

func (s *ReservationRequestService) SaveReservationRequest(createReservationRequest *model.CreateReservationRequest, ctx context.Context) (*model.ReservationRequest, error) {
	span := tracer.StartSpanFromContext(ctx, "saveReservationRequestService")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	if createReservationRequest.StartDate.Before(time.Now()) {
		return nil, errors.New("Start date cannot be in past")
	}

	if createReservationRequest.NumberOfDays <= 0 {
		return nil, errors.New("Number of days must be positive")
	}

	accommodationInfo, err := client.GetAccommodation(createReservationRequest.AccommodationID)
	if err != nil {
		tracer.LogError(span, err)
		return nil, err
	}

	for i := 0; uint(i) < createReservationRequest.NumberOfDays; i++ {
		if !s.isDateInAvailableTerms(createReservationRequest.StartDate.AddDate(0, 0, i), accommodationInfo.AvailableTerms, ctx) {
			return nil, errors.New("Accommodation is not available")
		}
	}

	var endDate = createReservationRequest.StartDate.AddDate(0, 0, int(createReservationRequest.NumberOfDays))

	acceptedReservationRequests := s.Repo.FindAcceptedReservationRequests(createReservationRequest.AccommodationID, ctx)
	for _, acceptedReservationRequest := range *acceptedReservationRequests {
		if createReservationRequest.StartDate.Before(acceptedReservationRequest.EndDate) && acceptedReservationRequest.StartDate.Before(endDate) {
			return nil, errors.New("Accomodation is reserved already")
		}
	}

	var status = model.SUBMITTED
	if accommodationInfo.AcceptReservationType == model.AUTOMATICALLY {
		status = model.ACCEPTED
	}

	var reservationRequest = model.ReservationRequest{
		StartDate:         createReservationRequest.StartDate,
		EndDate:           createReservationRequest.StartDate.AddDate(0, 0, int(createReservationRequest.NumberOfDays)),
		GuestID:           createReservationRequest.GuestID,
		GuestNumber:       createReservationRequest.GuestNumber,
		Status:            status,
		AccommodationID:   createReservationRequest.AccommodationID,
		OwnerID:           accommodationInfo.UserID,
		AccommodationName: accommodationInfo.Name}

	s.Repo.SaveReservationRequest(&reservationRequest, ctx)

	if status == model.ACCEPTED {
		resp, err := client.CreateReservedTerm(reservationRequest)
		if err == nil {
			reservationRequest.ReservedTermId = resp
			s.Repo.UpdateReservationRequestReservedTerm(&reservationRequest, ctx)
		}
	}

	return &reservationRequest, nil

}

func (s *ReservationRequestService) isDateInAvailableTerms(date time.Time, availableTerms []model.AvailableTerm, ctx context.Context) bool {
	span := tracer.StartSpanFromContext(ctx, "isDateInAvailableTermsService")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	for _, availableTerm := range availableTerms {
		if (availableTerm.StartDate.Before(date) || availableTerm.StartDate.Equal(date)) && availableTerm.EndDate.After(date) {
			return true
		}
	}

	return false
}

func (s *ReservationRequestService) GetGuestActiveReservations(guestID uint, ctx context.Context) *[]model.ReservationRequest {
	span := tracer.StartSpanFromContext(ctx, "getGuestActiveReservationsService")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	return s.Repo.FindGuestsActive(guestID, ctx)
}

func (s *ReservationRequestService) GetGuestAllReservations(guestID uint, ctx context.Context) *[]model.ReservationRequest {
	span := tracer.StartSpanFromContext(ctx, "getGuestAllReservationsService")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	return s.Repo.FindGuestsAllReservations(guestID, ctx)
}

func (s *ReservationRequestService) GetOwnersActiveReservations(ownerID uint, ctx context.Context) *[]model.ReservationRequest {
	span := tracer.StartSpanFromContext(ctx, "getOwnersActiveReservationsService")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	return s.Repo.FindOwnersActive(ownerID, ctx)
}

func (s *ReservationRequestService) GetOwnersAllReservations(ownerID uint, ctx context.Context, status []model.ReservationRequestStatus) *[]model.ReservationRequest {
	span := tracer.StartSpanFromContext(ctx, "getOwnersAllReservationsService")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	return s.Repo.FindOwnersReservations(ownerID, ctx, status)
}

func (s *ReservationRequestService) DeleteReservationRequest(reservationRequestID primitive.ObjectID, userID uint, ctx context.Context) error {
	span := tracer.StartSpanFromContext(ctx, "deleteReservationRequestService")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	reservationRequest := s.Repo.FindReservationRequest(reservationRequestID, ctx)

	if reservationRequest == nil {
		tracer.LogError(span, errors.New("Reservation request with given id does not exist."))
		return errors.New("Reservation request with given id does not exist.")
	}

	if reservationRequest.Status != model.SUBMITTED {
		tracer.LogError(span, errors.New("Reservation request can not be deleted - already submitted."))
		return errors.New("Reservation request can not be deleted - already submitted")
	}

	if reservationRequest.GuestID != userID {
		tracer.LogError(span, errors.New("You cannot access given entity."))
		return errors.New("You cannot access given entity.")
	}

	result := s.Repo.DeleteReservationRequest(reservationRequestID, ctx)
	if !result {
		tracer.LogError(span, errors.New("It's not possible to delete reservation request - repo error."))
		return errors.New("It's not possible to delete reservation request")
	}

	return nil
}

func (s *ReservationRequestService) AcceptReservationRequest(reservationRequestId primitive.ObjectID, hostId uint, ctx context.Context) (*model.ReservationRequest, error) {
	span := tracer.StartSpanFromContext(ctx, "acceptReservationRequestService")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	reservationRequest := s.Repo.FindReservationRequest(reservationRequestId, ctx)
	if reservationRequest == nil {
		tracer.LogError(span, errors.New("Given reservation request does not exist"))
		return nil, errors.New("Given reservation request does not exist")
	}

	if reservationRequest.OwnerID != hostId {
		tracer.LogError(span, errors.New("You can not acess to this entity."))
		return nil, errors.New("You can not access to this entity.")
	}

	if reservationRequest.Status != model.SUBMITTED {
		tracer.LogError(span, errors.New("You cannot update given reservation request - wrong status."))
		return nil, errors.New("You cannot update given reservation request - wrong status.")
	}

	reservationRequest.Status = model.ACCEPTED
	s.Repo.AcceptReservationRequest(reservationRequest, ctx)

	resp, err := client.CreateReservedTerm(*reservationRequest)
	if err == nil {
		reservationRequest.ReservedTermId = resp
		s.Repo.UpdateReservationRequestReservedTerm(reservationRequest, ctx)
	} else {
		tracer.LogError(span, err)
	}

	return reservationRequest, nil
}

func (s *ReservationRequestService) CancelReservationRequest(reservationRequestId primitive.ObjectID, guestId uint, ctx context.Context) (*model.ReservationRequest, error) {
	span := tracer.StartSpanFromContext(ctx, "cancelReservationRequestService")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	reservationRequest := s.Repo.FindReservationRequest(reservationRequestId, ctx)
	if reservationRequest == nil {
		tracer.LogError(span, errors.New("Given reservation request does not exist"))
		return nil, errors.New("Given reservation request does not exist")
	}

	if reservationRequest.GuestID != guestId {
		tracer.LogError(span, errors.New("You can not access to this entity"))
		return nil, errors.New("You can not access to this entity")
	}

	if reservationRequest.Status != model.ACCEPTED {
		tracer.LogError(span, errors.New("You cannot cancel given reservation request - wrong status"))
		return nil, errors.New("You cannot cancel given reservation request - wrong status")
	}

	if reservationRequest.StartDate.Before(time.Now().AddDate(0, -1, 0)) {
		tracer.LogError(span, errors.New("It's not possible to cancel reservation"))
		return nil, errors.New("It is not possible to cancel reservation.")
	}

	reservationRequest.Status = model.CANCELLED
	s.Repo.UpdateReservationRequestStatus(reservationRequest, ctx)

	resp, err := client.CreateReservedTerm(*reservationRequest)
	if err == nil {
		reservationRequest.ReservedTermId = resp
		s.Repo.UpdateReservationRequestReservedTerm(reservationRequest, ctx)
	} else {
		tracer.LogError(span, err)
	}

	client.DeleteReservedTerm(reservationRequest.ReservedTermId)

	return reservationRequest, nil
}

func (s *ReservationRequestService) CountCancelledReservations(guestId uint, ctx context.Context) int {
	span := tracer.StartSpanFromContext(ctx, "countCancelledReservationsService")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	return s.Repo.CountGuestsCancelled(guestId, ctx)
}
