package interfaces

import (
	"github.com/tanvirrb/event-app-go/src/events/models"
)

type EventService interface {
	Create(*models.Event) (*models.Event, error)
	Get(string) (*models.Event, error)
	GetAll() ([]*models.Event, error)
	Update(string, *models.Event) (*models.Event, error)
	//Delete(string) error
}
