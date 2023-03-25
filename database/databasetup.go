package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/mayuka-c/e-commerce-site/config"
)

var dbconfig config.DBConfig
var Client *mongo.Client = DBSet()

func DBSet() *mongo.Client {

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://" + dbconfig.DB_URL))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Successfully connected to mongoDB\n")
	return client
}

func UserData(client *mongo.Client, collectionName string) *mongo.Collection {

	collection := client.Database("Ecommerce").Collection(collectionName)
	return collection
}

func ProductData(client *mongo.Client, collectionName string) *mongo.Collection {

	collection := client.Database("Ecommerce").Collection(collectionName)
	return collection
}
