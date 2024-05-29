// In database.go

package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var UsersCollection *mongo.Collection
var OrdersCollection *mongo.Collection
var ProductsCollection *mongo.Collection
var PointsCollection *mongo.Collection

func InitDatabase(uri string) error {
	// Set client options
	clientOptions := options.Client().ApplyURI(uri)

	// Connect to MongoDB
	var err error
	client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return err
	}

	// Set a timeout for the context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Ping the primary
	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	log.Println("Connected to MongoDB!")
	return nil
}

func InitCollection() {
	// Use a consistent case for the database name
	const databaseName = "cafeshop"
	UsersCollection = client.Database(databaseName).Collection("Users")
	OrdersCollection = client.Database(databaseName).Collection("Orders")
	ProductsCollection = client.Database(databaseName).Collection("Products")
	PointsCollection = client.Database(databaseName).Collection("Points")
}

func Disconnect() {
	if client != nil {
		client.Disconnect(context.TODO())
		log.Println("Disconnected from MongoDB!")
	}
}
