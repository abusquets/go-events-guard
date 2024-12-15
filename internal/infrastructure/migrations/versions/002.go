package versions

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Up002(db *mongo.Database) error {
	eventCollection := db.Collection("events")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	index1 := mongo.IndexModel{
		Keys:    bson.D{{Key: "client_id", Value: 1}},
		Options: options.Index().SetName("events__client_id"),
	}

	index2 := mongo.IndexModel{
		Keys:    bson.D{{Key: "client_id", Value: 1}, {Key: "_id", Value: 1}},
		Options: options.Index().SetName("events__client_id__id"),
	}

	index3 := mongo.IndexModel{
		Keys:    bson.D{{Key: "client_id", Value: 1}, {Key: "type", Value: 1}},
		Options: options.Index().SetName("events__client_id__type"),
	}

	_, err := eventCollection.Indexes().CreateMany(ctx, []mongo.IndexModel{index1, index2, index3})
	if err != nil {
		return err
	}

	log.Print("Up002: created indexes for `events` collection")
	return nil
}

func Down002(db *mongo.Database) error {
	eventCollection := db.Collection("events")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := eventCollection.Indexes().DropOne(ctx, "events__client_id")
	if err != nil {
		return err
	}

	_, err = eventCollection.Indexes().DropOne(ctx, "events__client_id__id")
	if err != nil {
		return err
	}
	_, err = eventCollection.Indexes().DropOne(ctx, "events__client_id__type")
	if err != nil {
		return err
	}
	log.Print("Down002: indexes removed from `events` collection")
	return nil
}
