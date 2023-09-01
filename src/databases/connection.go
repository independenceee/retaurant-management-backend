package databases

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect() *mongo.Client {
	MONGO_DB_URL := ""
	fmt.Print(MONGO_DB_URL)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(MONGO_DB_URL))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connect to mongodb successfully !")
	return client
}

var Client *mongo.Client = Connect()

func Open(client *mongo.Client, collection string) *mongo.Collection {
	var query *mongo.Collection = client.Database("restaurant").Collection(collection)
	return query
}
