package controllers

import (
	"context"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (app *Application) AddToCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryID := c.Query("id")
		if productQueryID == "" {
			log.Error("Product ID is empty")
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "product id is empty"})
			return
		}

		userQueryID := c.Query("userID")
		if userQueryID == "" {
			log.Error("User ID is empty")
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "user id is empty"})
			return
		}

		product_id, err := primitive.ObjectIDFromHex(productQueryID)
		if err != nil {
			log.Error(err)
			c.IndentedJSON(http.StatusExpectationFailed, gin.H{"error": "productID provided is invalid"})
			return
		}

		user_id, err := primitive.ObjectIDFromHex(userQueryID)
		if err != nil {
			log.Error(err)
			c.IndentedJSON(http.StatusExpectationFailed, gin.H{"error": "userID provided is invalid"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		err = app.dbClient.AddProductToCart(ctx, product_id, user_id)
		if err != nil {
			log.Error(err)
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		c.IndentedJSON(http.StatusOK, gin.H{"msg": "Successfully added to the cart"})
	}
}

func (app *Application) RemoveItemFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryID := c.Query("id")
		if productQueryID == "" {
			log.Error("Product ID is empty")
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "product id is empty"})
			return
		}

		userQueryID := c.Query("userID")
		if userQueryID == "" {
			log.Error("User ID is empty")
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "user id is empty"})
			return
		}

		product_id, err := primitive.ObjectIDFromHex(productQueryID)
		if err != nil {
			log.Error(err)
			c.IndentedJSON(http.StatusExpectationFailed, gin.H{"error": "productID provided is invalid"})
			return
		}

		user_id, err := primitive.ObjectIDFromHex(userQueryID)
		if err != nil {
			log.Error(err)
			c.IndentedJSON(http.StatusExpectationFailed, gin.H{"error": "userID provided is invalid"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		err = app.dbClient.RemoveCartItem(ctx, product_id, user_id)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		c.IndentedJSON(http.StatusOK, gin.H{"msg": "Successfully removed from the cart"})
	}
}

func (app *Application) GetItemFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		userQueryID := c.Query("userID")
		if userQueryID == "" {
			log.Error("User ID is empty")
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "user id is empty"})
			return
		}

		user_id, err := primitive.ObjectIDFromHex(userQueryID)
		if err != nil {
			log.Error(err)
			c.IndentedJSON(http.StatusExpectationFailed, gin.H{"error": "userID provided is invalid"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		result, totalPrice, err := app.dbClient.GetItemFromCart(ctx, user_id)
		if err != nil {
			log.Error(err)
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		c.IndentedJSON(http.StatusOK, gin.H{"result": result.UserCart, "totalPrice": totalPrice})
	}
}

func (app *Application) BuyFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		userQueryID := c.Query("userID")
		if userQueryID == "" {
			log.Error("User ID is empty")
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "user id is empty"})
			return
		}

		user_id, err := primitive.ObjectIDFromHex(userQueryID)
		if err != nil {
			log.Error(err)
			c.IndentedJSON(http.StatusExpectationFailed, gin.H{"error": "userID provided is invalid"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		err = app.dbClient.BuyItemFromCart(ctx, user_id)
		if err != nil {
			log.Error(err)
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		c.IndentedJSON(http.StatusOK, gin.H{"msg": "Successfully placed the order!"})
	}
}

func (app *Application) InstantBuy() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryID := c.Query("id")
		if productQueryID == "" {
			log.Error("Product ID is empty")
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "product id is empty"})
			return
		}

		userQueryID := c.Query("userID")
		if userQueryID == "" {
			log.Error("User ID is empty")
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "user id is empty"})
			return
		}

		product_id, err := primitive.ObjectIDFromHex(productQueryID)
		if err != nil {
			log.Error(err)
			c.IndentedJSON(http.StatusExpectationFailed, gin.H{"error": "productID provided is invalid"})
			return
		}

		user_id, err := primitive.ObjectIDFromHex(userQueryID)
		if err != nil {
			log.Error(err)
			c.IndentedJSON(http.StatusExpectationFailed, gin.H{"error": "userID provided is invalid"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		err = app.dbClient.InstantBuyer(ctx, product_id, user_id)
		if err != nil {
			log.Error(err)
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		c.IndentedJSON(http.StatusOK, gin.H{"msg": "Successfully placed the order"})
	}
}
