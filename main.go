package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	initDB()
	router := mux.NewRouter()

	router.HandleFunc("/users", createUser).Methods("POST")
	router.HandleFunc("/users/{id}", getUser).Methods("GET")

	router.HandleFunc("/tasks", createTask).Methods("POST")
	router.HandleFunc("/tasks/user/{userId}", getTasksByUser).Methods("GET")
	router.HandleFunc("/tasks/{id}", updateTask).Methods("PUT")
	router.HandleFunc("/tasks/{id}", deleteTask).Methods("DELETE")

	http.ListenAndServe(":8080", router)
}
