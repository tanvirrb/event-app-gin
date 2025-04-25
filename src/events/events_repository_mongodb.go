package events

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/tanvirrb/event-app-gin/src/events/interfaces"
	"github.com/tanvirrb/event-app-gin/src/events/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDBEventRepository struct {
	collection *mongo.Collection
}

func NewMongoDBEventRepository(collection *mongo.Collection) interfaces.EventRepository {
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "uuid", Value: 1}},
		Options: nil,
	}
	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		fmt.Printf("Failed to create index on uuid field: %v\n", err)
	}

	return &MongoDBEventRepository{
		collection: collection,
	}
}

func (r *MongoDBEventRepository) Create(ctx context.Context, event *models.Event) (*models.Event, error) {
	event.Uuid = uuid.New()
	_, err := r.collection.InsertOne(ctx, event)
	if err != nil {
		return nil, fmt.Errorf("failed to create event: %v", err)
	}

	return event, nil
}

func (r *MongoDBEventRepository) Get(ctx context.Context, uuid uuid.UUID) (*models.Event, error) {
	var event models.Event
	err := r.collection.FindOne(ctx, bson.M{"uuid": uuid}).Decode(&event)
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

func (r *MongoDBEventRepository) Update(ctx context.Context, uuid uuid.UUID, event *models.Event) (*models.Event, error) {
	update := bson.M{
		"$set": bson.M{
			"name":  event.Name,
			"genre": event.Genre,
		},
	}

	result := r.collection.FindOneAndUpdate(ctx, bson.M{"uuid": uuid}, update)
	if result.Err() != nil {
		return nil, fmt.Errorf("failed to update event: %v", result.Err())
	}

	event.Uuid = uuid
	return event, nil
}

func (r *MongoDBEventRepository) Delete(ctx context.Context, uuid uuid.UUID) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"uuid": uuid})
	if err != nil {
		return fmt.Errorf("failed to delete event: %v", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("event not found")
	}

	return nil
}
