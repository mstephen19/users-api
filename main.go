package main

import (
	"fmt"
	"net/http"
	"regexp"

	db "github.com/mstephen19/users-api/db"
	lib "github.com/mstephen19/users-api/lib"
)

func main() {
	router := lib.NewRouter()

	mainExp, err := regexp.Compile(`^/$`)
	router.HandleFunc(http.MethodGet, mainExp, func(writer http.ResponseWriter, req *http.Request) {
		writer.Write([]byte("Hello world"))
	})

	wildcardExp, err := regexp.Compile(`.*`)
	router.HandleFunc(http.MethodGet, wildcardExp, func(writer http.ResponseWriter, req *http.Request) {
		writer.Write([]byte("Oops, 404!"))
		return
	})

	client, err := db.ConnectToDatabaseServer()
	if err != nil {
		panic(err)
	}
	db := client.Database("users-api")

	fmt.Println(db)

	http.ListenAndServe(":3000", router)
}
