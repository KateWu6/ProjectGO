package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	InitDB()
	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the API"))
	})

	router.HandleFunc("/users/", CreateUser).Methods("POST")
	router.HandleFunc("/users/{id}/", GetUser).Methods("GET")

	router.HandleFunc("/tasks/", CreateTask).Methods("POST")
	router.HandleFunc("/tasks/user/{userId}/", GetTasksByUser).Methods("GET")
	router.HandleFunc("/tasks/{id}/", UpdateTask).Methods("PUT")
	router.HandleFunc("/tasks/{id}/", DeleteTask).Methods("DELETE")

	http.ListenAndServe(":8080", router)
}
