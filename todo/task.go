package todo

import (
	"encoding/json"
	"net/http"
	"github.com/gorilla/mux"
)

type Task struct {
	ID string `json: "id"`
	Order int `json: "order"`
}

var tasks []Task

// ? get all tasks

// swagger:operation GET /tasks tasks GetTasks
// ---
// summary: Retrieve all tasks.
// description: Retrieves all tasks across all lists.
// responses:
//   '200':
//     description: A list of tasks.
//     schema:
//       type: array
//       items:
//         "$ref": "#/definitions/Task"
func GetTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json");
	json.NewEncoder(w).Encode(tasks)
}

// ? get task by id 

// swagger:operation GET /task/{id} tasks GetTask
// ---
// summary: Retrieve a task by id.
// description: Retrieves a specific task by its ID.
// parameters:
// - name: id
//   in: path
//   description: ID of the task to retrieve.
//   required: true
//   type: string
// responses:
//   '200':
//     description: A single task.
//     schema:
//       "$ref": "#/definitions/Task"
//   '404':
//     description: Task not found.
func GetTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json");
	params := mux.Vars(r);
	id := params["id"]
	for _, item := range tasks {
		if id == item.ID {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	http.Error(w, "Invalid Task ID", http.StatusBadRequest);
}

// Todo: make create task method

// swagger:operation POST /task tasks CreateTask
// ---
// summary: Create a new task.
// description: Creates a new task and returns the created task.
// responses:
//   '200':
//     description: Task created successfully.
//     schema:
//       "$ref": "#/definitions/Task"
func CreateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newTask Task;
	
	_ = json.NewDecoder(r.Body).Decode(&newTask)
	tasks = append(tasks, newTask)

	json.NewEncoder(w).Encode(newTask)
}

// Todo: update a task method

// swagger:operation PUT /task/{id} tasks UpdateTask
// ---
// summary: Update an existing task.
// description: Updates an existing task and returns the updated task.
// parameters:
// - name: id
//   in: path
//   description: ID of the task to update.
//   required: true
//   type: string
// responses:
//   '200':
//     description: Task updated successfully.
//     schema:
//       "$ref": "#/definitions/Task"
func UpdateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json");
	params := mux.Vars(r);
	id := params["id"]
	for idx, item := range tasks {
		if id == item.ID {
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

// swagger:operation DELETE /task/{id} tasks DeleteTask
// ---
// summary: Delete a task.
// description: Deletes a task and returns the remaining tasks.
// parameters:
// - name: id
//   in: path
//   description: ID of the task to delete.
//   required: true
//   type: integer
// responses:
//   '200':
//     description: Task deleted successfully.
//     schema:
//       type: array
//       items:
//         "$ref": "#/definitions/Task"
func DeleteTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r);
	id := params["id"]
	for idx, item := range tasks {
		if id == item.ID {
			tasks = append(tasks[:idx], tasks[idx + 1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(tasks)
}









