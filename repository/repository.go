package repository

import (
	"context"
	"github.com/windbnb/reservation-service/model"
	"github.com/windbnb/reservation-service/tracer"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Repository struct {
	Db *mongo.Database
}

func (r *Repository) FindAcceptedReservationRequests(accomodationId uint, ctx context.Context) *[]model.ReservationRequest {
	span := tracer.StartSpanFromContext(ctx, "findReservationRequestsRepository")
	defer span.Finish()

	reservationRequests := []model.ReservationRequest{}
	dbCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	filter := bson.D{
		{"accommodationID", accomodationId},
		{"status", model.ACCEPTED},
	}

	cursor, err := r.Db.Collection("reservation_request").Find(dbCtx, filter)
	if err != nil {
		tracer.LogError(span, err)
		return nil
	}

	for cursor.Next(dbCtx) {
		var reservationRequest model.ReservationRequest
		err := cursor.Decode(&reservationRequest)
		if err != nil {
			tracer.LogError(span, err)
			continue
		}

		reservationRequests = append(reservationRequests, reservationRequest)
	}

	return &reservationRequests
}

func (r *Repository) SaveReservationRequest(reservationRequest *model.ReservationRequest, ctx context.Context) *model.ReservationRequest {
	span := tracer.StartSpanFromContext(ctx, "saveReservationRequestRepository")
	defer span.Finish()

	dbCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	reservationRequest.ID = primitive.NewObjectID()
	_, err := r.Db.Collection("reservation_request").InsertOne(dbCtx, &reservationRequest)
	if err != nil {
		tracer.LogError(span, err)
		return nil
	}

	return reservationRequest
}

func (r *Repository) FindGuestsActive(guestID uint, ctx context.Context) *[]model.ReservationRequest {
	span := tracer.StartSpanFromContext(ctx, "findGuestsActiveRepository")
	defer span.Finish()

	reservationRequests := []model.ReservationRequest{}
	dbCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	filter := bson.D{
		{"guestID", guestID},
		{"status", model.ACCEPTED},
		{"endDate", bson.D{{"$gte", time.Now()}}},
	}
	cursor, err := r.Db.Collection("reservation_request").Find(dbCtx, filter)

	if err != nil {
		tracer.LogError(span, err)
		return nil
	}
	defer cursor.Close(dbCtx)

	for cursor.Next(dbCtx) {
		var reservationRequest model.ReservationRequest
		err := cursor.Decode(&reservationRequest)
		if err != nil {
			tracer.LogError(span, err)
			continue
		}

		reservationRequests = append(reservationRequests, reservationRequest)
	}

	return &reservationRequests
}

func (r *Repository) FindGuestsAllReservations(guestID uint, ctx context.Context) *[]model.ReservationRequest {
	span := tracer.StartSpanFromContext(ctx, "findGuestsAllRepository")
	defer span.Finish()

	reservationRequests := []model.ReservationRequest{}
	dbCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	filter := bson.D{
		{"guestID", guestID},
	}
	cursor, err := r.Db.Collection("reservation_request").Find(dbCtx, filter)

	if err != nil {
		tracer.LogError(span, err)
		return nil
	}
	defer cursor.Close(dbCtx)

	for cursor.Next(dbCtx) {
		var reservationRequest model.ReservationRequest
		err := cursor.Decode(&reservationRequest)
		if err != nil {
			tracer.LogError(span, err)
			continue
		}

		reservationRequests = append(reservationRequests, reservationRequest)
	}

	return &reservationRequests
}

func (r *Repository) FindOwnersActive(ownerID uint, ctx context.Context) *[]model.ReservationRequest {
	span := tracer.StartSpanFromContext(ctx, "findOwnersActiveRepository")
	defer span.Finish()

	reservationRequests := []model.ReservationRequest{}
	dbCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	filter := bson.D{
		{"ownerID", ownerID},
		{"status", model.ACCEPTED},
		{"endDate", bson.D{{"$gte", time.Now()}}},
	}
	cursor, err := r.Db.Collection("reservation_request").Find(dbCtx, filter)

	if err != nil {
		tracer.LogError(span, err)
		return nil
	}
	defer cursor.Close(dbCtx)

	for cursor.Next(dbCtx) {
		var reservationRequest model.ReservationRequest
		err := cursor.Decode(&reservationRequest)
		if err != nil {
			continue
		}

		reservationRequests = append(reservationRequests, reservationRequest)
	}

	return &reservationRequests
}

func (r *Repository) FindOwnersSubmitted(ownerID uint, ctx context.Context) *[]model.ReservationRequest {
	span := tracer.StartSpanFromContext(ctx, "findOwnersSubmittedRepository")
	defer span.Finish()

	reservationRequests := []model.ReservationRequest{}
	dbCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	filter := bson.D{
		{"ownerID", ownerID},
		{"status", model.SUBMITTED},
	}
	cursor, err := r.Db.Collection("reservation_request").Find(dbCtx, filter)

	if err != nil {
		tracer.LogError(span, err)
		return nil
	}
	defer cursor.Close(dbCtx)

	for cursor.Next(dbCtx) {
		var reservationRequest model.ReservationRequest
		err := cursor.Decode(&reservationRequest)
		if err != nil {
			continue
		}

		reservationRequests = append(reservationRequests, reservationRequest)
	}

	return &reservationRequests
}

func (r *Repository) FindOwnersAll(ownerID uint, ctx context.Context) *[]model.ReservationRequest {
	span := tracer.StartSpanFromContext(ctx, "findOwnersSubmittedRepository")
	defer span.Finish()

	reservationRequests := []model.ReservationRequest{}
	dbCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	filter := bson.D{
		{"ownerID", ownerID},
	}
	cursor, err := r.Db.Collection("reservation_request").Find(dbCtx, filter)

	if err != nil {
		tracer.LogError(span, err)
		return nil
	}
	defer cursor.Close(dbCtx)

	for cursor.Next(dbCtx) {
		var reservationRequest model.ReservationRequest
		err := cursor.Decode(&reservationRequest)
		if err != nil {
			continue
		}

		reservationRequests = append(reservationRequests, reservationRequest)
	}

	return &reservationRequests
}

func (r *Repository) DeleteReservationRequest(reservationRequestID primitive.ObjectID, ctx context.Context) bool {
	span := tracer.StartSpanFromContext(ctx, "saveAccomodationRepository")
	defer span.Finish()

	dbCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	filter := bson.D{
		{"_id", reservationRequestID},
	}

	one, err := r.Db.Collection("reservation_request").DeleteOne(dbCtx, filter)
	if err != nil {
		tracer.LogError(span, err)
		return false
	}

	return one.DeletedCount == 1
}

func (r *Repository) FindReservationRequest(reservationRequestID primitive.ObjectID, ctx context.Context) *model.ReservationRequest {
	span := tracer.StartSpanFromContext(ctx, "findReservationRequestRepository")
	defer span.Finish()

	dbCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	filter := bson.D{
		{"_id", reservationRequestID},
	}

	var reservationRequest model.ReservationRequest
	err := r.Db.Collection("reservation_request").FindOne(dbCtx, filter).Decode(&reservationRequest)
	if err != nil {
		tracer.LogError(span, err)
		return nil
	}

	return &reservationRequest
}

func (r *Repository) AcceptReservationRequest(reservationRequest *model.ReservationRequest, ctx context.Context) *model.ReservationRequest {
	span := tracer.StartSpanFromContext(ctx, "acceptReservationRequestRepository")
	defer span.Finish()

	dbCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	filter := bson.D{
		{
			"$and",
			bson.A{
				bson.D{
					{
						"$or",
						bson.A{
							bson.D{
								{"$and", bson.A{
									bson.D{{"startDate", bson.D{{"$gte", reservationRequest.StartDate}}}},
									bson.D{{"startDate", bson.D{{"$lte", reservationRequest.EndDate}}}},
								}},
							},
							bson.D{
								{"$and", bson.A{
									bson.D{{"endDate", bson.D{{"$gte", reservationRequest.StartDate}}}},
									bson.D{{"endDate", bson.D{{"$lte", reservationRequest.EndDate}}}},
								}},
							},
						},
					},
				},
				bson.D{
					{
						"status", model.SUBMITTED,
					},
				},
			},
		},
	}

	declinedReservationRequest := bson.D{{"$set", bson.D{{"status", model.DECLINED}}}}
	acceptedReservationRequest := bson.D{{"$set", bson.D{{"status", model.ACCEPTED}}}}

	_, err := r.Db.Collection("reservation_request").UpdateMany(dbCtx, filter, declinedReservationRequest)
	_, err = r.Db.Collection("reservation_request").UpdateByID(dbCtx, reservationRequest.ID, acceptedReservationRequest)
	if err != nil {
		tracer.LogError(span, err)
		return nil
	}

	return reservationRequest
}

func (r *Repository) UpdateReservationRequestReservedTerm(reservationRequest *model.ReservationRequest, ctx context.Context) *model.ReservationRequest {
	span := tracer.StartSpanFromContext(ctx, "updateReservationRequestReservedTermRepository")
	defer span.Finish()

	updateQuery := bson.D{{"$set", bson.D{{"reservedTermId", reservationRequest.ReservedTermId}}}}
	err := r.updateReservationRequest(reservationRequest, updateQuery, ctx)
	if err != nil {
		tracer.LogError(span, err)
		return nil
	}

	return reservationRequest
}

func (r *Repository) UpdateReservationRequestStatus(reservationRequest *model.ReservationRequest, ctx context.Context) *model.ReservationRequest {
	span := tracer.StartSpanFromContext(ctx, "updateReservationRequestStatusRepository")
	defer span.Finish()

	updateQuery := bson.D{{"$set", bson.D{{"status", model.CANCELLED}}}}
	err := r.updateReservationRequest(reservationRequest, updateQuery, ctx)
	if err != nil {
		tracer.LogError(span, err)
		return nil
	}

	return reservationRequest
}

func (r *Repository) updateReservationRequest(reservationRequest *model.ReservationRequest, updateQuery bson.D, ctx context.Context) error {
	span := tracer.StartSpanFromContext(ctx, "updateReservationRequestRepository")
	defer span.Finish()

	dbCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, err := r.Db.Collection("reservation_request").UpdateByID(dbCtx, reservationRequest.ID, updateQuery)
	return err
}

func (r *Repository) CountGuestsCancelled(guestId uint, ctx context.Context) int {
	span := tracer.StartSpanFromContext(ctx, "saveAccomodationRepository")
	defer span.Finish()

	dbCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	filter := bson.D{
		{"guestID", guestId},
		{"status", model.CANCELLED},
	}
	count, err := r.Db.Collection("reservation_request").CountDocuments(dbCtx, filter)

	if err != nil {
		tracer.LogError(span, err)
		return 0
	}

	return int(count)
}
