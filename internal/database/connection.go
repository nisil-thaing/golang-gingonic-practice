package database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var dbClient *mongo.Client

func SetupDBConnection(dbUri string) {
	if dbUri == "" {
		log.Fatal("Oops! Please specify the connection string first!")
	}

	clientOptions := options.Client().ApplyURI(dbUri)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("The MongoDB connection has been successfully established!")

	dbClient = client
}

func GetDBInstance() *mongo.Client {
	return dbClient
}
