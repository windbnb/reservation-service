package service

import (
	"errors"
	"github.com/windbnb/reservation-service/client"
	"github.com/windbnb/reservation-service/model"
	"github.com/windbnb/reservation-service/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
	"time"
)

var accommodationServiceUrl = os.Getenv("accommodationServiceUrl") + "/api/accomodation/"

type ReservationRequestService struct {
	Repo *repository.Repository
}

func (s *ReservationRequestService) SaveReservationRequest(createReservationRequest *model.CreateReservationRequest) (*model.ReservationRequest, error) {
	if createReservationRequest.StartDate.Before(time.Now()) {
		return nil, errors.New("Start date cannot be in past")
	}

	if createReservationRequest.NumberOfDays <= 0 {
		return nil, errors.New("Number of days must be positive")
	}

	accommodationInfo, err := client.GetAccommodation(createReservationRequest.AccommodationID)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	for i := 0; uint(i) < createReservationRequest.NumberOfDays; i++ {
		if !s.isDateInAvailableTerms(createReservationRequest.StartDate.AddDate(0, 0, i), accommodationInfo.AvailableTerms) {
			return nil, errors.New("Accommodation is not available")
		}
	}

	var endDate = createReservationRequest.StartDate.AddDate(0, 0, int(createReservationRequest.NumberOfDays))
	acceptedReservationRequests := s.Repo.FindAcceptedReservationRequests(createReservationRequest.AccommodationID)
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
		StartDate:       createReservationRequest.StartDate,
		EndDate:         createReservationRequest.StartDate.AddDate(0, 0, int(createReservationRequest.NumberOfDays)),
		GuestID:         createReservationRequest.GuestID,
		GuestNumber:     createReservationRequest.GuestNumber,
		Status:          status,
		AccommodationID: createReservationRequest.AccommodationID,
		OwnerID:         accommodationInfo.UserID}

	return s.Repo.SaveReservationRequest(&reservationRequest), nil
}

func (s *ReservationRequestService) isDateInAvailableTerms(date time.Time, availableTerms []model.AvailableTerm) bool {
	for _, availableTerm := range availableTerms {
		if (availableTerm.StartDate.Before(date) || availableTerm.StartDate.Equal(date)) && availableTerm.EndDate.After(date) {
			return true
		}
	}

	return false
}

func (s *ReservationRequestService) GetGuestActiveReservations(guestID uint) *[]model.ReservationRequest {
	return s.Repo.FindGuestsActive(guestID)
}

func (s *ReservationRequestService) GetOwnersActiveReservations(ownerID uint) *[]model.ReservationRequest {
	return s.Repo.FindOwnersActive(ownerID)
}

func (s *ReservationRequestService) DeleteReservationRequest(reservationRequestID primitive.ObjectID, userID uint) error {
	reservationRequest := s.Repo.FindReservationRequest(reservationRequestID)

	if reservationRequest == nil {
		return errors.New("Reservation request with given id does not exist")
	}

	if reservationRequest.Status != model.SUBMITTED {
		return errors.New("Reservation request can not be deleted")
	}

	if reservationRequest.GuestID != userID {
		return errors.New("You cannot access given entity.")
	}

	result := s.Repo.DeleteReservationRequest(reservationRequestID)
	if !result {
		return errors.New("It's not possible to delete reservation request")
	}

	client.DeleteReservedTerm(reservationRequest.ReservedTermId)

	return nil
}

func (s *ReservationRequestService) AcceptReservationRequest(reservationRequestId primitive.ObjectID, hostId uint) (*model.ReservationRequest, error) {
	reservationRequest := s.Repo.FindReservationRequest(reservationRequestId)
	if reservationRequest == nil {
		return nil, errors.New("Given reservation request does not exist")
	}

	if reservationRequest.OwnerID != hostId {
		return nil, errors.New("You can not access to this entity")
	}

	if reservationRequest.Status != model.SUBMITTED {
		return nil, errors.New("You cannot update given reservation request")
	}

	reservationRequest.Status = model.ACCEPTED
	s.Repo.AcceptReservationRequest(reservationRequest)

	resp, err := client.CreateReservedTerm(*reservationRequest)
	if err == nil {
		reservationRequest.ReservedTermId = resp
		s.Repo.UpdateReservationRequest(reservationRequest)
	}

	return reservationRequest, nil
}
