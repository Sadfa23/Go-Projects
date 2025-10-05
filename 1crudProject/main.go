package main

import (
	"crudProject-1/routes"
	"crudProject-1/utils"
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Server is starting")

	// ConnectDb
	err := db.ConnectDb("mongodb://localhost:27017")
	if err != nil {
		log.Fatal("❌ Failed to connect to MongoDB:", err)
	}

	r := router.Router()
	// Start the HTTP server
	log.Fatal(http.ListenAndServe(":3333", r))
	fmt.Println("Server is running on port 3333")
}
