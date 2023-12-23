package location_repo

import (
	"context"
	"final-project/internal/driver/model"
	"final-project/internal/driver/repository"
	"final-project/pkg/helpers"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type driverRepo struct {
	db *mongo.Database
}

func (d driverRepo) GetTripsList(ctx context.Context, userId string) (*[]model.Trip, error) {
	collection := d.db.Collection("trips")

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
			log.Fatal("Error decoding trip:", err)
			continue
		}
		result = append(result, trip)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return &result, nil
}

func (d driverRepo) GetTrip(ctx context.Context, userId string, tripId string) (*model.Trip, error) {
	collection := d.db.Collection("trips")

	var result *model.Trip
	filter := bson.M{"id": tripId, "driver_id": userId}
	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, helpers.ErrNotFound
		}
		return nil, err
	}

	return result, nil
}

func (d driverRepo) ChangeTripStatus(ctx context.Context, userId string, tripId string, newStatus string, cancelReason *string) error {
	collection := d.db.Collection("trips")

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
		return helpers.ErrNotFound
	}
	return nil
}

func New(db *mongo.Database, collectionName string) repository.Driver {
	return &driverRepo{
		db: db,
	}
}
