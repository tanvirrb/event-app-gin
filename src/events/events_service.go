package events

import (
	"context"

	"github.com/google/uuid"
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

func (s *EventService) Get(ctx context.Context, uuid uuid.UUID) (*models.Event, error) {
	return s.repo.Get(ctx, uuid)
}

func (s *EventService) GetAll(ctx context.Context) ([]*models.Event, error) {
	return s.repo.GetAll(ctx)
}

func (s *EventService) Update(ctx context.Context, uuid uuid.UUID, e *models.Event) (*models.Event, error) {
	return s.repo.Update(ctx, uuid, e)
}

func (s *EventService) Delete(ctx context.Context, uuid uuid.UUID) error {
	return s.repo.Delete(ctx, uuid)
}
