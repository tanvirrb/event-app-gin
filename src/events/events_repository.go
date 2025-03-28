package events

import (
	"context"
	"errors"
	"github.com/tanvirrb/event-app-go/src/configs"
	"github.com/tanvirrb/event-app-go/src/events/interfaces"
	"github.com/tanvirrb/event-app-go/src/events/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type EventRepository struct {
	collection *mongo.Collection
}

func NewEventRepository() interfaces.EventRepository {
	eventCollection := configs.GetCollection(configs.DB, "events")
	return &EventRepository{
		collection: eventCollection,
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
	defer func(eventListCursor *mongo.Cursor, ctx context.Context) {
		err := eventListCursor.Close(ctx)
		if err != nil {
			log.Printf("Error while closing cursor: %v", err)
		}
	}(eventListCursor, ctx)

	var events []*models.Event
	for eventListCursor.Next(ctx) {
		var event models.Event
		err := eventListCursor.Decode(&event)
		if err != nil {
			log.Printf("Error while decoding event: %v", err)
			return nil, err
		}
		events = append(events, &event)
	}
	return events, nil
}

//func (r *EventRepository) Update(id string, event *models.Event) (*models.Event, error) {
//	_, exists := r.storage[id]
//	if !exists {
//		return nil, errors.New("event not found")
//	}
//	event.ID = id
//	r.storage[id] = event
//	return event, nil
//}
//
//func (r *EventRepository) Delete(id string) error {
//	_, exists := r.storage[id]
//	if !exists {
//		return errors.New("event not found")
//	}
//	delete(r.storage, id)
//	return nil
//}
