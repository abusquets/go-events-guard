package versions

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Up001(db *mongo.Database) error {
	userCollection := db.Collection("clients")
	indexModel := mongo.IndexModel{
		Keys:    bson.M{"code": 1},
		Options: options.Index().SetUnique(true),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := userCollection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return err
	}
	log.Println("Up001: unique index created for clients.code field")
	return nil
}

func Down001(db *mongo.Database) error {
	userCollection := db.Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := userCollection.Indexes().DropOne(ctx, "code_1")
	if err != nil {
		return err
	}
	log.Println("Down001: unique index removed for clients.code field")
	return nil
}
