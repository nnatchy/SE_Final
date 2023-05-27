package todo

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"sort"
	"strconv"
)

type List struct {
	ID    int    `json: "id"`
	Title string `json: "task"`
	Order int    `json: "order"`
	Tasks []Task `json: "task"`
}

var lists []List

// ? get all lists
func GetLists(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lists)
}

// ? get list by id
func GetList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	for _, item := range lists {
		if id == item.ID {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Task{})
}

// Todo: make create list method

func CreateList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newList List

	_ = json.NewDecoder(r.Body).Decode(&newList)
	lists = append(lists, newList)

	json.NewEncoder(w).Encode(newList)
}

// Todo: update a list method

func UpdateList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
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

}

// Todo: move a task to another list

func moveTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	var wantTask *Task
	var curList *List

	// find wanted task
	taskId, err := strconv.Atoi(params["id"]);
	if (err != nil) {
		http.Error(w, "Invalid Task ID", http.StatusBadRequest);
		return;
	}
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
func reOrderTasks(reOrderedTasks []Task) []Task {

	sort.Slice(reOrderedTasks, func(i, j int) bool {
		return reOrderedTasks[i].Order < reOrderedTasks[j].Order
	})
	return reOrderedTasks
}

func reOrderTasksInList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"]);
	if (err != nil) {
		http.Error(w, "Invalid ID", http.StatusBadRequest);
		return;
	}
	for _, list := range lists {
		if id == list.ID { // find wantedList
			var reOrderedTasks []Task
			reOrderedTasks = reOrderTasks(list.Tasks)
			list.Tasks = reOrderedTasks
			json.NewEncoder(w).Encode(list)
			break
		}
	}
}

// Todo: reorder a list
func reOrderLists(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	sort.Slice(lists, func(i, j int) bool {
		return lists[i].Order < lists[j].Order
	})
	json.NewEncoder(w).Encode(lists)
}

// Todo: delete a list method + all tasks in it

func DeleteList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"]); 
	if (err != nil) {
		http.Error(w, "Invalid ID", http.StatusBadRequest);
		return;
	}
	for idx, list := range lists {
		if id == list.ID {
			// remove every tasks in it
			for idxTask, _ := range list.Tasks {
				list.Tasks = append(list.Tasks[:idxTask], list.Tasks[idxTask + 1:]...)
			}
			lists = append(lists[:idx], lists[idx+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(lists)
}
