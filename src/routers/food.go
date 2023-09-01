package routers

import (
	"restaurant-management-backend/src/controllers"

	"github.com/gin-gonic/gin"
)

func Food(router *gin.Engine) {
	router.GET("/foods", controllers.GetFoods())
	router.GET("/foods/:id", controllers.GetFood())
	router.POST("/foods", controllers.CreateFood())
	router.PATCH("/foods/:id", controllers.UpdateFood())
}
