package todo

import (
	"context"
	"log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Task struct {
	ID string `json: "id"`
	Order int `json: "order"`
}

type List struct {
	ID    string `json: "id"`
	Title string `json: "task"`
	Order int    `json: "order"`
	Tasks []Task `json: "task"`
}

// global variables for todo package
var (
	tasks []Task
	lists []List
	Client *mongo.Client
)

func Init() {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	var err error
	Client, err = mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = Client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB!")
}
