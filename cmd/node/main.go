package main

import (
	"cde/internal/handlers"
	"cde/internal/storage"
	"context"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	log.Println(os.Getenv("ADDR"))
	mux := http.NewServeMux()

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Ping(context.TODO(), nil); err != nil {
		log.Fatal(err)
	}

	db := client.Database("friends")
	mongoStore := &storage.MongoStorage{Store: db}

	service := handlers.Service{MongoStorage: mongoStore}
	mux.HandleFunc("/create", service.Create)
	mux.HandleFunc("/make_friends", service.MakeFriends)
	mux.HandleFunc("/deliete", service.DelieteUser)
	mux.HandleFunc("/friends/", service.GetFriends)
	mux.HandleFunc("/newAge/", service.NewAge)
	mux.HandleFunc("/allUsers", service.GetUsers)

	http.ListenAndServe(os.Getenv("ADDR"), mux)
}
