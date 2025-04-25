package configs

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type Database interface {
	Connect(ctx context.Context) error

	Disconnect(ctx context.Context) error

	GetCollection(name string) *mongo.Collection

	GetDatabase() *mongo.Database
}

type DatabaseConfig struct {
	Type     string
	URI      string
	Database string
	Options  map[string]interface{}
}
