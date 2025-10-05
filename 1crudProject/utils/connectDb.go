package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
This is the explanation on how ConnectDb works in Go projects
Imports
- Context is a package used to handle deadlines,  cancelletions and request-scoped values
- Example when talking to an external service like mongodb, you set a  deadline of 10 seconds
  afterwhich you cancel the process so that your application does not hang when mongodb is down or slow

- mongo is the main driver package fro mongodb in Go
- It provides you with Client, Database, Collections and CRUD methods such as InsertOne, Find, UpdateOne etc

- options helps you configure how the driver connects
- Illustration: -> options.Client().ApplyURI(uri)
-		options.Client() -> creates a blank client options object
-       .ApplyURI applies the connection string into the blank object created above
*/
var Client *mongo.Client
/*
- Once you connect, any part of your code can reuse Client instead of reconnecting.
- Why pointer? Because mongo.Client is a struct with lots of internals; we want to reference one
*/
func ConnectDb(uri string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10 *time.Second) // creates the context for db connection
	defer cancel() //cleans up resources when the function ends
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri)) // actual connection, catches the error in err variable

	// Error handler
	if err != nil {
		return  err
	}

	// Ping tests the connection btwn the server and the database
	if err := client.Ping(ctx, nil); err != nil { // This is a new err -> denoted by ":=" sign 
		return err
	}
	
	Client = client // saves the connected client  instance to your global Client
	return  nil
}

func TasksCollection () *mongo.Collection{
	return Client.Database("taskmanager_go").Collection("Task")
}
// This is a shortcut function: instead of repeating Client.Database("db_name").Collection("Collection_name"
// you just call Tasks.Collection
