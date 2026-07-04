package task

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoTaskRepo struct {
	col *mongo.Collection
}

func NewTaskRepository(db *mongo.Database) TaskRepository {
	return &MongoTaskRepo{col: db.Collection("tasks")}
}

func (r *MongoTaskRepo) FindActiveTasks(ctx context.Context) ([]Task, error) {
	cursor, err := r.col.Find(ctx, bson.M{"status": "active"})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tasks []Task
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

	var tasks []Task
	if err := cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *MongoTaskRepo) FindByID(ctx context.Context, id string) (*Task, error) {
	var task Task
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&task)
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

func (r *MongoTaskRepo) Delete(ctx context.Context, id string) error {
	_, err := r.col.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
