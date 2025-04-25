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

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/tanvirrb/event-app-gin/src/bootstrap"
	"github.com/tanvirrb/event-app-gin/src/configs"
	"github.com/tanvirrb/event-app-gin/src/events/models"
)

var (
	mongodbAPIApp *bootstrap.App
)

func TestMainMongoDBAPI(m *testing.M) {
	os.Setenv("DB_TYPE", "mongodb")

	db := configs.NewMongoDB(&configs.DatabaseConfig{
		Type:     "mongodb",
		URI:      os.Getenv("MONGODB_URI"),
		Database: "event-app-test",
	})

	mongodbAPIApp = bootstrap.NewApp("3002", db)

	go func() {
		if err := mongodbAPIApp.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			os.Exit(1)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	code := m.Run()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := mongodbAPIApp.Server.Shutdown(ctx); err != nil {
		os.Exit(1)
	}

	if mongodbAPIApp != nil {
		mongodbAPIApp.Cleanup()
	}

	os.Exit(code)
}

func TestCreateEventAPI(t *testing.T) {
	event := models.Event{
		Name:  "Test Event API",
		Genre: "Test Genre API",
	}
	body, _ := json.Marshal(event)
	req := httptest.NewRequest("POST", "/events", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mongodbAPIApp.Server.Handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var createdEvent models.Event
	err := json.Unmarshal(w.Body.Bytes(), &createdEvent)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, createdEvent.Uuid)
	assert.Equal(t, event.Name, createdEvent.Name)
	assert.Equal(t, event.Genre, createdEvent.Genre)
}

func TestGetEventAPI(t *testing.T) {
	event := models.Event{
		Name:  "Test Event API",
		Genre: "Test Genre API",
	}
	body, _ := json.Marshal(event)
	req := httptest.NewRequest("POST", "/events", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mongodbAPIApp.Server.Handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var createdEvent models.Event
	err := json.Unmarshal(w.Body.Bytes(), &createdEvent)
	assert.NoError(t, err)

	getReq := httptest.NewRequest("GET", "/events/"+createdEvent.Uuid.String(), nil)
	getW := httptest.NewRecorder()
	mongodbAPIApp.Server.Handler.ServeHTTP(getW, getReq)
	assert.Equal(t, http.StatusOK, getW.Code)

	var fetchedEvent models.Event
	err = json.Unmarshal(getW.Body.Bytes(), &fetchedEvent)
	assert.NoError(t, err)
	assert.Equal(t, createdEvent.Uuid, fetchedEvent.Uuid)
	assert.Equal(t, createdEvent.Name, fetchedEvent.Name)
	assert.Equal(t, createdEvent.Genre, fetchedEvent.Genre)
}

func TestGetAllEventsAPI(t *testing.T) {
	event := models.Event{
		Name:  "Test Event API",
		Genre: "Test Genre API",
	}
	body, _ := json.Marshal(event)
	req := httptest.NewRequest("POST", "/events", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mongodbAPIApp.Server.Handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	getAllReq := httptest.NewRequest("GET", "/events", nil)
	getAllW := httptest.NewRecorder()
	mongodbAPIApp.Server.Handler.ServeHTTP(getAllW, getAllReq)
	assert.Equal(t, http.StatusOK, getAllW.Code)

	var events []models.Event
	err := json.Unmarshal(getAllW.Body.Bytes(), &events)
	assert.NoError(t, err)
	assert.NotEmpty(t, events)
}

func TestUpdateEventAPI(t *testing.T) {
	event := models.Event{
		Name:  "Test Event API",
		Genre: "Test Genre API",
	}
	body, _ := json.Marshal(event)
	req := httptest.NewRequest("POST", "/events", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mongodbAPIApp.Server.Handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var createdEvent models.Event
	err := json.Unmarshal(w.Body.Bytes(), &createdEvent)
	assert.NoError(t, err)

	updateEvent := models.Event{
		Name:  "Updated Event API",
		Genre: "Updated Genre API",
	}
	updateBody, _ := json.Marshal(updateEvent)
	updateReq := httptest.NewRequest("PUT", "/events/"+createdEvent.Uuid.String(), bytes.NewBuffer(updateBody))
	updateReq.Header.Set("Content-Type", "application/json")
	updateW := httptest.NewRecorder()
	mongodbAPIApp.Server.Handler.ServeHTTP(updateW, updateReq)
	assert.Equal(t, http.StatusOK, updateW.Code)

	var updatedEvent models.Event
	err = json.Unmarshal(updateW.Body.Bytes(), &updatedEvent)
	assert.NoError(t, err)
	assert.Equal(t, createdEvent.Uuid, updatedEvent.Uuid)
	assert.Equal(t, updateEvent.Name, updatedEvent.Name)
	assert.Equal(t, updateEvent.Genre, updatedEvent.Genre)
}

func TestDeleteEventAPI(t *testing.T) {
	event := models.Event{
		Name:  "Test Event API",
		Genre: "Test Genre API",
	}
	body, _ := json.Marshal(event)
	req := httptest.NewRequest("POST", "/events", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mongodbAPIApp.Server.Handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var createdEvent models.Event
	err := json.Unmarshal(w.Body.Bytes(), &createdEvent)
	assert.NoError(t, err)

	deleteReq := httptest.NewRequest("DELETE", "/events/"+createdEvent.Uuid.String(), nil)
	deleteW := httptest.NewRecorder()
	mongodbAPIApp.Server.Handler.ServeHTTP(deleteW, deleteReq)
	assert.Equal(t, http.StatusOK, deleteW.Code)

	getReq := httptest.NewRequest("GET", "/events/"+createdEvent.Uuid.String(), nil)
	getW := httptest.NewRecorder()
	mongodbAPIApp.Server.Handler.ServeHTTP(getW, getReq)
	assert.Equal(t, http.StatusNotFound, getW.Code)
}

func TestInvalidUUIDAPI(t *testing.T) {
	getReq := httptest.NewRequest("GET", "/events/invalid-uuid", nil)
	getW := httptest.NewRecorder()
	mongodbAPIApp.Server.Handler.ServeHTTP(getW, getReq)
	assert.Equal(t, http.StatusBadRequest, getW.Code)

	updateEvent := models.Event{
		Name:  "Test Event API",
		Genre: "Test Genre API",
	}
	updateBody, _ := json.Marshal(updateEvent)
	updateReq := httptest.NewRequest("PUT", "/events/invalid-uuid", bytes.NewBuffer(updateBody))
	updateReq.Header.Set("Content-Type", "application/json")
	updateW := httptest.NewRecorder()
	mongodbAPIApp.Server.Handler.ServeHTTP(updateW, updateReq)
	assert.Equal(t, http.StatusBadRequest, updateW.Code)

	deleteReq := httptest.NewRequest("DELETE", "/events/invalid-uuid", nil)
	deleteW := httptest.NewRecorder()
	mongodbAPIApp.Server.Handler.ServeHTTP(deleteW, deleteReq)
	assert.Equal(t, http.StatusBadRequest, deleteW.Code)
}
