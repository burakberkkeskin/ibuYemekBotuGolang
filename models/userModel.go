package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID           primitive.ObjectID `json:"id,omitempty"`
	ChatID       int64              `json:"chatID,omitempty" validate:"required"`
	Username     string             `json:"username,omitempty" validate:"required"`
	Name         string             `json:"name,omitempty" validate:"required"`
	IsSubscribed bool               `json:"isSubscribed,omitempty" validate:"required"`
}
