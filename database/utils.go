package database

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (d *DBClient) CountDocuments(ctx context.Context, collection *mongo.Collection, filter interface{}) (int64, error) {

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return count, err
	}

	return count, err
}

func (d *DBClient) InsertOne(ctx context.Context, collection *mongo.Collection, document interface{}) error {

	_, inserterr := collection.InsertOne(ctx, document)
	if inserterr != nil {
		return inserterr
	}

	return nil
}

func (d *DBClient) FindOne(ctx context.Context, collection *mongo.Collection, filter interface{}) *mongo.SingleResult {

	result := collection.FindOne(ctx, filter)
	return result
}

func (d *DBClient) Find(ctx context.Context, collection *mongo.Collection, filter interface{}) *mongo.SingleResult {

	result := collection.FindOne(ctx, filter)
	return result
}

func (d *DBClient) UpdateOne(ctx context.Context, collection *mongo.Collection, filter interface{}, update interface{}, opt options.UpdateOptions) error {

	_, err := collection.UpdateOne(ctx, filter, update, &opt)
	if err != nil {
		return err
	}

	return nil
}
