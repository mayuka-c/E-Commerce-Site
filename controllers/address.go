package controllers

import (
	"context"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/mayuka-c/e-commerce-site/database"
	"github.com/mayuka-c/e-commerce-site/models"
)

func (app *Application) AddAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Query("id")
		if userId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user id not passed in query"})
			c.Abort()
			return
		}

		var address models.Address
		address.Address_ID = primitive.NewObjectID()

		if err := c.BindJSON(&address); err != nil {
			log.Error(err)
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user_id, err := primitive.ObjectIDFromHex(userId)
		if err != nil {
			log.Error(err)
			c.IndentedJSON(http.StatusExpectationFailed, gin.H{"error": "userID provided is invalid"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		err = app.dbClient.AddAddress(ctx, user_id, address)
		if err != nil {
			log.Error(err)
			if err == database.ErrAddAddress {
				c.IndentedJSON(http.StatusNotAcceptable, gin.H{"msg": database.ErrAddAddress.Error()})
			} else {
				c.AbortWithStatus(http.StatusInternalServerError)
			}
			return
		}

		c.IndentedJSON(http.StatusOK, gin.H{"msg": "Successfully added the new address!"})
	}
}

func (app *Application) EditHomeAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Query("id")
		if userId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user id not passed in query"})
			return
		}

		var editaddress models.Address
		if err := c.BindJSON(&editaddress); err != nil {
			log.Error(err)
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user_id, err := primitive.ObjectIDFromHex(userId)
		if err != nil {
			log.Error(err)
			c.IndentedJSON(http.StatusExpectationFailed, gin.H{"error": "userID provided is invalid"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		app.dbClient.EditHomeAddress(ctx, user_id, editaddress)
		if err != nil {
			log.Error(err)
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		c.IndentedJSON(http.StatusOK, gin.H{"msg": "Successfully Updated the Home address!"})
	}
}

func (app *Application) EditWorkAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Query("id")
		if userId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user id not passed in query"})
			return
		}

		var editaddress models.Address
		if err := c.BindJSON(&editaddress); err != nil {
			log.Error(err)
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user_id, err := primitive.ObjectIDFromHex(userId)
		if err != nil {
			log.Error(err)
			c.IndentedJSON(http.StatusExpectationFailed, gin.H{"error": "userID provided is invalid"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		app.dbClient.EditWorkAddress(ctx, user_id, editaddress)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		c.IndentedJSON(http.StatusOK, gin.H{"msg": "Successfully updated the Work Address!"})
	}
}

func (app *Application) DeleteAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Query("id")
		if userId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user id not passed in query"})
			return
		}

		user_id, err := primitive.ObjectIDFromHex(userId)
		if err != nil {
			log.Error(err)
			c.IndentedJSON(http.StatusExpectationFailed, gin.H{"error": "userID provided is invalid"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		app.dbClient.DeleteAddress(ctx, user_id)
		if err != nil {
			log.Error(err)
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		c.IndentedJSON(http.StatusOK, gin.H{"msg": "Successfully Deleted the address!"})
	}
}
