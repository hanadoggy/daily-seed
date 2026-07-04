package repository_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var testDB *mongo.Database

func TestMain(m *testing.M) {
	ctx := context.Background()

	mongodbContainer, err := mongodb.Run(ctx, "mongo:7.0")
	if err != nil {
		// Use fallback if the system doesn't have Docker running, or we just fail.
		log.Fatalf("failed to start container: %s", err)
	}

	defer func() {
		if err := mongodbContainer.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate container: %s", err)
		}
	}()

	uri, err := mongodbContainer.ConnectionString(ctx)
	if err != nil {
		log.Fatalf("failed to get connection string: %s", err)
	}

	clientOpts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatalf("failed to connect to mongo: %s", err)
	}

	testDB = client.Database("daily-seed-test")

	code := m.Run()

	if err := client.Disconnect(ctx); err != nil {
		log.Fatalf("failed to disconnect mongo: %s", err)
	}

	os.Exit(code)
}

func clearDB(ctx context.Context) {
	if testDB != nil {
		_ = testDB.Collection("tasks").Drop(ctx)
		_ = testDB.Collection("habits").Drop(ctx)
		_ = testDB.Collection("dailyRecords").Drop(ctx)
	}
}
