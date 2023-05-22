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

// TODO: make create task method

func createTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
}