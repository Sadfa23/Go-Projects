package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*The concept of JSON, BSON and STRUCTS
- Client sends JSON in a http request body
- Your Go server decodes that JSON into Go Struct using the json tags
- When you use the methods like InsertOne,
- the mongo driver converts the Go Structs to BSON using the bson tags and stores it in mongodb
- When you use methods like FindOne/Find ,
- the driver decodes BSON back to Go Structs and you encode JSON in the HTTP resonse
*/

type Task struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`
	Title string `json:"title" bson:"title"`
	Description string `json:"description" bson:"description"`
	Due_date time.Time `json:"due_date" bson:"due_date"`
	Completed bool `json:"completed" bson:"completed"`
	Created_at time.Time `json:"created_at" bson:"created_at"`
}

// Setting of Default Values is to be done after v1 of this project is completed
/*
// Constructor with defaults
func NewUser(name, email string) User {
    return User{
        Name:      name,
        Age:       18,                // default age
        IsActive:  true,              // default status
        CreatedAt: time.Now(),        // default timestamp
    }
}
*/