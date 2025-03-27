package main

import (
	"github.com/gin-gonic/gin"
	"github.com/tanvirrb/event-app-go/src/configs"
	"github.com/tanvirrb/event-app-go/src/router" // Update this import path to match your project structure
)

func main() {
	configs.ConnectDB()
	route := gin.Default()
	eventsRoutes := route.Group("/events")
	router.RegisterEventsRoutes(eventsRoutes)

	err := route.Run(":3001")
	if err != nil {
		println("Error starting server : ", err.Error())
		return
	}
}
