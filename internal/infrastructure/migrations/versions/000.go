package versions

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Up000(db *mongo.Database) error {
	userCollection := db.Collection("users")
	indexModel := mongo.IndexModel{
		Keys:    bson.M{"email": 1},
		Options: options.Index().SetUnique(true),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := userCollection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return err
	}
	log.Println("Up000: unique index created for users.email field")
	return nil
}

func Down000(db *mongo.Database) error {
	userCollection := db.Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := userCollection.Indexes().DropOne(ctx, "email_1")
	if err != nil {
		return err
	}
	log.Println("Down000: unique index removed for users.email field")
	return nil
}
