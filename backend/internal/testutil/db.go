package testutil

import (
	"bytes"
	"context"
	"daily-seed/internal/analytics"
	"daily-seed/internal/daily"
	"daily-seed/internal/habit"
	"daily-seed/internal/task"
	"log"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var DB *mongo.Database

func RunWithDB(m *testing.M) {
	ctx := context.Background()

	mongodbContainer, err := mongodb.Run(ctx, "mongo:5.0", mongodb.WithReplicaSet("rs0"))
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

	clientOpts := options.Client().ApplyURI(uri).SetDirect(true)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatalf("failed to connect to mongo: %s", err)
	}

	DB = client.Database("daily-seed-test")

	// Wait for replica set primary to be elected
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatalf("failed to ping primary: %s", err)
	}

	code := m.Run()

	if err := client.Disconnect(ctx); err != nil {
		log.Fatalf("failed to disconnect mongo: %s", err)
	}

	os.Exit(code)
}

func ClearDB(ctx context.Context) {
	if DB != nil {
		collections, err := DB.ListCollectionNames(ctx, bson.M{})
		if err == nil {
			for _, coll := range collections {
				_ = DB.Collection(coll).Drop(ctx)
			}
		}
	}
}

func SetupRouter(db *mongo.Database) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(gin.Recovery())

	habitStore := habit.NewHabitStore(db)
	taskStore := task.NewTaskStore(db)
	dailyStore := daily.NewDailyStore(db)

	ctx := context.Background()
	_ = habitStore.EnsureIndexes(ctx)
	_ = taskStore.EnsureIndexes(ctx)
	_ = dailyStore.EnsureIndexes(ctx)

	habitHandler := habit.NewHabitHandler(habitStore)
	taskHandler := task.NewTaskHandler(taskStore, dailyStore, dailyStore)
	dailyHandler := daily.NewDailyHandler(dailyStore, taskStore, habitStore)
	analyticsHandler := analytics.NewAnalyticsHandler(dailyStore, taskStore)

	v1 := r.Group("/api/v1")
	habitHandler.RegisterRoutes(v1)
	taskHandler.RegisterRoutes(v1)
	dailyHandler.RegisterRoutes(v1)
	analyticsHandler.RegisterRoutes(v1)

	return r
}

func DoRequest(r *gin.Engine, method, path string, body []byte) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func SeedHabit(ctx context.Context, db *mongo.Database, h habit.Habit) primitive.ObjectID {
	if h.ID.IsZero() {
		h.ID = primitive.NewObjectID()
	}
	if h.Status == "" {
		h.Status = "active"
	}
	_, _ = db.Collection("habits").InsertOne(ctx, h)
	return h.ID
}

func SeedTask(ctx context.Context, db *mongo.Database, t task.Task) primitive.ObjectID {
	if t.ID.IsZero() {
		t.ID = primitive.NewObjectID()
	}
	if t.Status == "" {
		t.Status = "active"
	}
	_, _ = db.Collection("tasks").InsertOne(ctx, t)
	return t.ID
}

func SeedDailyRecord(ctx context.Context, db *mongo.Database, rec daily.DailyRecord) primitive.ObjectID {
	if rec.ID.IsZero() {
		rec.ID = primitive.NewObjectID()
	}
	_, _ = db.Collection("dailyRecords").InsertOne(ctx, rec)
	return rec.ID
}
