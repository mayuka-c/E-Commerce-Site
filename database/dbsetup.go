package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	log "github.com/sirupsen/logrus"

	"github.com/mayuka-c/e-commerce-site/config"
	"github.com/mayuka-c/e-commerce-site/constants"
)

type DBClient struct {
	client            *mongo.Client
	userCollection    *mongo.Collection
	productCollection *mongo.Collection
}

func DBSet(dbConfig config.DBConfig) *DBClient {

	mongoClient, err := mongo.NewClient(options.Client().ApplyURI("mongodb://" + dbConfig.DB_URL))
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err = mongoClient.Connect(ctx)
	if err != nil {
		panic(err)
	}

	err = mongoClient.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}

	log.Println("Successfully connected to mongoDB")

	userCollection := mongoClient.Database("Ecommerce").Collection(constants.UserCollectionName)
	productCollection := mongoClient.Database("Ecommerce").Collection(constants.ProductCollectionName)

	return &DBClient{
		client:            mongoClient,
		userCollection:    userCollection,
		productCollection: productCollection,
	}
}

func (d *DBClient) GetUserCollection() *mongo.Collection {
	return d.userCollection
}

func (d *DBClient) GetProductCollection() *mongo.Collection {
	return d.productCollection
}
