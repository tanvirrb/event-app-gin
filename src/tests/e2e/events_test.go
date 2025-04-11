package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
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
	if os.Getenv("MONGODB_URI") == "" {
		err := os.Setenv("MONGODB_URI", "mongodb://localhost:27017")
		if err != nil {
			os.Exit(1)
		}
	}

	configs.SetDBName("event-app-test-db")
	err := configs.ConnectDB()
	if err != nil {
		os.Exit(1)
	}

	app = bootstrap.NewApp("3002")

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

func setupTestDB(t *testing.T) {
	testDB := configs.DB.Database("event-app-test-db")
	err := testDB.Drop(context.Background())
	assert.NoError(t, err, "Failed to drop test database")
}

func TestCreateEvent(t *testing.T) {
	setupTestDB(t)

	event := models.Event{
		Name:  "Test Event",
		Genre: "Test Genre",
	}

	jsonData, err := json.Marshal(event)
	assert.NoError(t, err)

	req := httptest.NewRequest("POST", "/events", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	app.Server.Handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response struct {
		Data models.Event `json:"data"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.NotEmpty(t, response.Data.Id)
	assert.Equal(t, event.Name, response.Data.Name)
	assert.Equal(t, event.Genre, response.Data.Genre)
}

func TestGetEvent(t *testing.T) {
	setupTestDB(t)

	collection := configs.GetCollection("events")
	repo := events.NewEventRepository(collection)
	eventService := events.NewEventService(repo)
	event := &models.Event{
		Name:  "Test Event",
		Genre: "Test Genre",
	}

	createdEvent, err := eventService.Create(event)
	assert.NoError(t, err)
	assert.NotEmpty(t, createdEvent.Id)

	req := httptest.NewRequest("GET", "/events/"+createdEvent.Id.Hex(), nil)
	w := httptest.NewRecorder()
	app.Server.Handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var getEventResponse struct {
		Data models.Event `json:"data"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &getEventResponse)
	assert.NoError(t, err)

	assert.Equal(t, event.Name, getEventResponse.Data.Name)
	assert.Equal(t, event.Genre, getEventResponse.Data.Genre)
}

func TestGetAllEvents(t *testing.T) {
	setupTestDB(t)

	collection := configs.GetCollection("events")
	repo := events.NewEventRepository(collection)
	eventService := events.NewEventService(repo)

	eventList := []*models.Event{
		{Name: "Event 1", Genre: "Genre 1"},
		{Name: "Event 2", Genre: "Genre 2"},
	}

	for _, event := range eventList {
		_, err := eventService.Create(event)
		assert.NoError(t, err)
	}

	req := httptest.NewRequest("GET", "/events", nil)
	w := httptest.NewRecorder()
	app.Server.Handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Data []models.Event `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.GreaterOrEqual(t, len(response.Data), len(eventList))
}

func TestUpdateEvent(t *testing.T) {
	setupTestDB(t)

	collection := configs.GetCollection("events")
	repo := events.NewEventRepository(collection)
	eventService := events.NewEventService(repo)

	event := &models.Event{
		Name:  "Test Event",
		Genre: "Test Genre",
	}

	createdEvent, err := eventService.Create(event)
	assert.NoError(t, err)
	assert.NotEmpty(t, createdEvent.Id)

	updatedEvent := &models.Event{
		Name:  "Updated Event",
		Genre: "Updated Genre",
	}

	jsonData, err := json.Marshal(updatedEvent)
	assert.NoError(t, err)

	req := httptest.NewRequest("PUT", "/events/"+createdEvent.Id.Hex(), bytes.NewBuffer(jsonData))

	w := httptest.NewRecorder()
	app.Server.Handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Data models.Event `json:"data"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, updatedEvent.Name, response.Data.Name)
	assert.Equal(t, updatedEvent.Genre, response.Data.Genre)
}

func TestDeleteEvent(t *testing.T) {
	setupTestDB(t)

	collection := configs.GetCollection("events")
	repo := events.NewEventRepository(collection)
	eventService := events.NewEventService(repo)

	event := &models.Event{
		Name:  "Test Event",
		Genre: "Test Genre",
	}

	createdEvent, err := eventService.Create(event)
	assert.NoError(t, err)
	assert.NotEmpty(t, createdEvent.Id)

	req := httptest.NewRequest("DELETE", "/events/"+createdEvent.Id.Hex(), nil)
	w := httptest.NewRecorder()
	app.Server.Handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Data string `json:"data"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.NotEmpty(t, response.Data)
	assert.IsType(t, "", response.Data)
}
