package repository

import (
	"context"

	"daily-seed/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type TaskRepository interface {
	FindActiveTasks(ctx context.Context) ([]model.Task, error)
}

type mongoTaskRepo struct {
	col *mongo.Collection
}

func NewTaskRepository(db *mongo.Database) TaskRepository {
	return &mongoTaskRepo{col: db.Collection("tasks")}
}

func (r *mongoTaskRepo) FindActiveTasks(ctx context.Context) ([]model.Task, error) {
	cursor, err := r.col.Find(ctx, bson.M{"status": "active"})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tasks []model.Task
	if err := cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}
