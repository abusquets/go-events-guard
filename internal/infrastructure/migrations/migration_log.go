package migrations

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	migrationLogCollection = "migration_logs"
)

type MigrationLog struct {
	Version   string    `bson:"version"`
	AppliedAt time.Time `bson:"appliedAt"`
}

func AddMigrationLog(db *mongo.Database, version string) error {
	logCollection := db.Collection(migrationLogCollection)

	_, err := logCollection.InsertOne(context.Background(), MigrationLog{
		Version:   version,
		AppliedAt: time.Now(),
	})
	return err
}

func RemoveMigrationLog(db *mongo.Database, version string) error {
	logCollection := db.Collection(migrationLogCollection)

	_, err := logCollection.DeleteOne(context.Background(), bson.M{"version": version})
	return err
}

func HasMigration(db *mongo.Database, version string) (bool, error) {
	logCollection := db.Collection(migrationLogCollection)

	count, err := logCollection.CountDocuments(context.Background(), bson.M{"version": version})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
