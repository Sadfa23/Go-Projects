package controller

import (
	"context"
	"crudProject-1/models"
	"crudProject-1/utils"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Specifying the structure of the request; interface in Typescript
type CreateTaskRequest struct {
	Title string `json:"title"`
	Description string `json:"description"`
	Due_date *time.Time `json:"due_date"`
	Created_at *time.Time `json:"created_at"`
	Completed bool `json:"completed"`
}

type createTaskResponse struct {
	ID string `json:"id"`
}

func CreateTask (w http.ResponseWriter, r *http.Request) {
	// set response header
	w.Header().Set("Content-Type", "application/json")

	var req  CreateTaskRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields() // helpfull to catch unknown fields
	if err := dec.Decode(&req); err != nil {
		http.Error(w, "Invalid request payload: " + err.Error(), http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Title == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}

	// Prepare Task Model, set defaults & timestamps
	now := time.Now().UTC()
	task := models.Task{
		ID : primitive.NewObjectID(),
		Title: req.Title,
		Description: req.Description,
		Completed: false, // default
		Due_date: *req.Due_date,
		Created_at: now,
	}

	// Insert into db with contect & timeout
	col := db.TasksCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	res, err := col.InsertOne(ctx, task)
	if err != nil {
		http.Error(w, "Failed to create task" +err.Error(), http.StatusInternalServerError)
		return 
	}

	// 7 return created response 201 with JSON body
	w.WriteHeader(http.StatusCreated)

	// res.InsertedID is usually primitive objectId
	idHex := ""
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		idHex = oid.Hex()
	} else {
		// fallback to fmt.sprint
		idHex = primitive.NewObjectID().Hex()
	}
	json.NewEncoder(w).Encode(createTaskResponse{ID:idHex})
}

func GetAllTasks (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	col := db.TasksCollection()
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second); defer cancel()
	cur, err := col.Find(ctx, bson.D{{}}) // MongoDB driver, please run this query — but if my context expires, stop immediately."
	if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
	defer cur.Close(ctx)
	var tasks []models.Task
	for cur.Next(ctx) {
		var t models.Task
		if err := cur.Decode(&t); err !=nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tasks = append(tasks, t)
	}
	if err := cur.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(tasks)
}
/*
Explanation of the GetAllTasks Controller above
- first the data format to be used is set, then the Tasks.Collection is retreived
- a context is created; this is sort of communication network allowing for DB and http (the r *http.Request)interactions
- and hence it has method as Timeout, cancel etc for managing the network
- in this example, network is created, with a 5 second window for requests , after a request or error the network is cancelled
- col.Find(ctx, bson.D{{}}) creates a cursor (pointer) when called;
- this is because mongo driver does not give you the list of documents from your filter directly;
- it gives you a cursor so that you loop through your collection one by one hence allowing efficiency in that
- incase the collection is huge, it does not have to fetch all the documents at once
- The ctx is included in the col.Find(ctx, bson.D{{}}) because the network still has to be managed e.g check
- if the DB request is taking too long which would eventually make the application slow etc
- in the loop, if there is still a document for the cursor to move to next, 
  - initialize a variable t models.Task; this is a Go Struct (template/ model/ interface in TS) 
  - cur.Decode(&t) reads the current document and decodes/reads it against the template above 
  - converts from MongoDB bson format to Go struct based on the t model.Tasks structure
  - if it doesn't match (decoding fails),  it returns an error
- if successful append the document to the tasks slice
*/

func GetSingleTaskById (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Fetching the id from the url
	idStr := mux.Vars(r)["id"] // the requestbody has many key value pairs and the one we want has the key of "id"
	oid, err := primitive.ObjectIDFromHex(idStr) // this creates an ObjectId (which is the one used by Mongo) from the id String
	if err != nil {
        http.Error(w, "invalid id", http.StatusBadRequest)
        return
    }
	col := db.TasksCollection()
	ctx, cancel := context.WithTimeout(r.Context(), 5 * time.Second)
	defer cancel()

	var task models.Task // declaring the Go Struct to decode against

	// find a document where _id matches the oid, then decode/read it against the Go Struct task
	// if an error occurs assign it to err variable -> then the error handler after the ; will execute
	if err := col.FindOne(ctx, bson.M{"_id": oid}).Decode(&task); err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "not found", http.StatusNotFound)
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(task)
}


func UpdateTask (w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	oid, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		http.Error(w, "invalid id for mongodb Objectd format", http.StatusBadRequest )
	}

	var patch map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	col := db.TasksCollection()
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	update := bson.M{"$set":patch}

	res, err := col.UpdateByID(ctx, oid, update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if res.MatchedCount == 0 {
		http.Error(w, "Not found", http.StatusNotFound)
	}
	w.WriteHeader(http.StatusNoContent)
}

func DeleteTask (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := mux.Vars(r)["id"]
	oid, err := primitive.ObjectIDFromHex(idStr)
	if err!=nil {
		panic(err)
	}
	col := db.TasksCollection()
	ctx, cancel := context.WithTimeout(r.Context(), 5 * time.Second)
	defer cancel()
	res, err := col.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if res.DeletedCount == 0 {
		http.Error(w, "not found", http.StatusNotFound)
        return
	}
	/* Equivalent of the above error handler
	if err = mongo.ErrNoDocuments {
		http.Error(w, "not found", http.StatusNotFound)
        return
	}
	*/
	
}

