package task

import (
	"context"

	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoTaskRepo struct {
	client *mongo.Client
	col    *mongo.Collection
}

func NewTaskRepository(db *mongo.Database) TaskRepository {
	return &MongoTaskRepo{
		client: db.Client(),
		col:    db.Collection("tasks"),
	}
}

func (r *MongoTaskRepo) EnsureIndexes(ctx context.Context) error {
	indexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "status", Value: 1}},
	}
	_, err := r.col.Indexes().CreateOne(ctx, indexModel)
	return err
}

func (r *MongoTaskRepo) FindActiveTasks(ctx context.Context) ([]Task, error) {
	cursor, err := r.col.Find(ctx, bson.M{"status": "active"})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	tasks := make([]Task, 0)
	if err := cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *MongoTaskRepo) FindAll(ctx context.Context) ([]Task, error) {
	cursor, err := r.col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	tasks := make([]Task, 0)
	if err := cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *MongoTaskRepo) FindByID(ctx context.Context, id string) (*Task, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid task id format: %w", err)
	}
	var task Task
	err = r.col.FindOne(ctx, bson.M{"_id": oid}).Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &task, nil
}

func (r *MongoTaskRepo) Create(ctx context.Context, task *Task) error {
	_, err := r.col.InsertOne(ctx, task)
	return err
}

func (r *MongoTaskRepo) Update(ctx context.Context, task *Task) error {
	filter := bson.M{"_id": task.ID}
	update := bson.M{"$set": task}
	_, err := r.col.UpdateOne(ctx, filter, update)
	return err
}

func (r *MongoTaskRepo) MigrateTaskAtomic(ctx context.Context, archivedTask *Task, newTask *Task) error {
	session, err := r.client.StartSession()
	if err != nil {
		return fmt.Errorf("starting session: %w", err)
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		// Archive old task
		filter := bson.M{"_id": archivedTask.ID}
		update := bson.M{"$set": bson.M{"status": "archived"}}
		if _, err := r.col.UpdateOne(sessCtx, filter, update); err != nil {
			return nil, fmt.Errorf("archiving old task: %w", err)
		}

		// Insert new task
		if _, err := r.col.InsertOne(sessCtx, newTask); err != nil {
			return nil, fmt.Errorf("inserting new task: %w", err)
		}
		return nil, nil
	})
	return err
}
