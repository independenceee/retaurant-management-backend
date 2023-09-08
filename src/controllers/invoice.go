package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"restaurant-management-backend/src/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
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

		invoiceView.InvoiceId = invoice.InvoiceId
		invoiceView.PaymentStatus = *&invoice.PaymentStatus
		invoiceView.PaymentDue = allOrderItems[0]["paymentDue"]
		invoiceView.TableNumber = allOrderItems[0]["tableNumber"]
		invoiceView.OrderDetails = allOrderItems[0]["orderItems"]

		ctx.JSON(http.StatusOK, invoiceView)
	}
}

func CreateInvoice() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var invoice models.Invoice
		if err := ctx.BindJSON(&invoice); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		var order models.Order

		err := orderCollection.FindOne(c, bson.M{
			"invoiceId": invoice.OrderId,
		}).Decode(&order)

		defer cancel()
		if err != nil {
			msg := fmt.Sprintf("order was not found")
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": msg,
			})
		}

		status := "PEDDING"
		if invoice.PaymentStatus == nil {
			invoice.PaymentStatus = &status
		}

		invoice.PaymentDueDate, _ = time.Parse(time.RFC3339, time.Now().AddDate(0, 0, 1).Format(time.RFC3339))
		invoice.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		invoice.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		invoice.ID = primitive.NewObjectID()
		invoice.InvoiceId = invoice.ID.Hex()
		validationErr := validate.Struct(invoice)

		if validationErr != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": validationErr.Error(),
			})
			return
		}

		result, insertErr := invoiceCollection.InsertOne(c, invoice)
		if insertErr != nil {
			msg := fmt.Sprintf("invoice item was not created")
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": msg,
			})

			return
		}

		defer cancel()

		ctx.JSON(http.StatusOK, result)
	}
}

func UpdateInvoice() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var invoice models.Invoice
		invoiceId := ctx.Param("id")

		if err := ctx.BindJSON(&invoice); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		filter := bson.M{
			"invoiceId": invoiceId,
		}

		var updateObject primitive.D

		if invoice.PaymentMethod != nil {

		}

		if invoice.PaymentStatus != nil {

		}

		invoice.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObject = append(updateObject, bson.E{
			"updatedAt", invoice.UpdatedAt,
		})

		upsert := true
		otp := options.UpdateOptions{
			Upsert: &upsert,
		}

		status := "PENDING"

		if invoice.PaymentStatus == nil {
			invoice.PaymentStatus = &status
		}

		result, err := invoiceCollection.UpdateOne(
			c,
			filter,
			bson.D{
				{"$set", updateObject},
			},
			&otp,
		)

		if err != nil {
			msg := fmt.Sprintf("invoice item update failed")
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": msg,
			})

			return
		}

		defer cancel()

		ctx.JSON(http.StatusOK, result)
	}
}
