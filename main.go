package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"

	db "github.com/mstephen19/users-api/db"
	"github.com/mstephen19/users-api/db/models"
	lib "github.com/mstephen19/users-api/lib"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	// Database
	client, err := db.ConnectToDatabaseServer()
	if err != nil {
		panic(err)
	}
	db := client.Database("users-api")

	// HTTP routes
	router := lib.NewRouter()

	// On the main route, simply serve up the static HTML page
	mainExp, err := regexp.Compile(`^/$`)
	router.Handle(http.MethodGet, mainExp, http.FileServer(http.Dir("./static")))

	// Handle get, post, and delete for the users route
	usersExp, err := regexp.Compile(`^/users$`)
	// Get all users, not paginated
	router.HandleFunc(http.MethodGet, usersExp, func(writer http.ResponseWriter, req *http.Request) {
		// ? Note to self - develop a better understanding of how bson.D works
		// * Just a slice of structs with Key and Value properties
		cursor, err := db.Collection("users").Find(context.TODO(), bson.D{{}})
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Header().Set("Content-Type", "application/json")
			bytes, _ := lib.NewJsonMessage("Failed to fetch users.")
			writer.Write(bytes)
			return
		}

		var results []models.User

		for cursor.Next(context.TODO()) {
			var elem models.User
			err := cursor.Decode(&elem)
			if err != nil {
				fmt.Println(err)
			}

			results = append(results, elem)
		}

		json, err := json.Marshal(results)

		writer.WriteHeader(http.StatusAccepted)
		writer.Header().Set("Content-Type", "application/json")
		writer.Write(json)
	})
	// Add a new user
	router.HandleFunc(http.MethodPost, usersExp, func(writer http.ResponseWriter, req *http.Request) {
		bytes, err := io.ReadAll(req.Body)
		result := &models.User{}
		// If either reading the request body or parsing it fails,
		// respond with a 400 status code.
		if err != nil || json.Unmarshal(bytes, result) != nil {
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusBadRequest)
			bytes, _ := lib.NewJsonMessage("Corrupt data provided.")
			writer.Write(bytes)
			return
		}

		// Add a created at date value
		result.CreatedAt = time.Now().Format(time.UnixDate)
		result.Id = primitive.NewObjectID()
		// Insert the item into the database and retrieve the result, which just contains the objectID
		// ! Handle this error later on
		_, err = db.Collection("users").InsertOne(context.TODO(), result)
		responseBody, _ := lib.NewJsonMessage(result.Id.Hex())
		writer.WriteHeader(http.StatusAccepted)
		writer.Write(responseBody)
		writer.Header().Set("Content-Type", "application/json")
	})

	// On any other routes, redirect back to the base route
	wildcardExp, err := regexp.Compile(`.*`)
	router.HandleFunc(http.MethodGet, wildcardExp, func(writer http.ResponseWriter, req *http.Request) {
		http.Redirect(writer, req, "/", http.StatusPermanentRedirect)
		return
	})

	http.ListenAndServe(":3000", router)
}
