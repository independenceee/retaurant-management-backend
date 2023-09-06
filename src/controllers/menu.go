package controllers

import (
	"context"
	"fmt"
	"net/http"
	"restaurant-management-backend/src/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Get All Menu
func GetMenus() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

// Get Menu
func GetMenu() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		menuId := ctx.Param("menu_id")
		var menu models.Menu

		err := menuCollection.FindOne(c, bson.M{
			"id": menuId,
		}).Decode(&menu)

		defer cancel()

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "error occurred while fetching the menu",
			})
		}

		ctx.JSON(http.StatusOK, menu)
	}
}

func CreateMenu() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var menu models.Menu
		var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		if err := ctx.BindJSON(&menu); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		validatorErr := validate.Struct(menu)
		if validatorErr != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": validatorErr.Error(),
			})
			return
		}

		menu.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		menu.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		menu.ID = primitive.NewObjectID()
		menu.MenuId = menu.ID.Hex()

		result, insertErr := menuCollection.InsertOne(c, menu)
		if insertErr != nil {
			msg := fmt.Sprintf("Menu item was not created")
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": msg,
			})
			return
		}

		defer cancel()
		ctx.JSON(http.StatusOK, result)
		defer cancel()
	}
}

func UpdateMenu() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var menu models.Menu

		if err := ctx.BindJSON(&menu); err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{
				"error": err.Error(),
			})
			return
		}

		menuId := ctx.Param("id")
		filter := bson.M{
			"menu_id": menuId,
		}

		var updateObject primitive.D

		if menu.StartDate != nil && menu.EndDate != nil {
			if !InTimeSpan(*menu.StartDate, *menu.EndDate, time.Now()) {
				msg := "Kindly retype the time"
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": msg,
				})
				defer cancel()
				return
			}

			updateObject = append(updateObject, bson.E{
				"startDate", menu.StartDate,
			})
			updateObject = append(updateObject, bson.E{
				"endDate", menu.EndDate,
			})

			if menu.Name != "" {
				updateObject = append(updateObject, bson.E{
					"name", menu.Name,
				})
			}

			if menu.Category != "" {
				updateObject = append(updateObject, bson.E{
					"name", menu.Category,
				})
			}

			menu.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			updateObject = append(updateObject, bson.E{
				"updatedAt", menu.UpdatedAt,
			})

			upsert := true

			otp := options.UpdateOptions{
				Upsert: &upsert,
			}

			result, err := menuCollection.UpdateOne(
				c,
				filter,
				bson.D{
					{"$set", updateObject},
				},
				&otp,
			)

			if err != nil {
				msg := "Menu update failed"
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": msg,
				})
			}
			defer cancel()
			ctx.JSON(http.StatusOK, result)
		}
	}
}

func InTimeSpan() bool {

}
