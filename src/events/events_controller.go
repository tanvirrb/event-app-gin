package events

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tanvirrb/event-app-gin/src/events/interfaces"
	"github.com/tanvirrb/event-app-gin/src/events/models"
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

	createdEvent, err := c.service.Create(ctx.Request.Context(), &event)
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
	uuidStr := ctx.Param("uuid")
	uuid, err := uuid.Parse(uuidStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid uuid format",
		})
		return
	}

	event, err := c.service.Get(ctx.Request.Context(), uuid)
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
	events, err := c.service.GetAll(ctx.Request.Context())
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
	uuidStr := ctx.Param("uuid")
	uuid, err := uuid.Parse(uuidStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid uuid format",
		})
		return
	}

	var event models.Event
	if err := ctx.ShouldBindJSON(&event); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	updatedEvent, err := c.service.Update(ctx.Request.Context(), uuid, &event)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": updatedEvent,
	})
}

func (c *EventController) Delete(ctx *gin.Context) {
	uuidStr := ctx.Param("uuid")
	uuid, err := uuid.Parse(uuidStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid uuid format",
		})
		return
	}

	err = c.service.Delete(ctx.Request.Context(), uuid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Event deleted successfully",
	})
}
