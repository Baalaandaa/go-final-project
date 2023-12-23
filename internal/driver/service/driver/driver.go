package driver

import (
	"context"
	"final-project/internal/driver/model"
	"final-project/internal/driver/service"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	StatusDriverSearch = "DRIVER_SEARCH"
	StatusDriverFound  = "DRIVER_FOUND"
	StatusOnPosition   = "ON_POSITION"
	StatusStarted      = "STARTED"
	StatusEnded        = "ENDED"
	StatusCanceled     = "CANCELED"
)

type driverService struct {
	db             *mongo.Database
	collectionName string
}

func (d driverService) ListTrips(ctx context.Context, userId string) (*[]model.Trip, error) {
	collection := d.db.Collection(d.collectionName)

	filter := bson.M{"driver_id": userId}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var result []model.Trip
	for cursor.Next(ctx) {
		var trip model.Trip
		err := cursor.Decode(&trip)
		if err != nil {
			log.Println("Error decoding trip:", err)
			continue
		}
		result = append(result, trip)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return &result, nil
}

func (d driverService) GetTrip(ctx context.Context, userId string, tripId string) (*model.Trip, error) {
	collection := d.db.Collection(d.collectionName)

	var result *model.Trip
	filter := bson.M{"id": tripId, "driver_id": userId}
	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, service.ErrNotFound
		}
		return nil, err
	}

	return result, nil
}

func (d driverService) CancelTrip(ctx context.Context, userId string, tripId string, reason *string) error {
	return d.changeTripStatus(ctx, userId, tripId, StatusCanceled, reason)
}

func (d driverService) AcceptTrip(ctx context.Context, userId string, tripId string) error {
	return d.changeTripStatus(ctx, userId, tripId, StatusDriverFound, nil)
}

func (d driverService) StartTrip(ctx context.Context, userId string, tripId string) error {
	return d.changeTripStatus(ctx, userId, tripId, StatusStarted, nil)
}

func (d driverService) EndTrip(ctx context.Context, userId string, tripId string) error {
	return d.changeTripStatus(ctx, userId, tripId, StatusEnded, nil)
}

func (d *driverService) changeTripStatus(ctx context.Context, userId string, tripId string, newStatus string, cancelReason *string) error {
	collection := d.db.Collection(d.collectionName)

	updateData := bson.M{"status": newStatus}
	if cancelReason != nil {
		updateData["cancel_reason"] = *cancelReason
	}

	filter := bson.M{"id": tripId, "driver_id": userId}
	update := bson.M{"$set": updateData}

	result, err := collection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(false))
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return service.ErrNotFound
	}
	return nil
}

func New(db *mongo.Database) service.Driver {
	return &driverService{db: db}
}
