package events

import (
	"context"
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/tanvirrb/event-app-go/src/events/interfaces"
	"github.com/tanvirrb/event-app-go/src/events/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type EventRepository struct {
	collection *mongo.Collection
}

func NewEventRepository(collection *mongo.Collection) interfaces.EventRepository {
	return &EventRepository{
		collection: collection,
	}
}

func (r *EventRepository) Create(event *models.Event) (*models.Event, error) {
	encodedId, err := r.collection.InsertOne(context.TODO(), event)
	if err != nil {
		log.Printf("Error while creating event: %v", err)
		return nil, err
	}

	eventId := encodedId.InsertedID.(primitive.ObjectID)

	var createdEvent models.Event
	err = r.collection.FindOne(context.Background(), primitive.M{"_id": eventId}).Decode(&createdEvent)
	if err != nil {
		log.Printf("Error while fetching created event: %v", err)
		return nil, err
	}

	return &createdEvent, nil
}

func (r *EventRepository) Get(id string) (*models.Event, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("Invalid Object ID: %v", err)
		return nil, err
	}

	var event models.Event
	err = r.collection.FindOne(context.Background(), primitive.M{"_id": objectId}).Decode(&event)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Printf("No event found with id: %v", id)
			return nil, nil
		}

		log.Printf("Error while fetching event: %v", err)
		return nil, err
	}
	return &event, nil
}

func (r *EventRepository) GetAll() ([]*models.Event, error) {
	ctx := context.Background()

	eventListCursor, err := r.collection.Find(ctx, primitive.M{})
	if err != nil {
		log.Printf("Error while fetching events: %v", err)
		return nil, err
	}

	var events []*models.Event
	if err = eventListCursor.All(ctx, &events); err != nil {
		log.Printf("Error while fetching events: %v", err)
		return nil, err
	}
	return events, nil
}

func (r *EventRepository) Update(id string, event *models.Event) (*models.Event, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("Invalid Object ID: %v", err)
		return nil, err
	}

	result := r.collection.FindOneAndReplace(context.Background(), primitive.M{"_id": objectId}, event, options.FindOneAndReplace().SetReturnDocument(options.After))
	if result.Err() != nil {
		return nil, result.Err()
	}

	var updatedEvent models.Event
	err = result.Decode(&updatedEvent)
	if err != nil {
		return nil, err
	}

	return &updatedEvent, nil
}

//func (r *EventRepository) Delete(id string) error {
//	_, exists := r.storage[id]
//	if !exists {
//		return errors.New("event not found")
//	}
//	delete(r.storage, id)
//	return nil
//}
