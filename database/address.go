package database

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/mayuka-c/e-commerce/models"
)

var (
	ErrAddAddress = errors.New("more than 2 addresses are already present, so cannot add new one")
)

func (d *DBClient) AddAddress(ctx context.Context, user_id primitive.ObjectID, address models.Address) error {

	match_filter := bson.D{{Key: "$match", Value: bson.D{{Key: "_id", Value: user_id}}}}
	unwind := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$address"}}}}
	group := bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: "$_id"}, {Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}}}}}

	pointcursor, err := d.userCollection.Aggregate(ctx, mongo.Pipeline{match_filter, unwind, group})
	if err != nil {
		return err
	}

	var addressInfo []bson.M
	if err = pointcursor.All(ctx, &addressInfo); err != nil {
		return err
	}

	var size int32
	for _, address_no := range addressInfo {
		count := address_no["count"]
		size = count.(int32)
	}

	if size < 2 {
		filter := bson.D{{Key: "_id", Value: user_id}}
		update := bson.D{{Key: "$push", Value: bson.D{{Key: "address", Value: address}}}}

		_, err := d.userCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			return err
		}
	} else {
		return ErrAddAddress
	}

	return nil
}

func (d *DBClient) EditHomeAddress(ctx context.Context, user_id primitive.ObjectID, address models.Address) error {

	filter := bson.D{primitive.E{Key: "_id", Value: user_id}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "address.0.house_name", Value: address.House}, {Key: "address.0.street_name", Value: address.Street}, {Key: "address.0.city_name", Value: address.City}, {Key: "address.0.pin_code", Value: address.Pincode}}}}

	_, err := d.userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (d *DBClient) EditWorkAddress(ctx context.Context, user_id primitive.ObjectID, address models.Address) error {

	filter := bson.D{primitive.E{Key: "_id", Value: user_id}}
	update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address.1.house_name", Value: address.House}, {Key: "address.1.street_name", Value: address.Street}, {Key: "address.1.city_name", Value: address.City}, {Key: "address.1.pin_code", Value: address.Pincode}}}}

	_, err := d.userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (d *DBClient) DeleteAddress(ctx context.Context, user_id primitive.ObjectID) error {

	// setting to empty slice
	address := make([]models.Address, 0)

	filter := bson.D{primitive.E{Key: "_id", Value: user_id}}
	update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address", Value: address}}}}

	_, err := d.userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}
