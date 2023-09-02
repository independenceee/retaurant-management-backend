package middlewares

import (
	"net/http"
	"restaurant-management-backend/src/helpers"

	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		clientToken := ctx.Request.Header.Get("token")
		if clientToken == "" {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Unauthorization",
			})
			ctx.Abort()
			return
		}

		claims, err := helpers.ValidateToken(clientToken)
		if err != "" {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		}

		ctx.Set("email", claims.Email)
		ctx.Set("firstName", claims.FirstName)
		ctx.Set("lastName", claims.LastName)
		ctx.Set("id", claims.ID)

		ctx.Next()
	}
}
