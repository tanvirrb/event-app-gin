package events

import (
	"errors"
	"github.com/google/uuid"
	"github.com/tanvirrb/event-app-go/src/events/interfaces"
	"github.com/tanvirrb/event-app-go/src/events/models"
)

type EventRepository struct {
	storage map[string]*models.Event
}

func NewEventRepository() interfaces.EventRepository {
	return &EventRepository{
		storage: make(map[string]*models.Event),
	}
}

func (r *EventRepository) Create(event *models.Event) (*models.Event, error) {
	event.ID = uuid.New().String()
	r.storage[event.ID] = event
	return event, nil
}

func (r *EventRepository) Get(id string) (*models.Event, error) {
	event, exists := r.storage[id]
	if !exists {
		return nil, errors.New("event not found")
	}
	return event, nil
}

func (r *EventRepository) GetAll() ([]*models.Event, error) {
	var events []*models.Event
	for _, event := range r.storage {
		events = append(events, event)
	}
	return events, nil
}

func (r *EventRepository) Update(id string, event *models.Event) (*models.Event, error) {
	_, exists := r.storage[id]
	if !exists {
		return nil, errors.New("event not found")
	}
	event.ID = id
	r.storage[id] = event
	return event, nil
}

func (r *EventRepository) Delete(id string) error {
	_, exists := r.storage[id]
	if !exists {
		return errors.New("event not found")
	}
	delete(r.storage, id)
	return nil
}
