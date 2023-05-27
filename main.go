package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nnatchy/SE_Final/todo"
)

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/tasks", todo.GetTasks).Methods("GET")
	r.HandleFunc("/tasks/{id}", todo.GetTask).Methods("GET")
	r.HandleFunc("/tasks", todo.CreateTask).Methods("POST")
	r.HandleFunc("/tasks/{id}", todo.UpdateTask).Methods("PUT")
	r.HandleFunc("/tasks/{id}", todo.DeleteTask).Methods("DELETE")

	r.PathPrefix("/swagger-ui/").Handler(http.StripPrefix("/swagger-ui/", http.FileServer(http.Dir("./swagger-ui/"))))
	
	log.Fatal(http.ListenAndServe(":8080", r))
}






