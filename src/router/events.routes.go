package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tanvirrb/event-app-go/src/events"
)

func RegisterEventsRoutes(router *gin.RouterGroup) {
	repo := events.NewEventRepository()
	service := events.NewEventService(repo)
	controller := events.NewEventController(service)

	router.POST("/", controller.Create)
	router.GET("/:id", controller.Get)
	router.GET("/", controller.GetAll)
	router.PUT("/:id", controller.Update)
	router.DELETE("/:id", controller.Delete)
}
