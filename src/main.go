package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tanvirrb/event-app-go/src/configs"
	"github.com/tanvirrb/event-app-go/src/helpers"
	"github.com/tanvirrb/event-app-go/src/router"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	port, err := helpers.GetPortFromEnv()
	if err != nil {
		log.Printf("Failed to get port: %v", err)
		os.Exit(1)
	}

	if err := configs.ConnectDB(); err != nil {
		log.Printf("Failed to connect to database: %v", err)
		os.Exit(1)
	}
	defer func() {
		err := configs.CloseDB()
		if err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	route := gin.Default()
	eventsRoutes := route.Group("/events")
	router.RegisterEventsRoutes(eventsRoutes)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: route,
	}

	serverErrors := make(chan error, 1)
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrors <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Printf("Server error: %v", err)
	case <-quit:
		log.Println("Shutdown signal received")
	}

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
