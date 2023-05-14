package repository

import (
	"context"
	"github.com/windbnb/reservation-service/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Repository struct {
	Db *mongo.Database
}

func (r *Repository) FindAcceptedReservationRequests(accomodationId uint) *[]model.ReservationRequest {
	reservationRequests := []model.ReservationRequest{}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	filter := bson.D{
		{"accommodationID", accomodationId},
		{"status", "ACCEPTED"},
	}
	cursor, err := r.Db.Collection("reservation_request").Find(ctx, filter)

	if err != nil {
		return nil
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var reservationRequest model.ReservationRequest
		err := cursor.Decode(&reservationRequest)
		if err != nil {
			continue
		}

		reservationRequests = append(reservationRequests, reservationRequest)
	}

	return &reservationRequests
}

func (r *Repository) SaveReservationRequest(reservationRequest *model.ReservationRequest) *model.ReservationRequest {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	reservationRequest.ID = primitive.NewObjectID()
	_, err := r.Db.Collection("reservation_request").InsertOne(ctx, &reservationRequest)
	if err != nil {
		return nil
	}

	return reservationRequest
}

func (r *Repository) FindGuestsActive(guestID uint) *[]model.ReservationRequest {
	reservationRequests := []model.ReservationRequest{}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	filter := bson.D{
		{"guestID", guestID},
		{"status", "ACCEPTED"},
		{"endDate", bson.D{{"$gte", time.Now()}}},
	}
	cursor, err := r.Db.Collection("reservation_request").Find(ctx, filter)

	if err != nil {
		return nil
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var reservationRequest model.ReservationRequest
		err := cursor.Decode(&reservationRequest)
		if err != nil {
			continue
		}

		reservationRequests = append(reservationRequests, reservationRequest)
	}

	return &reservationRequests
}

func (r *Repository) FindOwnersActive(ownerID uint) *[]model.ReservationRequest {
	reservationRequests := []model.ReservationRequest{}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	filter := bson.D{
		{"ownerID", ownerID},
		{"status", "ACCEPTED"},
		{"endDate", bson.D{{"$gte", time.Now()}}},
	}
	cursor, err := r.Db.Collection("reservation_request").Find(ctx, filter)

	if err != nil {
		return nil
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var reservationRequest model.ReservationRequest
		err := cursor.Decode(&reservationRequest)
		if err != nil {
			continue
		}

		reservationRequests = append(reservationRequests, reservationRequest)
	}

	return &reservationRequests
}

func (r *Repository) DeleteReservationRequest(reservationRequestID primitive.ObjectID) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	filter := bson.D{
		{"_id", reservationRequestID},
	}

	one, err := r.Db.Collection("reservation_request").DeleteOne(ctx, filter)
	if err != nil {
		return false
	}

	return one.DeletedCount == 1
}

func (r *Repository) FindReservationRequest(reservationRequestID primitive.ObjectID) *model.ReservationRequest {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	filter := bson.D{
		{"_id", reservationRequestID},
	}

	var reservationRequest model.ReservationRequest
	err := r.Db.Collection("reservation_request").FindOne(ctx, filter).Decode(&reservationRequest)
	if err != nil {
		return nil
	}

	return &reservationRequest
}
