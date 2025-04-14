package router

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tanvirrb/event-app-gin/src/configs"
	"github.com/tanvirrb/event-app-gin/src/events"
)

type Assets struct {
	Data map[string]interface{} `json:"-"`
}

func RegisterEventsRoutes(router *gin.RouterGroup) {
	collection := configs.GetCollection("events")
	repo := events.NewEventRepository(collection)
	service := events.NewEventService(repo)
	controller := events.NewEventController(service)

	router.POST("", controller.Create)
	router.GET("/:id", controller.Get)
	router.GET("", controller.GetAll)
	router.PUT("/:id", controller.Update)
	router.DELETE("/:id", controller.Delete)
	router.POST("/assets", func(ctx *gin.Context) {
		var assets Assets
		err := ctx.ShouldBindJSON(&assets.Data)
		fmt.Printf("ASSETS: %v", assets.Data)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusCreated, gin.H{
			"message": "Assets uploaded successfully",
		})
	})
}
