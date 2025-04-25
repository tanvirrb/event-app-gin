package configs

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type MockDB struct {
	collections map[string]*mongo.Collection
	database    *mongo.Database
}

func NewMockDB() *MockDB {
	return &MockDB{
		collections: make(map[string]*mongo.Collection),
	}
}

func (m *MockDB) Connect(ctx context.Context) error {
	return nil
}

func (m *MockDB) Disconnect(ctx context.Context) error {
	return nil
}

func (m *MockDB) GetCollection(name string) *mongo.Collection {
	return m.collections[name]
}

func (m *MockDB) GetDatabase() *mongo.Database {
	return m.database
}

func (m *MockDB) SetCollection(name string, collection *mongo.Collection) {
	m.collections[name] = collection
}
