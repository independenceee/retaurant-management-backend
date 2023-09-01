package routers

import (
	"restaurant-management-backend/src/controllers"

	"github.com/gin-gonic/gin"
)

func Invoice(router *gin.Engine) {
	router.GET("/invoices", controllers.GetFoods())
	router.GET("/invoices/:id", controllers.GetFood())
	router.POST("/invoices", controllers.CreateFood())
	router.PATCH("/invoices/:id", controllers.UpdateFood())
}
