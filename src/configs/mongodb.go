package configs

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	DB     *mongo.Client
	dbName = "event-app-db"
)

func SetDBName(name string) {
	dbName = name
}

func ConnectDB() error {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017" // fallback for local development
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	// Verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	DB = client
	fmt.Println("Connected to MongoDB!")
	return nil
}

func GetCollection(collection string) *mongo.Collection {
	return DB.Database(dbName).Collection(collection)
}

func CloseDB() error {
	if DB != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return DB.Disconnect(ctx)
	}
	return nil
}
