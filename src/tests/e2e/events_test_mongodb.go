package e2e

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/tanvirrb/event-app-gin/src/bootstrap"
	"github.com/tanvirrb/event-app-gin/src/configs"
	"github.com/tanvirrb/event-app-gin/src/events"
	"github.com/tanvirrb/event-app-gin/src/events/models"
)

var (
	mongodbApp *bootstrap.App
)

func TestMainMongoDB(m *testing.M) {
	os.Setenv("DB_TYPE", "mongodb")

	db := configs.NewMongoDB(&configs.DatabaseConfig{
		Type:     "mongodb",
		URI:      os.Getenv("MONGODB_URI"),
		Database: "event-app-test",
	})

	mongodbApp = bootstrap.NewApp("3002", db)

	go func() {
		if err := mongodbApp.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			os.Exit(1)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	code := m.Run()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := mongodbApp.Server.Shutdown(ctx); err != nil {
		os.Exit(1)
	}

	if mongodbApp != nil {
		mongodbApp.Cleanup()
	}

	os.Exit(code)
}

func getMongoDBEventService() *events.EventService {
	collection := mongodbApp.DB.GetCollection("events")
	repo := events.NewMongoDBEventRepository(collection)
	return events.NewEventService(repo).(*events.EventService)
}

func TestCreateEventMongoDB(t *testing.T) {
	event := &models.Event{
		Name:  "Test Event MongoDB",
		Genre: "Test Genre MongoDB",
	}

	ctx := context.Background()
	createdEvent, err := getMongoDBEventService().Create(ctx, event)
	assert.NoError(t, err)
	assert.NotNil(t, createdEvent)
	assert.NotEqual(t, uuid.Nil, createdEvent.Uuid)
	assert.Equal(t, event.Name, createdEvent.Name)
	assert.Equal(t, event.Genre, createdEvent.Genre)
}

func TestGetEventMongoDB(t *testing.T) {
	ctx := context.Background()
	event := &models.Event{
		Name:  "Test Event MongoDB",
		Genre: "Test Genre MongoDB",
	}

	createdEvent, err := getMongoDBEventService().Create(ctx, event)
	assert.NoError(t, err)

	fetchedEvent, err := getMongoDBEventService().Get(ctx, createdEvent.Uuid)
	assert.NoError(t, err)
	assert.NotNil(t, fetchedEvent)
	assert.Equal(t, createdEvent.Uuid, fetchedEvent.Uuid)
	assert.Equal(t, createdEvent.Name, fetchedEvent.Name)
	assert.Equal(t, createdEvent.Genre, fetchedEvent.Genre)
}

func TestGetAllEventsMongoDB(t *testing.T) {
	ctx := context.Background()
	event := &models.Event{
		Name:  "Test Event MongoDB",
		Genre: "Test Genre MongoDB",
	}

	createdEvent, err := getMongoDBEventService().Create(ctx, event)
	assert.NoError(t, err)

	events, err := getMongoDBEventService().GetAll(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, events)
	assert.Len(t, events, 1)
	assert.Equal(t, createdEvent.Uuid, events[0].Uuid)
	assert.Equal(t, createdEvent.Name, events[0].Name)
	assert.Equal(t, createdEvent.Genre, events[0].Genre)
}

func TestUpdateEventMongoDB(t *testing.T) {
	ctx := context.Background()
	event := &models.Event{
		Name:  "Test Event MongoDB",
		Genre: "Test Genre MongoDB",
	}

	createdEvent, err := getMongoDBEventService().Create(ctx, event)
	assert.NoError(t, err)

	updatedEvent := &models.Event{
		Name:  "Updated Event MongoDB",
		Genre: "Updated Genre MongoDB",
	}

	result, err := getMongoDBEventService().Update(ctx, createdEvent.Uuid, updatedEvent)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, createdEvent.Uuid, result.Uuid)
	assert.Equal(t, updatedEvent.Name, result.Name)
	assert.Equal(t, updatedEvent.Genre, result.Genre)
}

func TestDeleteEventMongoDB(t *testing.T) {
	ctx := context.Background()
	event := &models.Event{
		Name:  "Test Event MongoDB",
		Genre: "Test Genre MongoDB",
	}

	createdEvent, err := getMongoDBEventService().Create(ctx, event)
	assert.NoError(t, err)

	err = getMongoDBEventService().Delete(ctx, createdEvent.Uuid)
	assert.NoError(t, err)

	// Verify event is deleted
	_, err = getMongoDBEventService().Get(ctx, createdEvent.Uuid)
	assert.Error(t, err)
}
