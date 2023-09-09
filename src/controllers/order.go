package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"restaurant-management-backend/src/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetOrders() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		results, err := orderCollection.Find(context.TODO(), bson.M{})
		defer cancel()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "error occured while listing order items",
			})
		}

		var allOrders []bson.M
		if err = results.All(c, &allOrders); err != nil {
			log.Fatal(err)
		}

		ctx.JSON(http.StatusOK, allOrders)
	}
}

func GetOrder() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		foodId := ctx.Param("id")
		var food models.Order
		err := foodCollection.FindOne(c, bson.M{
			"orderId": foodId,
		}).Decode(&food)
		defer cancel()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "occured while  fetching the food item",
			})

		}

		ctx.JSON(http.StatusOK, food)
	}
}

func CreateOrder() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var table models.Table
		var order models.Order

		if err := ctx.BindJSON(&order); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		validateErr := validate.Struct(order)

		if validateErr != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": validateErr.Error(),
			})

			return
		}

		if order.TableId != nil {
			err := tableCollection.FindOne(c, bson.M{
				"tableId": order.TableId,
			}).Decode(&table)
			defer cancel()

			if err != nil {
				message := fmt.Sprintf("Table not found")
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": message,
				})

				return
			}
		}

		order.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		order.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		order.ID = primitive.NewObjectID()
		order.OrderId = order.ID.Hex()

		result, insertErr := orderCollection.InsertOne(c, order)
		if insertErr != nil {
			message := fmt.Sprintf("order item was not created")
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": message,
			})

			return
		}

		defer cancel()
		ctx.JSON(http.StatusOK, result)

	}
}

func UpdateOrder() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var table models.Table
		var order models.Order
		var food models.Food
		var updateObject primitive.D

		orderId := ctx.Param("orderId")

		if err := ctx.BindJSON(&order); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if order.TableId != nil {
			err := menuCollection.FindOne(ctx, bson.M{
				"tableId": order.TableId,
			}).Decode(&table)

			defer cancel()

			if err != nil {
				message := fmt.Sprintf("Menu was not found!")
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": message,
				})
				return
			}
			updateObject = append(updateObject, bson.E{
				"menu", order.TableId,
			})
		}

		order.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObject = append(updateObject, bson.E{
			"updatedAt", food.UpdateAt,
		})

		upsert := true
		filter := bson.M{
			"orderId": orderId,
		}

		otp := options.UpdateOptions{
			Upsert: &upsert,
		}
		result, err := orderCollection.UpdateOne(
			c,
			filter,
			bson.D{
				{"$st", updateObject},
			},
			&otp,
		)

		if err != nil {
			message := fmt.Sprintf("order item update failed")
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": message,
			})
			return
		}

		defer cancel()

		ctx.JSON(http.StatusOK, result)
	}
}

func OrderItemOrder(order models.Order) string {
	var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	order.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.ID = primitive.NewObjectID()
	order.OrderId = order.ID.Hex()

	orderCollection.InsertOne(c, order)
	defer cancel()
	return order.OrderId
}
