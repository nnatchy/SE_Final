package todo

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

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
//   '500':
//     description: Bad Server
func GetTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	collection := Client.Database("test").Collection("tasks")

	var tasks []Task
	cur, _ := collection.Find(context.Background(), bson.D{})

	for cur.Next(context.Background()) {
		var task Task
		err := cur.Decode(&task)
		if err != nil {
			log.Fatal(err)
		}

		tasks = append(tasks, task)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	cur.Close(context.Background())

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
//   '500':
//     description: Bad Server
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
//   '500':
//     description: Bad Server
func CreateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var newTask Task
	_ = json.NewDecoder(r.Body).Decode(&newTask)

	newTask.ID = uuid.New().String()

	collection := Client.Database("test").Collection("tasks")
	insertResult, err := collection.InsertOne(context.TODO(), newTask)

	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(insertResult.InsertedID)
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
//   '500':
//     description: Bad Server
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
//   type: string
// responses:
//   '200':
//     description: Task deleted successfully.
//     schema:
//       type: array
//       items:
//         "$ref": "#/definitions/Task"
//   '500':
//     description: Bad Server
func DeleteTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r);
	id := params["id"]
	for idx, item := range tasks {
		if id == item.ID {
			tasks = append(tasks[:idx], tasks[idx + 1:]...)
			json.NewEncoder(w).Encode(tasks)
			break
		}
	}
	http.Error(w, "Invalid ID", http.StatusBadRequest)
}









