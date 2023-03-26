package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/mayuka-c/e-commerce-site/models"
)

var Validate = validator.New()

func hashPassword(password string) string {

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panicln(err)
	}

	return string(bytes)
}

func verifyPassword(userPassword string, givenPassword string) (bool, string) {

	valid := true
	msg := ""
	err := bcrypt.CompareHashAndPassword([]byte(givenPassword), []byte(userPassword))
	if err != nil {
		msg = "Password is invalid"
		valid = false
	}

	return valid, msg
}

func (app *Application) SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		err := c.BindJSON(&user)
		if err != nil {
			log.Fatalln(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := Validate.Struct(user)
		if validationErr != nil {
			log.Fatalln(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		count, err := app.dbClient.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			log.Fatalln(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if count > 0 {
			log.Fatalln("User already exist with provided email")
			c.JSON(http.StatusBadRequest, gin.H{"error": "user already exist with provided email"})
			return
		}

		count, err = app.dbClient.CountDocuments(ctx, bson.M{"phone": user.Phone})
		if err != nil {
			log.Fatalln(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if count > 0 {
			log.Fatalln("Provided phone number is already in use")
			c.JSON(http.StatusBadRequest, gin.H{"error": "this phone number is already in use"})
			return
		}

		password := hashPassword(*user.Password)
		user.Password = &password

		user.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		user.ID = primitive.NewObjectID()
		userIDHex := new(string)
		*userIDHex = user.ID.Hex()
		user.User_ID = userIDHex

		token, refreshToken, _ := app.tokenClient.TokenGenerator(*user.Email, *user.First_Name, *user.Last_Name, *user.User_ID)
		user.Token = &token
		user.Refresh_Token = &refreshToken
		user.UserCart = make([]models.ProductUser, 0)
		user.Address_Details = make([]models.Address, 0)
		user.Order_Status = make([]models.Order, 0)

		err = app.dbClient.InsertOne(ctx, user)
		if err != nil {
			log.Fatalln(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "the user did not get created"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"msg": "Successfully signed in!"})
	}
}

func (app *Application) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		err := c.BindJSON(&user)
		if err != nil {
			log.Fatalln(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var founduser models.User
		err = app.dbClient.FindOne(ctx, bson.M{"email": user.Email}).Decode(&founduser)
		if err != nil {
			log.Fatalln(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "login is incorrect"})
			return
		}

		PasswordIsValid, msg := verifyPassword(*user.Password, *founduser.Password)
		if !PasswordIsValid {
			log.Fatalln(msg)
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}

		token, refreshToken, _ := app.tokenClient.TokenGenerator(*founduser.Email, *founduser.First_Name, *founduser.Last_Name, *founduser.User_ID)

		app.tokenClient.UpdateAllTokens(token, refreshToken, *founduser.User_ID)

		c.JSON(http.StatusFound, founduser)
	}
}

func (app *Application) ProductViewerAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var products models.Product

		if err := c.BindJSON(&products); err != nil {
			log.Fatalln(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		products.Product_ID = primitive.NewObjectID()
		err := app.dbClient.InsertOne(ctx, products)
		if err != nil {
			log.Fatalln(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Not Created"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"msg": "Successfully added our Product Admin!!"})
	}
}

func (app *Application) SearchProducts() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		productList, err := app.dbClient.SearchProducts(ctx)
		if err != nil {
			log.Fatalln(err)
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		c.IndentedJSON(http.StatusOK, productList)
	}
}

func (app *Application) SearchProductsByQuery() gin.HandlerFunc {
	return func(c *gin.Context) {
		productName := c.Query("name")
		if productName == "" {
			log.Println("product name is empty")
			c.JSON(http.StatusNotFound, gin.H{"error": "invalid search index"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		productList, err := app.dbClient.SearchProductsByQuery(ctx, productName)
		if err != nil {
			log.Fatalln(err)
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		c.IndentedJSON(http.StatusOK, productList)
	}
}
