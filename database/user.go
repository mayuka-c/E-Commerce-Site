package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/mayuka-c/e-commerce/models"
)

func (d *DBClient) SearchProducts(ctx context.Context) ([]models.Product, error) {

	var productList []models.Product

	cursor, err := d.productCollection.Find(ctx, bson.D{{}})
	if err != nil {
		return productList, err
	}

	defer func() {
		err := cursor.Close(ctx)
		if err != nil {
			panic(err)
		}
	}()

	err = cursor.All(ctx, &productList)
	if err != nil {
		return productList, err
	}

	if curserErr := cursor.Err(); curserErr != nil {
		return productList, err
	}

	return productList, err
}

func (d *DBClient) SearchProductsByQuery(ctx context.Context, productname string) ([]models.Product, error) {

	var productList []models.Product

	cursor, err := d.productCollection.Find(ctx, bson.M{"product_name": bson.M{"$regex": productname}})
	if err != nil {
		return productList, err
	}

	defer func() {
		err := cursor.Close(ctx)
		if err != nil {
			panic(err)
		}
	}()

	err = cursor.All(ctx, &productList)
	if err != nil {
		return productList, err
	}

	if curserErr := cursor.Err(); curserErr != nil {
		return productList, err
	}

	return productList, err
}
