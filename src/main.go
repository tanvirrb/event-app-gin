package main

import (
	"github.com/tanvirrb/event-app-go/src/bootstrap"
	"github.com/tanvirrb/event-app-go/src/helpers"
	"log"
)

func main() {
	port, err := helpers.GetPortFromEnv()
	if err != nil {
		log.Printf("Failed to get port: %v", err)
		return
	}

	app := bootstrap.NewApp(port)
	defer app.Cleanup()
	app.Start()
}
