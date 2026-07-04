package testutil

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func RunWithDB(m *testing.M) {
	ctx := context.Background()

	mongodbContainer, err := mongodb.Run(ctx, "mongo:7.0")
	if err != nil {
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

	DB = client.Database("daily-seed-test")

	code := m.Run()

	if err := client.Disconnect(ctx); err != nil {
		log.Fatalf("failed to disconnect mongo: %s", err)
	}

	os.Exit(code)
}

func ClearDB(ctx context.Context) {
	if DB != nil {
		_ = DB.Collection("tasks").Drop(ctx)
		_ = DB.Collection("habits").Drop(ctx)
		_ = DB.Collection("dailyRecords").Drop(ctx)
	}
}
