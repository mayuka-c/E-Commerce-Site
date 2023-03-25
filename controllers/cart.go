package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gin-gonic/gin"

	"github.com/mayuka-c/e-commerce-site/database"
	"github.com/mayuka-c/e-commerce-site/models"
)

type Application struct {
	productCollection *mongo.Collection
	userCollection    *mongo.Collection
}

func NewApplication(productCollection, userCollection *mongo.Collection) *Application {
	return &Application{
		productCollection: productCollection,
		userCollection:    userCollection,
	}
}

func (app *Application) AddToCart() gin.HandlerFunc {

	return func(c *gin.Context) {
		productQueryID := c.Query("id")
		if productQueryID == "" {
			log.Fatalln("Product ID is empty")
			c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return
		}

		userQueryID := c.Query("userID")
		if userQueryID == "" {
			log.Fatalln("User ID is empty")
			c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
			return
		}

		productID, err := primitive.ObjectIDFromHex(productQueryID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, errors.New("internal server error"))
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		err = database.AddProductToCart(ctx, app.productCollection, app.userCollection, productID, userQueryID)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "adding product to cart failed!")
			return
		}

		c.IndentedJSON(http.StatusOK, "Successfully added to the cart")
	}
}

func (app *Application) RemoveItemFromCart() gin.HandlerFunc {

	return func(c *gin.Context) {
		productQueryID := c.Query("id")
		if productQueryID == "" {
			log.Fatalln("Product ID is empty")
			c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return
		}

		userQueryID := c.Query("userID")
		if userQueryID == "" {
			log.Fatalln("User ID is empty")
			c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
			return
		}

		productID, err := primitive.ObjectIDFromHex(productQueryID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, errors.New("internal server error"))
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		err = database.RemoveCartItem(ctx, app.productCollection, app.userCollection, productID, userQueryID)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "removing product from cart failed!")
			return
		}

		c.IndentedJSON(http.StatusOK, "Successfully removed from the cart")
	}
}

func (app *Application) GetItemFromCart() gin.HandlerFunc {

	return func(c *gin.Context) {
		userQueryID := c.Query("userID")
		if userQueryID == "" {
			log.Println("User ID is empty")
			c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
			return
		}

		usert_id, err := primitive.ObjectIDFromHex(userQueryID)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "Internal Server error")
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var filledCart models.User
		err = UserCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: usert_id}}).Decode(&filledCart)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, "Internal Server error")
			return
		}

		filter_match := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: usert_id}}}}
		unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$usercart"}}}}
		grouping := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"}, {Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value: "$usercart.price"}}}}}}

		pointCursor, err := UserCollection.Aggregate(ctx, mongo.Pipeline{filter_match, unwind, grouping})
		if err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, "Internal Server error")
			return
		}

		var listing []bson.M
		err = pointCursor.All(ctx, &listing)
		if err != nil {
			log.Fatalln(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		for _, json := range listing {
			c.IndentedJSON(http.StatusOK, json["total"])
			c.IndentedJSON(http.StatusOK, filledCart.UserCart)
		}
	}
}

func (app *Application) BuyFromCart() gin.HandlerFunc {

	return func(c *gin.Context) {
		userQueryID := c.Query("userID")
		if userQueryID == "" {
			log.Println("User ID is empty")
			c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		err := database.BuyItemFromCart(ctx, app.userCollection, userQueryID)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "Buying Item from cart failed!")
			return
		}

		c.IndentedJSON(http.StatusOK, "Successfully placed the order")
	}
}

func (app *Application) InstantBuy() gin.HandlerFunc {

	return func(c *gin.Context) {
		productQueryID := c.Query("id")
		if productQueryID == "" {
			log.Fatalln("Product ID is empty")
			c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return
		}

		userQueryID := c.Query("userID")
		if userQueryID == "" {
			log.Fatalln("User ID is empty")
			c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
			return
		}

		productID, err := primitive.ObjectIDFromHex(productQueryID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, errors.New("internal server error"))
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		err = database.InstantBuyer(ctx, app.productCollection, app.userCollection, productID, userQueryID)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "Placing the order failed!")
			return
		}

		c.IndentedJSON(http.StatusOK, "Successfully placed the order")
	}
}
