package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tanvirrb/event-app-gin/src/configs"
	"github.com/tanvirrb/event-app-gin/src/events"
	"github.com/tanvirrb/event-app-gin/src/events/interfaces"
)

func RegisterEventsRoutes(router *gin.RouterGroup, db configs.Database) {
	collection := db.GetCollection("events")

	// Create repository based on database type
	var repo interfaces.EventRepository
	switch db.(type) {
	case *configs.MongoDB:
		repo = events.NewMongoDBEventRepository(collection)
	case *configs.MockDB:
		repo = events.NewMockEventRepository()
	default:
		panic("unsupported database type")
	}

	service := events.NewEventService(repo)
	controller := events.NewEventController(service)

	router.POST("", controller.Create)
	router.GET("/:id", controller.Get)
	router.GET("", controller.GetAll)
	router.PUT("/:id", controller.Update)
	router.DELETE("/:id", controller.Delete)
}
