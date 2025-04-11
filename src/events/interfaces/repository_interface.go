package interfaces

import (
	"github.com/tanvirrb/event-app-gin/src/events/models"
)

type EventRepository interface {
	Create(*models.Event) (*models.Event, error)
	Get(string) (*models.Event, error)
	GetAll() ([]*models.Event, error)
	Update(string, *models.Event) (*models.Event, error)
	//Delete(string) error
}
