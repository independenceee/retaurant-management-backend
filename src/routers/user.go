package routers

import (
	"restaurant-management-backend/src/controllers"

	"github.com/gin-gonic/gin"
)

func User(router *gin.Engine) {
	router.GET("/users", controllers.GetUsers())
	router.GET("/users/:id", controllers.GetUser())
	router.POST("/users/login", controllers.Login())
	router.POST("/users/register", controllers.Register())
}
