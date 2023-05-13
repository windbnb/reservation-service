package service

import (
	"encoding/json"
	"errors"
	"github.com/windbnb/reservation-service/model"
	"github.com/windbnb/reservation-service/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"os"
	"strconv"
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

	client := &http.Client{}

	req, _ := http.NewRequest("GET", accommodationServiceUrl+strconv.FormatUint(uint64(createReservationRequest.AccommodationID), 10), nil)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	var accommodationInfo = &model.AccommodationInfo{}
	decoderErr := json.NewDecoder(resp.Body).Decode(&accommodationInfo)

	if decoderErr != nil {
		return nil, decoderErr
	}

	for i := 0; uint(i) < createReservationRequest.NumberOfDays; i++ {
		if !s.isDateInAvailableTerms(createReservationRequest.StartDate.AddDate(0, 0, i), accommodationInfo.AvailableTerms) {
			return nil, errors.New("Accommodaion is not available")
		}
	}

	var endDate = createReservationRequest.StartDate.AddDate(0, 0, int(createReservationRequest.NumberOfDays))
	acceptedReservationRequests := s.Repo.FindAcceptedReservationRequests(createReservationRequest.AccommodationID)
	for _, acceptedReservationRequest := range *acceptedReservationRequests {
		if createReservationRequest.StartDate.Before(acceptedReservationRequest.EndDate) && acceptedReservationRequest.StartDate.Before(endDate) {
			return nil, errors.New("Accomodation is reserved already")
		}
	}

	var reservationRequest = model.ReservationRequest{
		StartDate:       createReservationRequest.StartDate,
		EndDate:         createReservationRequest.StartDate.AddDate(0, 0, int(createReservationRequest.NumberOfDays)),
		GuestID:         createReservationRequest.GuestID,
		GuestNumber:     createReservationRequest.GuestNumber,
		Status:          model.SUBMITTED,
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

func (s *ReservationRequestService) DeleteReservationRequest(reservationRequestID primitive.ObjectID) error {
	reservationRequest := s.Repo.FindReservationRequest(reservationRequestID)

	if reservationRequest == nil {
		return errors.New("Reservation request with given id does not exist")
	}

	if reservationRequest.Status != model.SUBMITTED {
		return errors.New("Reservation request can not be deleted")
	}

	result := s.Repo.DeleteReservationRequest(reservationRequestID)
	if !result {
		return errors.New("It's not possible to delete reservation request")
	}

	return nil
}
