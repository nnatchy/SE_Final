package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nnatchy/SE_Final/todo"
)

func main() {

	todo.Init();
	r := mux.NewRouter()

	// Tasks
	r.HandleFunc("/tasks", todo.GetTasks).Methods("GET")
	r.HandleFunc("/task/{id}", todo.GetTask).Methods("GET")
	r.HandleFunc("/task", todo.CreateTask).Methods("POST")
	r.HandleFunc("/task/{id}", todo.UpdateTask).Methods("PUT")
	r.HandleFunc("/task/{id}", todo.DeleteTask).Methods("DELETE")

	// Lists
	r.HandleFunc("/lists", todo.GetLists).Methods("GET")
	r.HandleFunc("/list/{id}", todo.GetList).Methods("GET")
	r.HandleFunc("/list", todo.CreateList).Methods("POST")
	r.HandleFunc("/list/{id}", todo.UpdateList).Methods("PUT")
	r.HandleFunc("/list/{id}", todo.DeleteList).Methods("DELETE")
	r.HandleFunc("/list/{listId}/task/{taskId}", todo.MoveTask).Methods("PUT")
	r.HandleFunc("/list/{id}/tasks", todo.ReOrderTasksInList).Methods("PUT")
	r.HandleFunc("/list", todo.ReOrderLists).Methods("PUT")

	r.PathPrefix("/swagger-ui/").Handler(http.StripPrefix("/swagger-ui/", http.FileServer(http.Dir("./swagger-ui/"))))

	// swagger API route
	fmt.Println("Swagger API: http://localhost:8080/swagger-ui/")

	// normal API routes
	fmt.Println("Normal API routes: http://localhost:8080")

	log.Fatal(http.ListenAndServe(":8080", r))
}
