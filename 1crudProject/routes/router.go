package router

import (
	controller "crudProject-1/controllers"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/api/tasks/create", controller.CreateTask).Methods("POST")
	r.HandleFunc("/api/tasks/{id}", controller.GetSingleTaskById).Methods("GET")
	r.HandleFunc("/api/tasks", controller.GetAllTasks).Methods("GET")
	r.HandleFunc("/api/task/update/{id}", controller.UpdateTask).Methods("PUT")
	r.HandleFunc("/api/task/delete/{id}", controller.DeleteTask).Methods("DELETE")

	return r
}