package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Event struct {
	Id    primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name  string             `json:"name,omitempty" validate:"required"`
	Genre string             `json:"genre,omitempty" validate:"required"`
}
