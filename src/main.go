package main

import (
	"log"
	"os"

	"github.com/tanvirrb/event-app-gin/src/bootstrap"
	"github.com/tanvirrb/event-app-gin/src/configs"
	"github.com/tanvirrb/event-app-gin/src/helpers"
)

func main() {
	port, err := helpers.GetPortFromEnv()
	if err != nil {
		log.Printf("Failed to get port: %v", err)
		return
	}

	dbType := os.Getenv("DB_TYPE")
	if dbType == "" {
		dbType = "mongodb" // default to mongodb
	}

	var db configs.Database
	switch dbType {
	case "mongodb":
		dbConfig := &configs.DatabaseConfig{
			Type:     "mongodb",
			URI:      os.Getenv("MONGODB_URI"),
			Database: os.Getenv("DB_NAME"),
		}
		if dbConfig.URI == "" {
			dbConfig.URI = "mongodb://localhost:27017"
		}
		if dbConfig.Database == "" {
			dbConfig.Database = "event-app-db"
		}
		db = configs.NewMongoDB(dbConfig)
	case "mock":
		db = configs.NewMockDB()
	default:
		log.Printf("Unsupported database type: %s", dbType)
		return
	}

	app := bootstrap.NewApp(port, db)
	defer app.Cleanup()
	app.Start()
}
