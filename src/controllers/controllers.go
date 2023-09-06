package controllers

import (
	"restaurant-management-backend/src/databases"

	"go.mongodb.org/mongo-driver/mongo"
)

var foodCollection *mongo.Collection = databases.Open(databases.Client, "food")
var menuCollection *mongo.Collection = databases.Open(databases.Client, "menu")
var invoiceCollection *mongo.Collection = databases.Open(databases.Client, "invoice")
var orderItemCollection *mongo.Collection = databases.Open(databases.Client, "orderItem")
var orderCollection *mongo.Collection = databases.Open(databases.Client, "order")
var tableCollection *mongo.Collection = databases.Open(databases.Client, "table")
var userCollection *mongo.Collection = databases.Open(databases.Client, "user")
