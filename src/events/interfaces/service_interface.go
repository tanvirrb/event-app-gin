package interfaces

import (
	"context"

	"github.com/tanvirrb/event-app-gin/src/events/models"
)

type EventService interface {
	Create(ctx context.Context, event *models.Event) (*models.Event, error)
	Get(ctx context.Context, id string) (*models.Event, error)
	GetAll(ctx context.Context) ([]*models.Event, error)
	Update(ctx context.Context, id string, event *models.Event) (*models.Event, error)
	Delete(ctx context.Context, id string) error
}
