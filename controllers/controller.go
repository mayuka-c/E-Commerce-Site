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

	"github.com/mayuka-c/e-commerce-site/database"
	"github.com/mayuka-c/e-commerce-site/models"
	generate "github.com/mayuka-c/e-commerce-site/tokens"
)

var UserCollection = database.UserData(database.Client, "Users")
var ProductCollection = database.ProductData(database.Client, "Products")
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

func SignUp() gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		err := c.BindJSON(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := Validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		count, err := UserCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user already exist with provided email"})
			return
		}

		count, err = UserCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if count > 0 {
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

		token, refreshToken, _ := generate.TokenGenerator(*user.Email, *user.First_Name, *user.Last_Name, *user.User_ID)
		user.Token = &token
		user.Refresh_Token = &refreshToken
		user.UserCart = make([]models.ProductUser, 0)
		user.Address_Details = make([]models.Address, 0)
		user.Order_Status = make([]models.Order, 0)

		_, inserterr := UserCollection.InsertOne(ctx, user)
		if inserterr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "the user did not get created"})
			return
		}

		c.JSON(http.StatusCreated, "Successfully signed in!")
	}
}

func Login() gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		err := c.BindJSON(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var founduser models.User
		err = UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&founduser)
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "login is incorrect"})
			return
		}

		PasswordIsValid, msg := verifyPassword(*user.Password, *founduser.Password)
		if !PasswordIsValid {
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			log.Println(msg)
			return
		}

		token, refreshToken, _ := generate.TokenGenerator(*founduser.Email, *founduser.First_Name, *founduser.Last_Name, *founduser.User_ID)

		generate.UpdateAllTokens(token, refreshToken, *founduser.User_ID)

		c.JSON(http.StatusFound, founduser)
	}
}

func ProductViewerAdmin() gin.HandlerFunc {

	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var products models.Product

		if err := c.BindJSON(&products); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		products.Product_ID = primitive.NewObjectID()
		_, anyerr := ProductCollection.InsertOne(ctx, products)
		if anyerr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Not Created"})
			return
		}

		c.JSON(http.StatusOK, "Successfully added our Product Admin!!")
	}
}

func SearchProducts() gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var productList []models.Product

		cursor, err := ProductCollection.Find(ctx, bson.D{{}})

		defer func() {
			err := cursor.Close(ctx)
			if err != nil {
				panic(err)
			}
		}()

		if err != nil {
			log.Fatalln(err)
			c.IndentedJSON(http.StatusInternalServerError, "something went wrong, please retry after some time")
			return
		}

		err = cursor.All(ctx, &productList)
		if err != nil {
			log.Fatalln(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if curserErr := cursor.Err(); curserErr != nil {
			log.Fatalln(curserErr)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.IndentedJSON(http.StatusOK, productList)
	}
}

func SearchProductsByQuery() gin.HandlerFunc {
	return func(c *gin.Context) {
		queryParam := c.Query("name")
		if queryParam == "" {
			log.Println("query is empty")
			c.JSON(http.StatusNotFound, gin.H{"error": "invalid search index"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var productList []models.Product

		cursor, err := ProductCollection.Find(ctx, bson.M{"product_name": bson.M{"$regex": queryParam}})
		defer func() {
			err := cursor.Close(ctx)
			if err != nil {
				panic(err)
			}
		}()
		if err != nil {
			log.Fatalln(err)
			c.IndentedJSON(http.StatusInternalServerError, "something went wrong, please retry after some time")
			return
		}

		err = cursor.All(ctx, &productList)
		if err != nil {
			log.Fatalln(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if curserErr := cursor.Err(); curserErr != nil {
			log.Fatalln(curserErr)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.IndentedJSON(http.StatusOK, productList)
	}
}
