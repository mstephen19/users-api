package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id        primitive.ObjectID `json:"_id", bson:"_id"`
	Name      string             `json:"name" bson:"name"`
	CreatedAt string             `json:"created_at" bson:"created_at"`
}
