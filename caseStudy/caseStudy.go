package caseStudy

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

)

type Book struct {
	ID int `json: "id"`
	ISBN string `json: "isbn"`
	Title string `json: "title"`
	Author *Author `json: "author"`
}

type Author struct {
	FirstName string `json: "firstname"`
	LastName string `json: "lastname"`
}

var books []Book

// get all books
func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

// get book by id
func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range books {
		item.ID, _ = strconv.Atoi(params["id"]);
		if item.ID == item.ID {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Book{});
}

// create book
func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json");
	var book Book;
	
	_ = json.NewDecoder(r.Body).Decode(&book)
	books = append(books, book)

	json.NewEncoder(w).Encode(book)
}

// Update
func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json");
	params := mux.Vars(r);

	for idx, item := range books {
		if item.ID, _ = strconv.Atoi(params["id"]); item.ID == item.ID {
			books = append(books[:idx], books[idx + 1:]...)
			var book Book
			_ = json.NewDecoder(r.Body).Decode(&book)
			book.ID = item.ID
			books = append(books, book)
			json.NewEncoder(w).Encode(book)
			return;
		}
	}
}

// Delete
func deleteBook (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r);
	for idx, item := range books {
		if item.ID, _ = strconv.Atoi(params["id"]); item.ID == item.ID {
			books = append(books[:idx], books[idx + 1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(books)
}

func main() {
	r := mux.NewRouter()

	books = append(books, Book{ID: 1, ISBN: "1221321", Title: "Book 1", Author: &Author{FirstName: "Ayyo", LastName: "That bitch"}});
	books = append(books, Book{ID: 2, ISBN: "4214142", Title: "Book 2", Author: &Author{FirstName: "Suppo", LastName: "Ma nibba"}});

	r.HandleFunc("/api/books", getBooks).Methods("GET")
	r.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	r.HandleFunc("/api/books", createBook).Methods("POST")
	r.HandleFunc("/api/books/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", r))
}