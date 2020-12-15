package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client mongo.Client
var todosCollection mongo.Collection

func main() {
	fmt.Println("Starting mongo todo app...")

	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27011,localhost:27012,localhost:27013/todos?replicaSet=rs0")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	todosCollection := client.Database("test").Collection("todos")

	fmt.Println("Connected to MongoDB!")

	router := NewRouter()
	fmt.Println("Starting HTTP Server on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

type Todo struct {
	Name      string    `json:"name"`
	Completed bool      `json:"completed"`
	Due       time.Time `json:"due"`
	Created   time.Time `json:"created"`
	Updated   time.Time `json:"updated"`
}

type Todos []Todo
