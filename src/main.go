package main

import (
	"os"
	"restaurant-management-backend/src/middlewares"
	"restaurant-management-backend/src/routers"

	"github.com/gin-gonic/gin"
)

func main() {
	var PORT string = os.Getenv("PORT")
	if PORT == "" {
		PORT = "8000"
	}

	router := gin.New()
	router.Use(gin.Logger())

	routers.User(router)

	router.Use(middlewares.Authentication())

	routers.Food(router)
	routers.Menu(router)
	routers.Table(router)
	routers.Order(router)
	routers.OrderItem(router)
	routers.Invoice(router)
	router.Run(":" + PORT)

}
