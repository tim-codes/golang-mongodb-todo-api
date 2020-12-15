package main

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func main() {
	a := App{}
	a.ConnectMongo()
	a.SetupCollections()
	a.StartServer()
}

type Todo struct {
	Due       time.Time          `json:"due" bson:"due"`
	Created   time.Time          `json:"created" bson:"created"`
	Updated   time.Time          `json:"updated" bson:"updated"`
	Title     string             `json:"title" bson:"title"`
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Completed bool               `json:"completed" bson:"completed"`
}
