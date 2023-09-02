package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderItem struct {
	ID          primitive.ObjectID `bson: "_id"`
	Quantity    *string            `json:"quanlity" 		validate:"required, eq=S|eq=M|eq=L"`
	UnitPrice   *float64           `json:"unitPrice" 		validate:"required"`
	CreatedAt   time.Time          `json:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt"`
	FoodId      *string            `json:"foodId" 			validate:"required"`
	OrderItemId string             `json: "orderItemId"`
	OrderId     string             `json: "orderId" 		validate: "required"`
}
