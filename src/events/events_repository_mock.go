package events

import (
	"context"
	"fmt"
	"sync"

	"github.com/tanvirrb/event-app-gin/src/events/interfaces"
	"github.com/tanvirrb/event-app-gin/src/events/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockEventRepository struct {
	events map[string]*models.Event
	mu     sync.RWMutex
}

func NewMockEventRepository() interfaces.EventRepository {
	return &MockEventRepository{
		events: make(map[string]*models.Event),
	}
}

func (r *MockEventRepository) Create(ctx context.Context, event *models.Event) (*models.Event, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	event.Id = primitive.NewObjectID()
	r.events[event.Id.Hex()] = event
	return event, nil
}

func (r *MockEventRepository) Get(ctx context.Context, id string) (*models.Event, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	event, exists := r.events[id]
	if !exists {
		return nil, fmt.Errorf("event not found")
	}
	return event, nil
}

func (r *MockEventRepository) GetAll(ctx context.Context) ([]*models.Event, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	events := make([]*models.Event, 0, len(r.events))
	for _, event := range r.events {
		events = append(events, event)
	}
	return events, nil
}

func (r *MockEventRepository) Update(ctx context.Context, id string, event *models.Event) (*models.Event, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.events[id]; !exists {
		return nil, fmt.Errorf("event not found")
	}

	event.Id, _ = primitive.ObjectIDFromHex(id)
	r.events[id] = event
	return event, nil
}

func (r *MockEventRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.events[id]; !exists {
		return fmt.Errorf("event not found")
	}

	delete(r.events, id)
	return nil
}
