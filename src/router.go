package main

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
)

func CreateRouter(collections *AppCollections) *mux.Router {
	r := mux.NewRouter()

	// health check endpoint
	r.HandleFunc("/status", HealthCheckHandler)

	// todos CRUD
	r.Methods("GET").Path("/todos").HandlerFunc(FetchTodosHandler(collections))
	r.Methods("POST").Path("/todos").HandlerFunc(AddTodoHandler(collections))
	r.Methods("PATCH").Path("/todos").HandlerFunc(UpdateTodoHandler(collections))
	r.Methods("DELETE").Path("/todos").HandlerFunc(DeleteTodoHandler(collections))

	return r
}

func HealthCheckHandler(w http.ResponseWriter, _ *http.Request) {
	if err := json.NewEncoder(w).Encode(map[string]bool{"ok": true}); err != nil {
		log.Printf("Error responding to health check: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type FetchTodosResponse struct {
	Items []*Todo `json:"items"`
}

func FetchTodosHandler(collections *AppCollections) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		var cursor *mongo.Cursor

		filter := bson.D{}
		opts := &options.FindOptions{}

		cursor, err = collections.todos.Find(context.TODO(), filter, opts)
		if err != nil {
			log.Printf("Find error: %s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		CloseCursor := func() {
			if err := cursor.Close(context.TODO()); err != nil {
				log.Printf("Close cursor error: %s", err)
			}
		}

		defer CloseCursor()

		var items []*Todo
		for cursor.Next(context.TODO()) {
			var item *Todo

			err = cursor.Decode(&item)
			if err != nil {
				log.Printf("Decode BSON error: %s", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			items = append(items, item)
		}

		err = json.NewEncoder(w).Encode(FetchTodosResponse{items})
		if err != nil {
			log.Printf("Encode JSON response error: %s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

type AddTodoPayload struct {
	Title string `json:"title"`
}

type AddTodoResponse struct {
	ID string `json:"id"`
}

func AddTodoHandler(collections *AppCollections) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		var payload AddTodoPayload
		var writeResult *mongo.InsertOneResult

		err = json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			log.Printf("Error decoding body: %s", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		writeResult, err = collections.todos.InsertOne(context.TODO(), payload)
		if err != nil {
			log.Printf("Insert todo error: %s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		id := writeResult.InsertedID.(primitive.ObjectID).Hex()
		err = json.NewEncoder(w).Encode(AddTodoResponse{id})
		if err != nil {
			log.Printf("Encode JSON response error: %s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func UpdateTodoHandler(collections *AppCollections) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func DeleteTodoHandler(collections *AppCollections) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {}
}
