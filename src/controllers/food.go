package controllers

import (
	"context"
	"fmt"
	"net/http"
	"restaurant-management-backend/src/databases"
	"restaurant-management-backend/src/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var foodCollection *mongo.Collection = databases.Open(databases.Client, "food")
var menuCollection *mongo.Collection = databases.Open(databases.Client, "menu")
var validate = validator.New()

func GetFoods() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

func GetFood() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		foodId := ctx.Param("foodId")
		var food models.Food

		err := foodCollection.FindOne(c, bson.M{"foodId": foodId}).Decode(&food)
		defer cancel()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "error occured while fetching the food item",
			})
		}
		ctx.JSON(http.StatusOK, food)
	}
}

func CreateFood() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var menu models.Menu
		var food models.Food
		if err := ctx.BindJSON((&food)); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		validationError := validate.Struct(food)
		if validationError != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": validationError.Error(),
			})
			return
		}

		err := menuCollection.FindOne(c, bson.M{
			"menuId": food.MenuId,
		})
		defer cancel()
		if err != nil {
			message := fmt.Sprintf("Menu was not found.")
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": message,
			})
			return
		}

		food.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		food.UpdateAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		food.ID = primitive.NewObjectID()
		food.FoodId = food.ID.Hex()

		var number = ToFixed(*food.Price, 2)
		food.Price = &number

		result, insertErr := foodCollection.InsertOne(c, food)
		if insertErr != nil {
			message := fmt.Sprintf("foodItem was not created")
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": message,
			})
			return
		}

		defer cancel()
		ctx.JSON(http.StatusOK, result)
	}
}

func UpdateFood() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

// func Round(number float64) int {

// }

func ToFixed(number float64, precision int) float64 {

}
