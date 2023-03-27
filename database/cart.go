package database

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/mayuka-c/e-commerce-site/models"
)

var (
	ErrCantFindProduct    = errors.New("can't find the product")
	ErrCantDecodeProducts = errors.New("can't find the product")
	ErrCantUpdateUser     = errors.New("cannot add this product to the cart")
	ErrCantRemoveItem     = errors.New("cannot remove this product from the cart")
	ErrCantGetItem        = errors.New("unable to get the item from the cart")
	ErrCantBuyCartItem    = errors.New("cannot update the purchase")
	ErrCantDoInstantBuyer = errors.New("cannot update the purchase")
)

func (d *DBClient) AddProductToCart(ctx context.Context, product_id, user_id primitive.ObjectID) error {

	searchfromdb, err := d.productCollection.Find(ctx, bson.M{"_id": product_id})
	if err != nil {
		return ErrCantFindProduct
	}

	var productcart []models.ProductUser
	err = searchfromdb.All(ctx, &productcart)
	if err != nil {
		return ErrCantDecodeProducts
	}

	filter := bson.D{{Key: "_id", Value: user_id}}
	update := bson.D{{Key: "$push", Value: bson.D{{Key: "usercart", Value: bson.D{{Key: "$each", Value: productcart}}}}}}

	_, err = d.userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return ErrCantUpdateUser
	}

	return nil
}

func (d *DBClient) RemoveCartItem(ctx context.Context, product_id, user_id primitive.ObjectID) error {

	filter := bson.D{{Key: "_id", Value: user_id}}
	update := bson.M{"$pull": bson.M{"usercart": bson.M{"_id": product_id}}}

	_, err := d.userCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		return ErrCantRemoveItem
	}

	return nil
}

func (d *DBClient) GetItemFromCart(ctx context.Context, user_id primitive.ObjectID) (models.User, string, error) {

	var filledCart models.User
	var totalPrice string

	err := d.userCollection.FindOne(ctx, bson.D{{Key: "_id", Value: user_id}}).Decode(&filledCart)
	if err != nil {
		return filledCart, totalPrice, err
	}

	filter_match := bson.D{{Key: "$match", Value: bson.D{{Key: "_id", Value: user_id}}}}
	unwind := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$usercart"}}}}
	grouping := bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: "$_id"}, {Key: "total", Value: bson.D{{Key: "$sum", Value: "$usercart.price"}}}}}}

	pointCursor, err := d.userCollection.Aggregate(ctx, mongo.Pipeline{filter_match, unwind, grouping})
	if err != nil {
		return filledCart, totalPrice, err
	}

	var listing []bson.M
	err = pointCursor.All(ctx, &listing)
	if err != nil {
		return filledCart, totalPrice, err
	}

	for _, json := range listing {
		totalPrice = json["total"].(string)
	}

	return filledCart, totalPrice, nil
}

func (d *DBClient) BuyItemFromCart(ctx context.Context, user_id primitive.ObjectID) error {

	var getCartItems models.User
	var orderCart models.Order

	orderCart.Order_ID = primitive.NewObjectID()
	orderCart.Ordered_At = time.Now()
	orderCart.Order_Cart = make([]models.ProductUser, 0)
	orderCart.Payment_Method.CashOnDelivery = true

	unwind := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$usercart"}}}}
	grouping := bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: "$_id"}, {Key: "total_price", Value: bson.D{{Key: "$sum", Value: "$usercart.price"}}}}}}

	totalPriceCursor, err := d.userCollection.Aggregate(ctx, mongo.Pipeline{unwind, grouping})
	if err != nil {
		return ErrCantBuyCartItem
	}

	var totalPriceDoc []bson.M

	if err = totalPriceCursor.All(ctx, &totalPriceDoc); err != nil {
		return ErrCantBuyCartItem
	}

	var total_price int32

	for _, user_item := range totalPriceDoc {
		total_price = user_item["total_price"].(int32)
	}

	orderCart.Price = int(total_price)

	filter := bson.D{{Key: "_id", Value: user_id}}
	update := bson.D{{Key: "$push", Value: bson.D{{Key: "orders", Value: orderCart}}}}

	_, err = d.userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return ErrCantBuyCartItem
	}

	err = d.userCollection.FindOne(ctx, bson.D{{Key: "_id", Value: user_id}}).Decode(&getCartItems)
	if err != nil {
		return ErrCantBuyCartItem
	}

	filter2 := bson.D{{Key: "_id", Value: user_id}}
	update2 := bson.M{"$push": bson.M{"orders.order_list": bson.M{"$each": getCartItems.UserCart}}}

	_, err = d.userCollection.UpdateOne(ctx, filter2, update2)
	if err != nil {
		return ErrCantBuyCartItem
	}

	usercart_empty := make([]models.ProductUser, 0)
	filtered := bson.D{{Key: "_id", Value: user_id}}
	updated := bson.D{{Key: "$set", Value: bson.D{{Key: "usercart", Value: usercart_empty}}}}

	_, err = d.userCollection.UpdateOne(ctx, filtered, updated)
	if err != nil {
		return ErrCantBuyCartItem
	}

	return nil
}

func (d *DBClient) InstantBuyer(ctx context.Context, product_id, user_id primitive.ObjectID) error {

	var product_details models.ProductUser
	var orders_detail models.Order

	orders_detail.Order_ID = primitive.NewObjectID()
	orders_detail.Ordered_At = time.Now()
	orders_detail.Order_Cart = make([]models.ProductUser, 0)
	orders_detail.Payment_Method.CashOnDelivery = true

	err := d.productCollection.FindOne(ctx, bson.D{{Key: "_id", Value: product_id}}).Decode(&product_details)
	if err != nil {
		return ErrCantDoInstantBuyer
	}

	orders_detail.Price = product_details.Price
	filter := bson.D{{Key: "_id", Value: user_id}}
	update := bson.D{{Key: "$push", Value: bson.D{{Key: "orders", Value: orders_detail}}}}

	_, err = d.userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return ErrCantDoInstantBuyer
	}

	filter2 := bson.D{{Key: "_id", Value: user_id}}
	update2 := bson.M{"$push": bson.M{"orders.order_list": product_details}}

	_, err = d.userCollection.UpdateOne(ctx, filter2, update2)
	if err != nil {
		return ErrCantDoInstantBuyer
	}

	return nil
}
