package controllers

import (
	"context"
	"log"
	"net/http"
	"restaurant-management-backend/src/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type InvoiceViewFormat struct {
	InvoiceId      string
	PaymentMethod  string
	OrderId        string
	PaymentStatus  *string
	PaymentDue     interface{}
	TableNumber    interface{}
	PaymentDueDate time.Time
	OrderDetails   interface{}
}

func GetInvoices() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		result, err := invoiceCollection.Find(context.TODO(), bson.M{})

		defer cancel()

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "error occured while invoice items",
			})
		}

		var allInvoices []bson.M

		if err = result.All(c, &allInvoices); err != nil {
			log.Fatal(err)

		}

		ctx.JSON(http.StatusOK, allInvoices)
	}
}

func GetInvoice() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		invoiceId := ctx.Param("id")
		var invoice models.Invoice

		err := invoiceCollection.FindOne(c, bson.M{
			"id": invoiceId,
		}).Decode(&invoice)

		defer cancel()

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "error occured while listing invoice items",
			})
		}

		var invoiceView InvoiceViewFormat

		allOrderItems, err := ItemsByOrder(invoice.OrderId)
		invoiceView.OrderId = invoice.OrderId
		invoiceView.PaymentDueDate = invoice.PaymentDueDate
		invoiceView.PaymentMethod = "null"

		if invoice.PaymentMethod != nil {
			invoiceView.PaymentMethod = *invoice.PaymentMethod
		}

		// 11
	}
}

func CreateInvoice() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

func UpdateInvoice() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}
