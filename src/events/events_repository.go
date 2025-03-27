package events

import (
	"context"
	"fmt"
	"github.com/tanvirrb/event-app-go/src/configs"
	"github.com/tanvirrb/event-app-go/src/events/interfaces"
	"github.com/tanvirrb/event-app-go/src/events/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type EventRepository struct {
	collection *mongo.Collection
}

func NewEventRepository() interfaces.EventRepository {
	var eventCollection *mongo.Collection = configs.GetCollection(configs.DB, "events")
	return &EventRepository{
		collection: eventCollection,
	}
}

func GetObjectId(result *mongo.InsertOneResult) string {
	return result.InsertedID.(primitive.ObjectID).Hex()

}

func (r *EventRepository) Create(event *models.Event) (*models.Event, error) {
	encodedId, err := r.collection.InsertOne(context.TODO(), event)
	if err != nil {
		return nil, err
	}

	eventId := encodedId.InsertedID.(primitive.ObjectID)
	fmt.Printf("Event _Id after save: %v", eventId)
	var createdEvent models.Event
	err = r.collection.FindOne(context.Background(), primitive.M{"_id": eventId}).Decode(&createdEvent)
	if err != nil {
		return nil, err
	}

	return &createdEvent, nil
}

func (r *EventRepository) Get(id string) (*models.Event, error) {
	var event models.Event
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	err = r.collection.FindOne(context.Background(), primitive.M{"_id": objectId}).Decode(&event)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Event: %v", event)
	return &event, nil
}

//
//func (r *EventRepository) GetAll() ([]*models.Event, error) {
//	var events []*models.Event
//	for _, event := range r.storage {
//		events = append(events, event)
//	}
//	return events, nil
//}
//
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
