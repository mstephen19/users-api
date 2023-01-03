package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	db "github.com/mstephen19/users-api/db"
	"github.com/mstephen19/users-api/db/models"
	lib "github.com/mstephen19/users-api/lib"
	"go.mongodb.org/mongo-driver/bson"
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
	router.HandleFunc(http.MethodGet, usersExp, func(writer http.ResponseWriter, req *http.Request) {
		// ! Note to self - develop a better understanding of how bson.D works
		cursor, err := db.Collection("users").Find(context.TODO(), bson.D{{}})
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			bytes, _ := lib.NewJsonError("Failed to fetch users.")
			writer.Write(bytes)
			return
		}

		var results []models.User

		for cursor.Next(context.TODO()) {
			var elem models.User
			// ! After better understanding bson.D, develop a better understanding
			// ! of what the .Decode function does.
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

	// On any other routes, redirect back to the base route
	wildcardExp, err := regexp.Compile(`.*`)
	router.HandleFunc(http.MethodGet, wildcardExp, func(writer http.ResponseWriter, req *http.Request) {
		http.Redirect(writer, req, "/", http.StatusPermanentRedirect)
		return
	})

	http.ListenAndServe(":3000", router)
}
