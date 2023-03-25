package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gin-gonic/gin"

	"github.com/mayuka-c/e-commerce-site/models"
)

func AddAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")
		if user_id == "" {
			c.JSON(http.StatusNotFound, gin.H{"error": "user id not passed in query"})
			c.Abort()
			return
		}
		address, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, "Internal Server Error")
			return
		}
		var addresses models.Address
		addresses.Address_ID = primitive.NewObjectID()
		if err = c.BindJSON(&addresses); err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusNotAcceptable, err.Error())
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		match_filter := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: address}}}}
		unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$address"}}}}
		group := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"}, {Key: "count", Value: bson.D{primitive.E{Key: "$sum", Value: 1}}}}}}

		pointcursor, err := UserCollection.Aggregate(ctx, mongo.Pipeline{match_filter, unwind, group})
		if err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, "Internal Server Error")
			return
		}

		var addressinfo []bson.M
		if err = pointcursor.All(ctx, &addressinfo); err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var size int32
		for _, address_no := range addressinfo {
			count := address_no["count"]
			size = count.(int32)
		}
		if size < 2 {
			filter := bson.D{primitive.E{Key: "_id", Value: address}}
			update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "address", Value: addresses}}}}
			_, err := UserCollection.UpdateOne(ctx, filter, update)
			if err != nil {
				log.Println(err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		} else {
			c.IndentedJSON(http.StatusBadRequest, "Not Allowed!. only 2 addresses are allowed")
			return
		}

		c.IndentedJSON(http.StatusOK, "Successfully added the new address!")
	}
}

func EditHomeAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")
		if user_id == "" {
			c.JSON(http.StatusNotFound, gin.H{"error": "user id not passed in query"})
			return
		}

		usert_id, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, "Internal Server Error")
			return
		}

		var editaddress models.Address
		if err := c.BindJSON(&editaddress); err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusBadRequest, err.Error())
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		filter := bson.D{primitive.E{Key: "_id", Value: usert_id}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address.0.house_name", Value: editaddress.House}, {Key: "address.0.street_name", Value: editaddress.Street}, {Key: "address.0.city_name", Value: editaddress.City}, {Key: "address.0.pin_code", Value: editaddress.Pincode}}}}
		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, "Something Went Wrong")
			return
		}

		c.IndentedJSON(http.StatusOK, "Successfully Updated the Home address!")
	}
}

func EditWorkAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")
		if user_id == "" {
			c.JSON(http.StatusNotFound, gin.H{"error": "user id not passed in query"})
			return
		}

		usert_id, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, "Internal Server Error")
			return
		}

		var editaddress models.Address
		if err := c.BindJSON(&editaddress); err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusBadRequest, err.Error())
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		filter := bson.D{primitive.E{Key: "_id", Value: usert_id}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address.1.house_name", Value: editaddress.House}, {Key: "address.1.street_name", Value: editaddress.Street}, {Key: "address.1.city_name", Value: editaddress.City}, {Key: "address.1.pin_code", Value: editaddress.Pincode}}}}
		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, "something Went wrong")
			return
		}

		c.IndentedJSON(http.StatusOK, "Successfully updated the Work Address!")
	}
}

func DeleteAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")
		if user_id == "" {
			c.JSON(http.StatusNotFound, gin.H{"error": "user id not passed in query"})
			return
		}

		addresses := make([]models.Address, 0)
		usert_id, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, "Internal Server Error")
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		filter := bson.D{primitive.E{Key: "_id", Value: usert_id}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address", Value: addresses}}}}
		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, "something Went wrong")
			return
		}

		c.IndentedJSON(http.StatusOK, "Successfully Deleted the address!")
	}
}
