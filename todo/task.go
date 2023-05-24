package todo

import (
	"encoding/json"
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
)

type Task struct {
	ID int `json: "id"`
	Order string `json: "order"`
}

var tasks []Task

// ? get all tasks (why not)
func GetTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json");
	json.NewEncoder(w).Encode(tasks)
}

// ? get task by id (why not)
func GetTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json");
	params := mux.Vars(r)
	for _, item := range tasks {
		item.ID, _ = strconv.Atoi(params["id"]);
		if item.ID == item.ID {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Task{});
}

// Todo: make create task method

func CreateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newTask Task;
	
	_ = json.NewDecoder(r.Body).Decode(&newTask)
	tasks = append(tasks, newTask)

	json.NewEncoder(w).Encode(newTask)
}

// Todo: update a task method

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json");
	params := mux.Vars(r);

	for idx, item := range tasks {
		if item.ID, _ = strconv.Atoi(params["id"]); item.ID == item.ID {
			tasks = append(tasks[:idx], tasks[idx + 1:]...)
			var newTask Task
			_ = json.NewDecoder(r.Body).Decode(&newTask)
			newTask.ID = item.ID
			tasks = append(tasks, newTask)
			json.NewEncoder(w).Encode(newTask)
			return;
		}
	}
}

// Todo: delete a task method

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r);
	for idx, item := range tasks {
		if item.ID, _ = strconv.Atoi(params["id"]); item.ID == item.ID {
			tasks = append(tasks[:idx], tasks[idx + 1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(tasks)
}









