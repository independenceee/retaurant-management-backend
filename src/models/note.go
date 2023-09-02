package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Note struct {
	ID        primitive.ObjectID `bson: "_id"`
	Text      string             `json: "text" validate:"required, min=2, max=100"`
	Title     string             `json: "title" validate: "required, min=2, max=100"`
	CreatedAt time.Time          `json: "createdAt"`
	UpdatedAt time.Time          `json: "updatedAt"`
	NodeId    string             `json: "userId"`
}
