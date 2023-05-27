package todo

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Task struct {
	ID    primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Order int `bson:"order,omitempty" json:"order,omitempty"`
}

type List struct {
	ID    primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title string `bson:"title,omitempty" json:"title,omitempty"`
	Order int    `bson:"order,omitempty" json:"order,omitempty"`
	Tasks []Task `bson:"tasks,omitempty" json:"tasks,omitempty"`
}


// global variables for todo package
var (
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
	// var tasksCollection = Client.Database("test").Collection("tasks")
	log.Println("Connected to MongoDB!")
}
