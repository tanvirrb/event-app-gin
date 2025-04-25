package events

import (
	"context"
	"fmt"

	"github.com/tanvirrb/event-app-gin/src/events/interfaces"
	"github.com/tanvirrb/event-app-gin/src/events/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDBEventRepository struct {
	collection *mongo.Collection
}

func NewMongoDBEventRepository(collection *mongo.Collection) interfaces.EventRepository {
	return &MongoDBEventRepository{
		collection: collection,
	}
}

func (r *MongoDBEventRepository) Create(ctx context.Context, event *models.Event) (*models.Event, error) {
	result, err := r.collection.InsertOne(ctx, event)
	if err != nil {
		return nil, fmt.Errorf("failed to create event: %v", err)
	}

	event.Id = result.InsertedID.(primitive.ObjectID)
	return event, nil
}

func (r *MongoDBEventRepository) Get(ctx context.Context, id string) (*models.Event, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id format: %v", err)
	}

	var event models.Event
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&event)
	if err != nil {
		return nil, fmt.Errorf("failed to get event: %v", err)
	}

	return &event, nil
}

func (r *MongoDBEventRepository) GetAll(ctx context.Context) ([]*models.Event, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %v", err)
	}
	defer cursor.Close(ctx)

	var events []*models.Event
	if err = cursor.All(ctx, &events); err != nil {
		return nil, fmt.Errorf("failed to decode events: %v", err)
	}

	return events, nil
}

func (r *MongoDBEventRepository) Update(ctx context.Context, id string, event *models.Event) (*models.Event, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id format: %v", err)
	}

	update := bson.M{
		"$set": bson.M{
			"name":  event.Name,
			"genre": event.Genre,
		},
	}

	result := r.collection.FindOneAndUpdate(ctx, bson.M{"_id": objectID}, update)
	if result.Err() != nil {
		return nil, fmt.Errorf("failed to update event: %v", result.Err())
	}

	event.Id = objectID
	return event, nil
}

func (r *MongoDBEventRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid id format: %v", err)
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("failed to delete event: %v", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("event not found")
	}

	return nil
}
