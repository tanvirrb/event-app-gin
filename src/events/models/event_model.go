package models

import "github.com/google/uuid"

type Event struct {
	Uuid  uuid.UUID `bson:"uuid" json:"uuid"`
	Name  string    `json:"name,omitempty" validate:"required"`
	Genre string    `json:"genre,omitempty" validate:"required"`
}
