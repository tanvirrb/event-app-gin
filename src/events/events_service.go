package events

import (
	"context"

	"github.com/tanvirrb/event-app-gin/src/events/interfaces"
	"github.com/tanvirrb/event-app-gin/src/events/models"
)

type EventService struct {
	repo interfaces.EventRepository
}

func NewEventService(repo interfaces.EventRepository) interfaces.EventService {
	return &EventService{
		repo: repo,
	}
}

func (s *EventService) Create(ctx context.Context, e *models.Event) (*models.Event, error) {
	return s.repo.Create(ctx, e)
}

func (s *EventService) Get(ctx context.Context, id string) (*models.Event, error) {
	return s.repo.Get(ctx, id)
}

func (s *EventService) GetAll(ctx context.Context) ([]*models.Event, error) {
	return s.repo.GetAll(ctx)
}

func (s *EventService) Update(ctx context.Context, id string, e *models.Event) (*models.Event, error) {
	return s.repo.Update(ctx, id, e)
}

func (s *EventService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
