package main

import (
	func_go "Project/func"
	"Project/handler"
	"Project/sessions"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", sessions.СheckSession(handler.Home_page))

	router.HandleFunc("/register", handler.RegisterHandler)
	router.HandleFunc("/login", handler.LoginHandler)

	router.HandleFunc("/users/", func_go.CreateUser).Methods("POST")
	router.HandleFunc("/users/{id}/", func_go.GetUser).Methods("GET")

	router.HandleFunc("/tasks/", func_go.CreateTask).Methods("POST")
	router.HandleFunc("/tasks/user/{userId}/", func_go.GetTasksByUser).Methods("GET")
	router.HandleFunc("/tasks/{id}/", func_go.UpdateTask).Methods("PUT")
	router.HandleFunc("/tasks/{id}/", func_go.DeleteTask).Methods("DELETE")

	http.ListenAndServe(":8080", router)
}
