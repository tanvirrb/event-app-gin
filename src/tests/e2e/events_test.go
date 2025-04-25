package e2e

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tanvirrb/event-app-gin/src/bootstrap"
	"github.com/tanvirrb/event-app-gin/src/configs"
	"github.com/tanvirrb/event-app-gin/src/events"
	"github.com/tanvirrb/event-app-gin/src/events/models"
)

var (
	app *bootstrap.App
)

func TestMain(m *testing.M) {
	// Use mock database for testing
	db := configs.NewMockDB()
	app = bootstrap.NewApp("3002", db)

	go func() {
		if err := app.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			os.Exit(1)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	code := m.Run()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := app.Server.Shutdown(ctx); err != nil {
		os.Exit(1)
	}

	if app != nil {
		app.Cleanup()
	}

	os.Exit(code)
}

func TestCreateEvent(t *testing.T) {
	repo := events.NewMockEventRepository()
	eventService := events.NewEventService(repo)
	event := &models.Event{
		Name:  "Test Event",
		Genre: "Test Genre",
	}

	ctx := context.Background()
	createdEvent, err := eventService.Create(ctx, event)
	assert.NoError(t, err)
	assert.NotNil(t, createdEvent)
	assert.NotEmpty(t, createdEvent.Id)
	assert.Equal(t, event.Name, createdEvent.Name)
	assert.Equal(t, event.Genre, createdEvent.Genre)
}

func TestGetEvent(t *testing.T) {
	repo := events.NewMockEventRepository()
	eventService := events.NewEventService(repo)

	ctx := context.Background()
	event := &models.Event{
		Name:  "Test Event",
		Genre: "Test Genre",
	}

	createdEvent, err := eventService.Create(ctx, event)
	assert.NoError(t, err)

	fetchedEvent, err := eventService.Get(ctx, createdEvent.Id.Hex())
	assert.NoError(t, err)
	assert.NotNil(t, fetchedEvent)
	assert.Equal(t, createdEvent.Id, fetchedEvent.Id)
	assert.Equal(t, createdEvent.Name, fetchedEvent.Name)
	assert.Equal(t, createdEvent.Genre, fetchedEvent.Genre)
}

func TestGetAllEvents(t *testing.T) {
	repo := events.NewMockEventRepository()
	eventService := events.NewEventService(repo)

	ctx := context.Background()
	event := &models.Event{
		Name:  "Test Event",
		Genre: "Test Genre",
	}

	createdEvent, err := eventService.Create(ctx, event)
	assert.NoError(t, err)

	events, err := eventService.GetAll(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, events)
	assert.Len(t, events, 1)
	assert.Equal(t, createdEvent.Id, events[0].Id)
	assert.Equal(t, createdEvent.Name, events[0].Name)
	assert.Equal(t, createdEvent.Genre, events[0].Genre)
}

func TestUpdateEvent(t *testing.T) {
	repo := events.NewMockEventRepository()
	eventService := events.NewEventService(repo)

	ctx := context.Background()
	event := &models.Event{
		Name:  "Test Event",
		Genre: "Test Genre",
	}

	createdEvent, err := eventService.Create(ctx, event)
	assert.NoError(t, err)

	updatedEvent := &models.Event{
		Name:  "Updated Event",
		Genre: "Updated Genre",
	}

	result, err := eventService.Update(ctx, createdEvent.Id.Hex(), updatedEvent)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, createdEvent.Id, result.Id)
	assert.Equal(t, updatedEvent.Name, result.Name)
	assert.Equal(t, updatedEvent.Genre, result.Genre)
}

func TestDeleteEvent(t *testing.T) {
	repo := events.NewMockEventRepository()
	eventService := events.NewEventService(repo)

	ctx := context.Background()
	event := &models.Event{
		Name:  "Test Event",
		Genre: "Test Genre",
	}

	createdEvent, err := eventService.Create(ctx, event)
	assert.NoError(t, err)

	err = eventService.Delete(ctx, createdEvent.Id.Hex())
	assert.NoError(t, err)

	// Verify event is deleted
	_, err = eventService.Get(ctx, createdEvent.Id.Hex())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "event not found")
}
