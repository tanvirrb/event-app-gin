package interfaces

import (
	"context"

	"github.com/google/uuid"
	"github.com/tanvirrb/event-app-gin/src/events/models"
)

type EventRepository interface {
	Create(ctx context.Context, event *models.Event) (*models.Event, error)

	Get(ctx context.Context, uuid uuid.UUID) (*models.Event, error)

	GetAll(ctx context.Context) ([]*models.Event, error)

	Update(ctx context.Context, uuid uuid.UUID, event *models.Event) (*models.Event, error)

	Delete(ctx context.Context, uuid uuid.UUID) error
}
