package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (d *DBClient) CountDocuments(ctx context.Context, filter interface{}) (int64, error) {

	count, err := d.userCollection.CountDocuments(ctx, filter)
	if err != nil {
		return count, err
	}

	return count, err
}

func (d *DBClient) InsertOne(ctx context.Context, document interface{}) error {

	_, inserterr := d.userCollection.InsertOne(ctx, document)
	if inserterr != nil {
		return inserterr
	}

	return nil
}

func (d *DBClient) FindOne(ctx context.Context, filter interface{}) *mongo.SingleResult {

	result := d.userCollection.FindOne(ctx, filter)
	return result
}

func (d *DBClient) Find(ctx context.Context, filter interface{}) *mongo.SingleResult {

	result := d.userCollection.FindOne(ctx, filter)
	return result
}

func (d *DBClient) UpdateOne(ctx context.Context, filter interface{}, setValue interface{}, opt options.UpdateOptions) error {

	_, err := d.userCollection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: setValue}}, &opt)
	if err != nil {
		return err
	}

	return nil
}
