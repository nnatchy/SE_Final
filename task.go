package main

import (
	"encoding/json"
	"net/http"
)

type Task struct {
	Order string `json: "order"`
}

var tasks []Task

// ? get all tasks (why not)
func getTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json");
	json.NewEncoder(w).Encode(tasks)
}

// ? get task by id (why not)
func getTask(w http.ResponseWriter, r *http.Request) {
	
}

// TODO: make create task method

func createTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
}

// Todo: update a task method


// Todo: delete a task method

