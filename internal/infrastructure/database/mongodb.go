package database

import (
	"context"
	"eventsguard/internal/infrastructure/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectToMongo(config *config.AppConfig) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(config.MongoDBUri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}
	return client, nil
}
