package controllers

import (
	"context"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func GetUsers() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		recordPerPage, err := strconv.Atoi(ctx.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}

		page, err1 := strconv.Atoi(ctx.Query("page"))
		if err1 != nil || page < 1 {
			page = 1
		}

		startIndex := (page -1 )* recordPerPage
		startIndex, err = strconv.Atoi(ctx.Query("startIndex"))
	}
}

func GetUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

func Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

func Register() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}
