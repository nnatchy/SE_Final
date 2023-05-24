package todo

import (
	"encoding/json"
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
)

type List struct {
	ID int `json: "id"`
	Title string `json: "task"`
	Order string `json: "order"`
	Task *Task `json: "task"`
}

var lists []List;

// ? get all lists (why not)
func GetLists(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json");
	json.NewEncoder(w).Encode(lists)
}

// ? get list by id (why not)
func GetList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json");
	params := mux.Vars(r)
	for _, item := range lists {
		item.ID, _ = strconv.Atoi(params["id"]);
		if item.ID == item.ID {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Task{});
}

// Todo: make create list method

func CreateList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newList List;
	
	_ = json.NewDecoder(r.Body).Decode(&newList)
	lists = append(lists, newList)

	json.NewEncoder(w).Encode(newList)
}

// Todo: update a list method

func UpdateList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json");
	params := mux.Vars(r);

	for idx, item := range lists {
		if item.ID, _ = strconv.Atoi(params["id"]); item.ID == item.ID {
			lists = append(lists[:idx], lists[idx + 1:]...)
			var newList List
			_ = json.NewDecoder(r.Body).Decode(&newList)
			newList.ID = item.ID
			lists = append(lists, newList)
			json.NewEncoder(w).Encode(newList)
			return;
		}
	}
}

// Todo: move a task to another list
func moveTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "applciation/json");

	
}

// Todo: reorder a task in list

// Todo: reorder a list

// Todo: delete a list method + all tasks in it

func DeleteList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r);
	for idx, item := range lists {
		if item.ID, _ = strconv.Atoi(params["id"]); item.ID == item.ID {
			lists = append(lists[:idx], lists[idx + 1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(lists)
}



