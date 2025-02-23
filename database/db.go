package database

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func StartDB() *mongo.Client {
	err := godotenv.Load(".env")

	if err != nil {
		log.Println("Info: .env file not found, will use environment variables")
	}

	MongoURI := os.Getenv("MONGOURI")

	if MongoURI == "" {
		log.Fatal("MONGODB_URL environment variable is not set")
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(MongoURI))

	if err != nil {

		log.Fatal("Error while connecting to Mongo: ", err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	err = client.Connect(ctx)

	if err != nil {
		log.Fatal("Error while connecting to Mongo: ", err.Error())
	}

	log.Println("Connected to MongoDB Successfully")
	return client
}

var Client *mongo.Client = StartDB()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	MONGO_DB_NAME := os.Getenv("MONGO_DATABASE_NAME")

	var collection *mongo.Collection = client.Database(MONGO_DB_NAME).Collection(collectionName)
	return collection
}
