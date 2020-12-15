package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
)

type AppCollections struct {
	todos *mongo.Collection
}

type App struct {
	mongo       *mongo.Client
	collections *AppCollections
	router      *mux.Router
}

func (app *App) ConnectMongo() {
	fmt.Println("Starting mongo todo app...")

	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

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

	app.mongo = client
	fmt.Println("Connected to MongoDB!")
}

func (app *App) SetupCollections() {
	// Setup collections
	app.collections = new(AppCollections)
	app.collections.todos = app.mongo.Database("reminders-test").Collection("todos")
	fmt.Println("Initialised collections: [todos]")
}

func (app *App) StartServer() {
	app.router = CreateRouter(app.collections)
	fmt.Println("Starting HTTP Server on :8080")
	log.Fatal(http.ListenAndServe(":8080", app.router))
}
