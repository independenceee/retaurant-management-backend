package helpers

import (
	"context"
	"fmt"
	"log"
	"os"
	"restaurant-management-backend/src/databases"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SignedDetails struct {
	Email     string
	FirstName string
	LastName  string
	ID        string
	jwt.StandardClaims
}

var userCollection *mongo.Collection = databases.Open(databases.Client, "user")
var ACCESS_TOKEN_SECRET = os.Getenv("ACCESS_TOKEN_SECRET")
var REFRESH_TOKEN_SECRET = os.Getenv("REFRESH_TOKEN_SECRET")

func GenerateAllTokens(id string, email string, firstName string, lastName string) (accessToken string, refreshToken string, err error) {
	claims := &SignedDetails{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		ID:        id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	tokenAccess, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(ACCESS_TOKEN_SECRET))
	tokenRefresh, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(REFRESH_TOKEN_SECRET))

	if err != nil {
		log.Panic(err)
		return
	}

	return tokenAccess, tokenRefresh, err
}

func UpdateAllTokens(tokenAccess string, tokenRefresh string, id string) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var updateObject primitive.D

	updateObject = append(updateObject, bson.E{"accessToken", tokenAccess})
	updateObject = append(updateObject, bson.E{"refreshToken", tokenRefresh})

	updateAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObject = append(updateObject, bson.E{"updatedAt", updateAt})
	upsert := true

	filter := bson.M{"id": id}

	otp := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err := userCollection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{"$set", updateObject},
		},
		&otp,
	)

	defer cancel()

	if err != nil {
		log.Panic(err)
	}

	return
}

func ValidateToken(signedToken string) (claims *SignedDetails, message string) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(ACCESS_TOKEN_SECRET), nil
		},
	)

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		message = fmt.Sprintf("The token is invalid.")
		message = err.Error()
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		message = fmt.Sprintf("The token is invalid.")
		message = err.Error()
		return
	}

	return claims, message
}
