package events

import (
	"github.com/gin-gonic/gin"
	"github.com/tanvirrb/event-app-gin/src/events/interfaces"
	"github.com/tanvirrb/event-app-gin/src/events/models"
	"net/http"
)

type EventController struct {
	service interfaces.EventService
}

func NewEventController(service interfaces.EventService) interfaces.EventController {
	return &EventController{
		service: service,
	}
}

func (c *EventController) Create(ctx *gin.Context) {
	var event models.Event

	if err := ctx.ShouldBindJSON(&event); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	createdEvent, err := c.service.Create(&event)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"data": createdEvent,
	})

}

func (c *EventController) Get(ctx *gin.Context) {
	id := ctx.Param("id")
	event, err := c.service.Get(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data": event,
	})
}

func (c *EventController) GetAll(ctx *gin.Context) {
	events, err := c.service.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data": events,
	})
}

func (c *EventController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var event models.Event
	if err := ctx.ShouldBindJSON(&event); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	updatedEvent, err := c.service.Update(id, &event)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": updatedEvent,
	})
}

//func (c *EventController) Delete(ctx *gin.Context) {
//	id := ctx.Param("id")
//	err := c.service.Delete(id)
//	if err != nil {
//		ctx.JSON(http.StatusInternalServerError, gin.H{
//			"error": err.Error(),
//		})
//		return
//	}
//	ctx.JSON(http.StatusOK, gin.H{
//		"message": "event deleted",
//	})
//}
