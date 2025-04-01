package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tanvirrb/event-app-go/src/bootstrap"
	"github.com/tanvirrb/event-app-go/src/configs"
	"github.com/tanvirrb/event-app-go/src/events"
	"github.com/tanvirrb/event-app-go/src/events/models"
)

var (
	app *bootstrap.App
)

func setupTestApp(t *testing.T) {
	err := os.Setenv("MONGODB_URI", "mongodb://mongodb:27017")
	if err != nil {
		return
	}

	configs.SetDBName("event-app-test-db")

	err = configs.ConnectDB()
	assert.NoError(t, err, "Failed to connect to MongoDB")

	testDB := configs.DB.Database("event-app-test-db")
	err = testDB.Drop(context.Background())
	assert.NoError(t, err, "Failed to drop test database")

	app = bootstrap.NewApp("3002")
}

func teardownTestApp() {
	if app != nil {
		app.Cleanup()
	}
}

func TestCreateEvent(t *testing.T) {
	setupTestApp(t)
	defer teardownTestApp()

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
	setupTestApp(t)
	defer teardownTestApp()

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

//func TestGetAllEvents(t *testing.T) {
//	setupTestApp(t)
//	defer app.Cleanup()
//
//	// Create multiple events
//	events := []models.Event{
//		{
//			Name:  "Event 1",
//			Genre: "Genre 1",
//		},
//		{
//			Name:  "Event 2",
//			Genre: "Genre 2",
//		},
//	}
//
//	for _, event := range events {
//		jsonData, err := json.Marshal(event)
//		assert.NoError(t, err)
//
//		req := httptest.NewRequest("POST", "/events", bytes.NewBuffer(jsonData))
//		req.Header.Set("Content-Type", "application/json")
//		w := httptest.NewRecorder()
//		app.Server.Handler.ServeHTTP(w, req)
//		assert.Equal(t, http.StatusCreated, w.Code)
//	}
//
//	// Get all events
//	req := httptest.NewRequest("GET", "/events", nil)
//	w := httptest.NewRecorder()
//	app.Server.Handler.ServeHTTP(w, req)
//
//	assert.Equal(t, http.StatusOK, w.Code)
//
//	var response struct {
//		Data []models.Event `json:"data"`
//	}
//	err := json.Unmarshal(w.Body.Bytes(), &response)
//	assert.NoError(t, err)
//
//	// Assert we have at least the events we created
//	assert.GreaterOrEqual(t, len(response.Data), len(events))
//}
//
//func TestGetNonExistentEvent(t *testing.T) {
//	setupTestApp(t)
//	defer app.Cleanup()
//
//	// Try to get a non-existent event with a valid ObjectID format
//	req := httptest.NewRequest("GET", "/events/507f1f77bcf86cd799439011", nil)
//	w := httptest.NewRecorder()
//	app.Server.Handler.ServeHTTP(w, req)
//
//	assert.Equal(t, http.StatusNotFound, w.Code)
//}
