package routers

import (
	"restaurant-management-backend/src/controllers"

	"github.com/gin-gonic/gin"
)

func Table(router *gin.Engine) {
	router.GET("/tables", controllers.GetTables())
	router.GET("/tables/:id", controllers.GetTable())
	router.POST("/tables", controllers.CreateTable())
	router.PATCH("/tables/:id", controllers.UpdateTable())
}
