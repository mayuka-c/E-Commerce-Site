package tokens

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	jwt "github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"

	"github.com/mayuka-c/e-commerce/database"
)

type SignedDetails struct {
	Email     string
	FirstName string
	LastName  string
	UUID      string
	jwt.StandardClaims
}

var JWT_Key = "my_secret_value"

type TokenGenrator struct {
	dbClient *database.DBClient
}

func NewTokenGenerator(dbClient *database.DBClient) *TokenGenrator {
	return &TokenGenrator{
		dbClient: dbClient,
	}
}

func (t *TokenGenrator) TokenGenerator(email, firstName, lastName, uuid string) (signedToken string, signedRefreshToken string, err error) {

	claims := &SignedDetails{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		UUID:      uuid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(JWT_Key))
	if err != nil {
		return "", "", err
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(JWT_Key))
	if err != nil {
		return "", "", err
	}

	return token, refreshToken, err
}

func (t *TokenGenrator) ValidateToken(signedtoken string) (claims *SignedDetails, msg string) {

	claims = &SignedDetails{}

	token, err := jwt.ParseWithClaims(signedtoken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWT_Key), nil
	})
	if err != nil {
		msg = err.Error()
		return
	}

	if !token.Valid {
		msg = "unauthorized access"
		return
	}

	if claims.StandardClaims.ExpiresAt < time.Now().Local().Unix() {
		msg = "token is expired"
		return
	}

	return claims, msg
}

func (t *TokenGenrator) UpdateAllTokens(signedtoken, signedrefreshtoken, user_id string) {

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var updateobj primitive.D
	updateobj = append(updateobj, bson.E{Key: "token", Value: signedtoken})
	updateobj = append(updateobj, bson.E{Key: "refresh_token", Value: signedrefreshtoken})
	updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateobj = append(updateobj, bson.E{Key: "updated_at", Value: updated_at})
	upsert := true
	filter := bson.M{"user_id": user_id}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	err := t.dbClient.UpdateOne(ctx, t.dbClient.GetUserCollection(), filter, bson.D{{Key: "$set", Value: updateobj}}, opt)
	if err != nil {
		log.Panic(err)
	}
}
