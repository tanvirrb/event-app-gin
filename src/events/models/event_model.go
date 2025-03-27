package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Event struct {
	_Id   primitive.ObjectID `json:"_id,omitempty"`
	Name  string             `json:"name,omitempty" validate:"required"`
	Genre string             `json:"genre,omitempty" validate:"required"`
}
