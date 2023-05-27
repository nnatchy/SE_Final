package todo

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"sort"
	"github.com/google/uuid"
)

// ? get all lists

// swagger:operation GET /lists lists GetLists
// ---
// summary: Get all lists.
// description: Get all lists.
// responses:
//   '200':
//     description: A list of lists.
//     schema:
//       type: array
//       items:
//         "$ref": "#/definitions/List"
//   '500':
//     description: Bad Server
func GetLists(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lists)
}

// ? get list by id

// swagger:operation GET /list/{id} lists GetList
// ---
// summary: Get a list by id.
// description: Get a list by id.
// parameters:
// - name: id
//   in: path
//   description: ID of the list to retrieve.
//   required: true
//   type: string
// responses:
//   '200':
//     description: A single list.
//     schema:
//       "$ref": "#/definitions/List"
//   '404':
//     description: List not found.
//   '500':
//     description: Bad Server
func GetList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id := params["id"]

	for _, item := range lists {
		if id == item.ID {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	http.Error(w, "Invalid ID", http.StatusBadRequest);
}

// Todo: make create list method

// swagger:operation POST /list lists CreateList
// ---
// summary: Create a new list.
// description: Creates a new list and returns the created list.
// responses:
//   '200':
//     description: List created successfully.
//     schema:
//       "$ref": "#/definitions/List"
//   '500':
//     description: Bad Server
func CreateList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var newList List
	_ = json.NewDecoder(r.Body).Decode(&newList)

	newList.ID = uuid.New().String()

	for i := range newList.Tasks {
		newList.Tasks[i].ID = uuid.New().String()
		tasks = append(tasks, newList.Tasks[i])  // append the tasks to global tasks
	}

	lists = append(lists, newList)

	json.NewEncoder(w).Encode(newList)
}

// Todo: update a list method

// swagger:operation PUT /list/{id} lists UpdateList
// ---
// summary: Update an existing list.
// description: Updates an existing list and returns the updated list.
// parameters:
// - name: id
//   in: path
//   description: ID of the list to update.
//   required: true
//   type: string
// responses:
//   '200':
//     description: List updated successfully.
//     schema:
//       "$ref": "#/definitions/List"
//   '404':
//     description: List not found.
//   '500':
//     description: Bad Server
func UpdateList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id := params["id"]
	for idx, item := range lists {
		if id == item.ID {
			lists = append(lists[:idx], lists[idx+1:]...) // remove wantedList
			var newList List
			_ = json.NewDecoder(r.Body).Decode(&newList)
			newList.ID = item.ID
			lists = append(lists, newList)
			json.NewEncoder(w).Encode(newList)
			return
		}
	}
	http.Error(w, "List not found !", http.StatusNotFound);
}

// Todo: move a task to another list

// swagger:operation PUT /list/{listId}/task/{taskId} lists moveTask
// ---
// summary: Move a task to another list.
// description: Moves a task to another list and returns the updated list.
// parameters:
// - name: destinationListId
//   in: path
//   description: ID of the destination list.
//   required: true
//   type: string
// - name: taskId
//   in: path
//   description: ID of the task to move.
//   required: true
//   type: string
// responses:
//   '200':
//     description: Task moved successfully.
//     schema:
//       "$ref": "#/definitions/List"
//   '404':
//     description: List or task not found.
//   '500':
//     description: Bad Server
func MoveTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	var wantTask *Task
	var curList *List

	// find wanted task
	taskId := params["id"]
	for _, l := range lists {
		for idx, task := range l.Tasks {
			if taskId == task.ID {
				wantTask = &l.Tasks[idx]
				curList = &l
				break
			}
		}
	}

	if wantTask == nil || curList == nil {
		http.Error(w, "Task Not Found !", http.StatusNotFound)
		return
	}

	var newList List
	_ = json.NewDecoder(r.Body).Decode(&newList)

	// find destination list
	var destinationList *List
	for _, l := range lists {
		if l.ID == newList.ID {
			destinationList = &l
			break
		}
	}

	if destinationList == nil {
		http.Error(w, "List not found !", http.StatusNotFound)
		return
	}

	// append the task to destination list
	destinationList.Tasks = append(destinationList.Tasks, *wantTask)

	// remove wantedTask in old list
	for idx, task := range curList.Tasks {
		if task.ID == wantTask.ID {
			curList.Tasks = append(curList.Tasks[:idx], curList.Tasks[idx+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(destinationList.Tasks)
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
// - name: id
//   in: path
//   description: ID of the list whose tasks need to be reordered.
//   required: true
//   type: string
// responses:
//   '200':
//     description: Tasks reordered successfully.
//     schema:
//       "$ref": "#/definitions/List"
//   '404':
//     description: List not found.
//   '500':
//     description: Bad Server
func ReOrderTasksInList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id := params["id"]
	for _, list := range lists {
		if id == list.ID { // find wantedList
			var reOrderedTasks []Task
			reOrderedTasks = ReOrderTasks(list.Tasks)
			list.Tasks = reOrderedTasks
			json.NewEncoder(w).Encode(list)
			break
		}
	}
	http.Error(w, "List not Found !", http.StatusNotFound)
}

// Todo: reorder a list

// swagger:operation PUT /list lists reOrderLists
// ---
// summary: Reorder the lists.
// description: Reorders all the lists and returns the reordered lists.
// responses:
//   '200':
//     description: Lists reordered successfully.
//     schema:
//       type: array
//       items:
//         "$ref": "#/definitions/List"
//   '500':
//     description: Bad Server
func ReOrderLists(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	sort.Slice(lists, func(i, j int) bool {
		return lists[i].Order < lists[j].Order
	})
	json.NewEncoder(w).Encode(lists)
}

// Todo: delete a list method + all tasks in it


// swagger:operation DELETE /list/{id} lists DeleteList
// ---
// summary: Delete a list.
// description: Deletes a list and returns the remaining lists.
// parameters:
// - name: id
//   in: path
//   description: ID of the list to delete.
//   required: true
//   type: string
// responses:
//   '200':
//     description: List deleted successfully.
//     schema:
//       type: array
//       items:
//         "$ref": "#/definitions/List"
//   '404':
//     description: List not found.
//   '500':
//     description: Bad Server
func DeleteList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id := params["id"]
	for idx, list := range lists {
		if id == list.ID {
			// remove every tasks in it
			for idxTask, _ := range list.Tasks {
				list.Tasks = append(list.Tasks[:idxTask], list.Tasks[idxTask + 1:]...)
			}
			lists = append(lists[:idx], lists[idx+1:]...)
			json.NewEncoder(w).Encode(lists)
			return;
		}
	}
	http.Error(w, "Invalid List ID", http.StatusBadRequest);
}
