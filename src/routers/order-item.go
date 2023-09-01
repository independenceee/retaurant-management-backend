package routers

import (
	"restaurant-management-backend/src/controllers"

	"github.com/gin-gonic/gin"
)

func OrderItem(router *gin.Engine) {
	router.GET("/order-items", controllers.GetOrderItems())
	router.GET("/order-items/:id", controllers.GetOrderItem())
	router.GET("/order-items-order/:id", controllers.GetOrderItemByOrder())
	router.POST("/order-items", controllers.CreateOrderItems())
	router.PATCH("/order-items/:id", controllers.UpdateOrderItem())
}
