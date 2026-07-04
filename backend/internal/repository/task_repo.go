package repository

import (
	"context"

	"daily-seed/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)



type mongoTaskRepo struct {
	col *mongo.Collection
}

func NewTaskRepository(db *mongo.Database) model.TaskRepository {
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

func (r *mongoTaskRepo) FindAll(ctx context.Context) ([]model.Task, error) {
	cursor, err := r.col.Find(ctx, bson.M{})
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

func (r *mongoTaskRepo) FindByID(ctx context.Context, id string) (*model.Task, error) {
	var task model.Task
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &task, nil
}

func (r *mongoTaskRepo) Create(ctx context.Context, task *model.Task) error {
	_, err := r.col.InsertOne(ctx, task)
	return err
}

func (r *mongoTaskRepo) Update(ctx context.Context, task *model.Task) error {
	filter := bson.M{"_id": task.ID}
	update := bson.M{"$set": task}
	_, err := r.col.UpdateOne(ctx, filter, update)
	return err
}

func (r *mongoTaskRepo) Delete(ctx context.Context, id string) error {
	_, err := r.col.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
