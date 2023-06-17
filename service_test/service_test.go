package service_test

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/windbnb/reservation-service/model"
	"github.com/windbnb/reservation-service/repository"
	"github.com/windbnb/reservation-service/service"
	"github.com/windbnb/reservation-service/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
	"time"
)

func TestDeleteReservationRequest_DoesNotExist(t *testing.T) {
	// Create a mock repository with desired behavior
	mockRepo := &MockRepo{
		FindReservationRequestFn: func(reservationRequestID primitive.ObjectID, ctx context.Context) *model.ReservationRequest {
			return nil
		},
	}

	reservationService := service.ReservationRequestService{
		Repo: mockRepo,
	}

	reservationRequest := reservationService.DeleteReservationRequest(primitive.NewObjectID(), 1, context.Background())

	assert.EqualError(t, errors.New("Reservation request with given id does not exist."), reservationRequest.Error())
}

func TestDeleteReservationRequest_ReservationRequestAccepted(t *testing.T) {
	// Create a mock repository with desired behavior
	mockRepo := &MockRepo{
		FindReservationRequestFn: func(reservationRequestID primitive.ObjectID, ctx context.Context) *model.ReservationRequest {
			return &model.ReservationRequest{
				ID:              primitive.NewObjectID(),
				StartDate:       time.Now(),
				EndDate:         time.Now(),
				AccommodationID: 1,
				GuestID:         1,
				GuestNumber:     3,
				Status:          model.ACCEPTED,
				OwnerID:         1,
				ReservedTermId:  1,
			}
		},
	}

	reservationService := service.ReservationRequestService{
		Repo: mockRepo,
	}

	reservationRequest := reservationService.DeleteReservationRequest(primitive.NewObjectID(), 1, context.Background())

	assert.EqualError(t, errors.New("Reservation request can not be deleted - only reservation request with status SUBMITTED can be deleted."), reservationRequest.Error())
}

func TestDeleteReservationRequest_ReservationRequestCancelled(t *testing.T) {
	// Create a mock repository with desired behavior
	mockRepo := &MockRepo{
		FindReservationRequestFn: func(reservationRequestID primitive.ObjectID, ctx context.Context) *model.ReservationRequest {
			return &model.ReservationRequest{
				ID:              primitive.NewObjectID(),
				StartDate:       time.Now(),
				EndDate:         time.Now(),
				AccommodationID: 1,
				GuestID:         1,
				GuestNumber:     3,
				Status:          model.CANCELLED,
				OwnerID:         1,
				ReservedTermId:  1,
			}
		},
	}

	reservationService := service.ReservationRequestService{
		Repo: mockRepo,
	}

	reservationRequest := reservationService.DeleteReservationRequest(primitive.NewObjectID(), 1, context.Background())

	assert.EqualError(t, errors.New("Reservation request can not be deleted - only reservation request with status SUBMITTED can be deleted."), reservationRequest.Error())
}

func TestDeleteReservationRequest_ReservationRequestDeclined(t *testing.T) {
	// Create a mock repository with desired behavior
	mockRepo := &MockRepo{
		FindReservationRequestFn: func(reservationRequestID primitive.ObjectID, ctx context.Context) *model.ReservationRequest {
			return &model.ReservationRequest{
				ID:              primitive.NewObjectID(),
				StartDate:       time.Now(),
				EndDate:         time.Now(),
				AccommodationID: 1,
				GuestID:         1,
				GuestNumber:     3,
				Status:          model.DECLINED,
				OwnerID:         1,
				ReservedTermId:  1,
			}
		},
	}

	reservationService := service.ReservationRequestService{
		Repo: mockRepo,
	}

	reservationRequest := reservationService.DeleteReservationRequest(primitive.NewObjectID(), 1, context.Background())

	assert.EqualError(t, errors.New("Reservation request can not be deleted - only reservation request with status SUBMITTED can be deleted."), reservationRequest.Error())
}

func TestDeleteReservationRequest_WrongGuestId(t *testing.T) {
	// Create a mock repository with desired behavior
	mockRepo := &MockRepo{
		FindReservationRequestFn: func(reservationRequestID primitive.ObjectID, ctx context.Context) *model.ReservationRequest {
			return &model.ReservationRequest{
				ID:              primitive.NewObjectID(),
				StartDate:       time.Now(),
				EndDate:         time.Now(),
				AccommodationID: 1,
				GuestID:         1,
				GuestNumber:     3,
				Status:          model.SUBMITTED,
				OwnerID:         1,
				ReservedTermId:  1,
			}
		},
	}

	reservationService := service.ReservationRequestService{
		Repo: mockRepo,
	}

	reservationRequest := reservationService.DeleteReservationRequest(primitive.NewObjectID(), 2, context.Background())

	assert.EqualError(t, errors.New("You cannot access given entity."), reservationRequest.Error())
}

func TestDeleteReservationRequest_RepoError(t *testing.T) {
	// Create a mock repository with desired behavior
	mockRepo := &MockRepo{
		FindReservationRequestFn: func(reservationRequestID primitive.ObjectID, ctx context.Context) *model.ReservationRequest {
			return &model.ReservationRequest{
				ID:              primitive.NewObjectID(),
				StartDate:       time.Now(),
				EndDate:         time.Now(),
				AccommodationID: 1,
				GuestID:         1,
				GuestNumber:     3,
				Status:          model.SUBMITTED,
				OwnerID:         1,
				ReservedTermId:  1,
			}
		},
		DeleteReservationRequestFn: func(reservationRequestID primitive.ObjectID, ctx context.Context) bool {
			return false
		},
	}

	reservationService := service.ReservationRequestService{
		Repo: mockRepo,
	}

	reservationRequest := reservationService.DeleteReservationRequest(primitive.NewObjectID(), 1, context.Background())

	assert.EqualError(t, errors.New("It's not possible to delete reservation request"), reservationRequest.Error())
}

func TestDeleteReservationRequest_Successfully(t *testing.T) {
	// Create a mock repository with desired behavior
	mockRepo := &MockRepo{
		FindReservationRequestFn: func(reservationRequestID primitive.ObjectID, ctx context.Context) *model.ReservationRequest {
			return &model.ReservationRequest{
				ID:              primitive.NewObjectID(),
				StartDate:       time.Now(),
				EndDate:         time.Now(),
				AccommodationID: 1,
				GuestID:         1,
				GuestNumber:     3,
				Status:          model.SUBMITTED,
				OwnerID:         1,
				ReservedTermId:  1,
			}
		},
		DeleteReservationRequestFn: func(reservationRequestID primitive.ObjectID, ctx context.Context) bool {
			return true
		},
	}

	reservationService := service.ReservationRequestService{
		Repo: mockRepo,
	}

	reservationRequest := reservationService.DeleteReservationRequest(primitive.NewObjectID(), 1, context.Background())

	assert.Equal(t, nil, reservationRequest)
}

func TestDeleteReservationRequest_DoesNotExist_Integration(t *testing.T) {
	// Given
	db := util.ConnectToDatabase()

	reservationService := service.ReservationRequestService{
		Repo: &repository.Repository{
			Db: db,
		},
	}

	// When
	reservationRequest := reservationService.DeleteReservationRequest(primitive.NewObjectIDFromTimestamp(time.Now()), 1, context.Background())

	// Then
	assert.EqualError(t, errors.New("Reservation request with given id does not exist."), reservationRequest.Error())
}

func TestDeleteReservationRequest_WrongStatus_Integration(t *testing.T) {
	// Given
	db := util.ConnectToDatabase()

	repo := &repository.Repository{
		Db: db,
	}
	reservationService := service.ReservationRequestService{
		Repo: repo,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	savedReservationRequest := repo.SaveReservationRequest(&model.ReservationRequest{
		ID:              primitive.NewObjectID(),
		StartDate:       time.Now(),
		EndDate:         time.Now(),
		AccommodationID: 1,
		GuestID:         1,
		GuestNumber:     3,
		Status:          model.ACCEPTED,
		OwnerID:         1,
		ReservedTermId:  1,
	}, ctx)
	objectId := savedReservationRequest.ID

	// When
	reservationRequest := reservationService.DeleteReservationRequest(objectId, 1, context.Background())

	// Then
	assert.EqualError(t, errors.New("Reservation request can not be deleted - only reservation request with status SUBMITTED can be deleted."), reservationRequest.Error())
}

func TestDeleteReservationRequest_Successfully_Integration(t *testing.T) {
	// Given
	db := util.ConnectToDatabase()

	repo := &repository.Repository{
		Db: db,
	}
	reservationService := service.ReservationRequestService{
		Repo: repo,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	savedReservationRequest := repo.SaveReservationRequest(&model.ReservationRequest{
		ID:              primitive.NewObjectID(),
		StartDate:       time.Now(),
		EndDate:         time.Now(),
		AccommodationID: 1,
		GuestID:         1,
		GuestNumber:     3,
		Status:          model.SUBMITTED,
		OwnerID:         1,
		ReservedTermId:  1,
	}, ctx)
	objectId := savedReservationRequest.ID

	// When
	reservationRequest := reservationService.DeleteReservationRequest(objectId, 1, context.Background())

	// Then
	assert.Equal(t, nil, reservationRequest)
}

type MockRepo struct {
	repository.Repository
	FindReservationRequestFn   func(reservationRequestID primitive.ObjectID, ctx context.Context) *model.ReservationRequest
	DeleteReservationRequestFn func(reservationRequestID primitive.ObjectID, ctx context.Context) bool
}

func (m *MockRepo) FindReservationRequest(reservationRequestId primitive.ObjectID, ctx context.Context) *model.ReservationRequest {
	return m.FindReservationRequestFn(reservationRequestId, ctx)
}

func (m *MockRepo) DeleteReservationRequest(reservationRequestID primitive.ObjectID, ctx context.Context) bool {
	return m.DeleteReservationRequestFn(reservationRequestID, ctx)
}
