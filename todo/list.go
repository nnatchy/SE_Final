package todo

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sort"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Response struct {
	Status string `json "status"`
}

// ? get all lists

// swagger:operation GET /lists lists GetLists
// ---
// summary: Get all lists.
// description: Get all lists.
// responses:
//
//	'200':
//	  description: A list of lists.
//	  schema:
//	    type: array
//	    items:
//	      "$ref": "#/definitions/List"
//	'500':
//	  description: Bad Server
func GetLists(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	collection := Client.Database("test").Collection("lists")

	var lists []List
	cur, _ := collection.Find(context.Background(), bson.D{})

	for cur.Next(context.Background()) {
		var list List
		err := cur.Decode(&list)
		if err != nil {
			log.Fatal(err)
		}

		lists = append(lists, list)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	cur.Close(context.Background())

	json.NewEncoder(w).Encode(lists)
}

// ? get list by id

// swagger:operation GET /list/{id} lists GetList
// ---
// summary: Get a list by id.
// description: Get a list by id.
// parameters:
//   - name: id
//     in: path
//     description: ID of the list to retrieve.
//     required: true
//     type: string
//
// responses:
//
//	'200':
//	  description: A single list.
//	  schema:
//	    "$ref": "#/definitions/List"
//	'404':
//	  description: List not found.
//	'500':
//	  description: Bad Server
func GetList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, "Invalid List ID", http.StatusBadRequest)
		return
	}

	var list List
	var listsCollection = Client.Database("test").Collection("lists")
	err = listsCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&list)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "No Task Found", http.StatusNotFound)
			return
		}
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(list)
}

// Todo: make create list method

// swagger:operation POST /list lists CreateList
// ---
// summary: Create a new list.
// description: Creates a new list and returns the created list.
// responses:
//
//	'200':
//	  description: List created successfully.
//	  schema:
//	    "$ref": "#/definitions/List"
//	'500':
//	  description: Bad Server
func CreateList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var newList List
	_ = json.NewDecoder(r.Body).Decode(&newList)

	collection := Client.Database("test").Collection("lists")
	insertResult, err := collection.InsertOne(context.TODO(), newList)

	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(insertResult.InsertedID)
}

// Todo: update a list method

// swagger:operation PUT /list/{id} lists UpdateList
// ---
// summary: Update an existing list.
// description: Updates an existing list and returns the updated list.
// parameters:
//   - name: id
//     in: path
//     description: ID of the list to update.
//     required: true
//     type: string
//
// responses:
//
//	'200':
//	  description: List updated successfully.
//	  schema:
//	    "$ref": "#/definitions/List"
//	'404':
//	  description: List not found.
//	'500':
//	  description: Bad Server
func UpdateList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	var list List
	var listsCollection = Client.Database("test").Collection("lists")
	err = json.NewDecoder(r.Body).Decode(&list)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	updatedList := bson.D{
		{"$set", bson.D{
			{"title", list.Title},
			{"order", list.Order},
			{"tasks", list.Tasks},
		}},
	}
	res, err := listsCollection.UpdateOne(context.Background(), bson.M{"_id": id}, updatedList)
	if err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}
	if res.MatchedCount == 0 {
		http.Error(w, "No List found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(res)
}

// Todo: move a task to another list

// swagger:operation PUT /list/{listId}/task/{taskId} lists moveTask
// ---
// summary: Move a task to another list.
// description: Moves a task to another list and returns the updated list.
// parameters:
//   - name: destinationListId
//     in: path
//     description: ID of the destination list.
//     required: true
//     type: string
//   - name: taskId
//     in: path
//     description: ID of the task to move.
//     required: true
//     type: string
//
// responses:
//
//	'200':
//	  description: Task moved successfully.
//	  schema:
//	    "$ref": "#/definitions/List"
//	'404':
//	  description: List or task not found.
//	'500':
//	  description: Bad Server
func MoveTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	taskId, err := primitive.ObjectIDFromHex(params["taskId"])
	if (err != nil) {
		http.Error(w, "Invalid task Id", http.StatusBadRequest);
		return;
	}
	destinationListId, err := primitive.ObjectIDFromHex(params["destinationListId"]);
	if (err != nil) {
		http.Error(w, "Invalid list Id", http.StatusBadRequest);
		return;
	}
	// find task
	var wantedTask Task;
	var currentList *List;
	var destinationList *List;
	var lists []List;
	listsCollection := Client.Database("test").Collection("lists");

	cur, _ := listsCollection.Find(context.Background(), bson.D{})

	for cur.Next(context.Background()) {
		var list List
		err := cur.Decode(&list)
		if err != nil {
			log.Fatal(err)
		}

		lists = append(lists, list)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	tasksCollection := Client.Database("test").Collection("tasks")
	err = tasksCollection.FindOne(context.Background(), bson.M{"_id": taskId}).Decode(&wantedTask)
	if (err != nil) {
		http.Error(w, "Task not found", http.StatusNotFound);
		return;
	}
	// find task in list, destination list
	for _, list := range lists {
		if (destinationListId == list.ID) {
			destinationList = &list
		}
		for _, task := range list.Tasks {
			if (wantedTask == task) {
				currentList = &list;
			}
		}
	}
	if (currentList == nil || destinationList == nil) {
		http.Error(w, "Task not in any list or Destination List not found", http.StatusNotFound);
		return;
	}
	
	// append
	destinationList.Tasks = append(destinationList.Tasks, wantedTask);

	// remove task from old one
	for idx, task := range currentList.Tasks {
		if (task == wantedTask) {
			currentList.Tasks = append(currentList.Tasks[:idx], currentList.Tasks[idx + 1:]...)
			break;
		}
	}
	
	// assign to the mongodb

	// Start a session for transaction.
	session, err := Client.StartSession()
	if err != nil {
		http.Error(w, "Failed to start session", http.StatusInternalServerError)
		return
	}
	session.StartTransaction()

	// update
	updateDes := bson.D{{"$set", bson.D{{"tasks", destinationList.Tasks}}}}
	_, err = listsCollection.UpdateOne(context.Background(), bson.M{"_id": currentList.ID}, updateDes)
    if err != nil {
        http.Error(w, "Failed to update destination list", http.StatusInternalServerError)
        session.AbortTransaction(context.Background())
        return
    }

	updateCur := bson.D{{"$set", bson.D{{"tasks", currentList.Tasks}}}}
	_, err = listsCollection.UpdateOne(context.Background(), bson.M{"_id": currentList.ID}, updateCur)
    if err != nil {
        http.Error(w, "Failed to update current list", http.StatusInternalServerError)
        session.AbortTransaction(context.Background())
        return
    }

	// Commit the transaction.
	err = session.CommitTransaction(context.Background())
	if err != nil {
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}
	session.EndSession(context.Background())

	json.NewEncoder(w).Encode(destinationList)
}

// Todo: reorder tasks in list
func ReOrderTasks(reOrderedTasks []Task) []Task {

	sort.Slice(reOrderedTasks, func(i, j int) bool {
		return reOrderedTasks[i].Order < reOrderedTasks[j].Order
	})
	return reOrderedTasks
}

// swagger:operation PUT /list/{id}/tasks lists reOrderTasksInList
// ---
// summary: Reorder tasks in a list.
// description: Reorders tasks in a list and returns the updated list.
// parameters:
//   - name: id
//     in: path
//     description: ID of the list whose tasks need to be reordered.
//     required: true
//     type: string
//
// responses:
//
//	'200':
//	  description: Tasks reordered successfully.
//	  schema:
//	    "$ref": "#/definitions/List"
//	'404':
//	  description: List not found.
//	'500':
//	  description: Bad Server
func ReOrderTasksInList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"]);
	if (err != nil) {
		http.Error(w, "Invalid ID", http.StatusBadRequest);
		return;
	}
	var list List;
	listCollections := Client.Database("test").Collection("lists");
	err = listCollections.FindOne(context.Background(), bson.M{"_id": id}).Decode(&list);
	if (err != nil) {
		http.Error(w, "List not found", http.StatusNotFound);
		return;
	}
	var newTasks = ReOrderTasks(list.Tasks)

	// Start a session for transaction.
	session, err := Client.StartSession()
	if err != nil {
		http.Error(w, "Failed to start session", http.StatusInternalServerError)
		return
	}
	session.StartTransaction()

	// Update tasks in wantedList

	updateTasks := bson.D{{"$set", bson.D{{"tasks", newTasks}}}}
	_, err = listCollections.UpdateOne(context.Background(), bson.M{"_id": list.ID}, updateTasks)
    if err != nil {
        http.Error(w, "Failed to update current list", http.StatusInternalServerError)
        session.AbortTransaction(context.Background())
        return
    }

	// Commit the transaction.
	err = session.CommitTransaction(context.Background())
	if err != nil {
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}
	session.EndSession(context.Background())

	json.NewEncoder(w).Encode(newTasks)

}

// Todo: reorder a list

// swagger:operation PUT /list lists reOrderLists
// ---
// summary: Reorder the lists.
// description: Reorders all the lists and returns the reordered lists.
// responses:
//
//	'200':
//	  description: Lists reordered successfully.
//	  schema:
//	    type: array
//	    items:
//	      "$ref": "#/definitions/List"
//	'500':
//	  description: Bad Server
func ReOrderLists(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	listsCollection := Client.Database("test").Collection("lists")
	var newLists []List

	cur, _ := listsCollection.Find(context.Background(), bson.D{})

	for cur.Next(context.Background()) {
		var list List
		err := cur.Decode(&list)
		if err != nil {
			log.Fatal(err)
		}

		newLists = append(newLists, list)
	}

	sort.Slice(newLists, func(i, j int) bool {
		return newLists[i].Order < newLists[j].Order
	})

	// Start a session for transaction.
	session, err := Client.StartSession()
	if err != nil {
		http.Error(w, "Failed to start session", http.StatusInternalServerError)
		return
	}
	session.StartTransaction()

	// Update each document in MongoDB.
	for _, list := range newLists {
		update := bson.D{
			{"$set", bson.D{
				{"title", list.Title},
				{"order", list.Order},
				{"tasks", list.Tasks},
			}},
		}
		_, err := listsCollection.UpdateOne(context.Background(), bson.M{"_id": list.ID}, update)
		if err != nil {
			http.Error(w, "Failed to update list", http.StatusInternalServerError)
			return
		}
	}

	// Commit the transaction.
	err = session.CommitTransaction(context.Background())
	if err != nil {
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}
	session.EndSession(context.Background())

	json.NewEncoder(w).Encode(newLists)
}

// Todo: delete a list method + all tasks in it

// swagger:operation DELETE /list/{id} lists DeleteList
// ---
// summary: Delete a list.
// description: Deletes a list and returns the remaining lists.
// parameters:
//   - name: id
//     in: path
//     description: ID of the list to delete.
//     required: true
//     type: string
//
// responses:
//
//	'200':
//	  description: List deleted successfully.
//	  schema:
//	    type: array
//	    items:
//	      "$ref": "#/definitions/List"
//	'404':
//	  description: List not found.
//	'500':
//	  description: Bad Server
func DeleteList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"])

	if err != nil {
		http.Error(w, "Invalid List ID", http.StatusBadRequest)
		return
	}

	listsCollection := Client.Database("test").Collection("lists")
	tasksCollection := Client.Database("test").Collection("tasks")

	// Delete tasks within the list first
	var list List
	err = listsCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&list)
	if err != nil {
		http.Error(w, "List not found", http.StatusNotFound)
		return
	}
	for _, task := range list.Tasks {
		_, err := tasksCollection.DeleteOne(context.Background(), bson.M{"_id": primitive.ObjectID(task.ID)})
		if err != nil {
			http.Error(w, "Server Error", http.StatusInternalServerError)
			return
		}
	}

	// Delete the list
	_, err = listsCollection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}

	var res = Response{
		Status: "Deleted successfully",
	}

	json.NewEncoder(w).Encode(res)

}
