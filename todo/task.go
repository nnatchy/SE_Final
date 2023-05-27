package todo

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, "Invalid Task ID", http.StatusBadRequest)
		return
	}

	var task Task
	var tasksCollection = Client.Database("test").Collection("tasks")
	err = tasksCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "No Task Found", http.StatusNotFound)
			return
		}
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(task)
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
	id, err := primitive.ObjectIDFromHex(params["id"]);
	if (err != nil) {
		http.Error(w, "Invalid ID", http.StatusBadRequest);
		return;
	}
	var task Task;
	var tasksCollection = Client.Database("test").Collection("tasks")
	err = json.NewDecoder(r.Body).Decode(&task);
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest);
		return;
	}
	updatedTask := bson.D {
		{"$set", bson.D{
			{"order", task.Order},
		}},
	}
	res, err := tasksCollection.UpdateOne(context.Background(), bson.M{"_id": id}, updatedTask);
	if (err != nil) {
		http.Error(w, "Server Error", http.StatusInternalServerError);
		return;
	}
	if (res.MatchedCount == 0) {
		http.Error(w, "No task found", http.StatusNotFound);
		return;
	}
	json.NewEncoder(w).Encode(res);
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
	params := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, "Invalid Task ID", http.StatusBadRequest)
		return
	}
	var tasksCollection = Client.Database("test").Collection("tasks")
	res, err := tasksCollection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}

	if res.DeletedCount == 0 {
		http.Error(w, "No task found", http.StatusNotFound);
		return;
	}

	json.NewEncoder(w).Encode(res)
}









