package events

import (
	"github.com/tanvirrb/event-app-go/src/events/interfaces"
	"github.com/tanvirrb/event-app-go/src/events/models"
)

type EventService struct {
	repo interfaces.EventRepository
}

func NewEventService(repo interfaces.EventRepository) interfaces.EventService {
	return &EventService{
		repo: repo,
	}
}

func (s *EventService) Create(e *models.Event) (*models.Event, error) {
	return s.repo.Create(e)
}

func (s *EventService) Get(id string) (*models.Event, error) {
	return s.repo.Get(id)
}

//
//func (s *EventService) GetAll() ([]*models.Event, error) {
//	return s.repo.GetAll()
//}
//
//func (s *EventService) Update(id string, e *models.Event) (*models.Event, error) {
//	return s.repo.Update(id, e)
//}
//
//func (s *EventService) Delete(id string) error {
//	return s.repo.Delete(id)
//}
