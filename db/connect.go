package db

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const MongoURI = "mongodb://localhost:27017"

// Connect to the MongoDB database server.
func ConnectToDatabaseServer() (client *mongo.Client, err error) {
	client, err = mongo.Connect(
		context.TODO(),
		options.Client().ApplyURI(MongoURI),
	)
	// If the initial connection failed, return out early so the error
	// can be handled.
	if err != nil {
		return
	}
	// Attempt to ping the server.
	err = client.Ping(context.TODO(), readpref.Primary())

	return
}
