package repository

import (
	"context"

	"daily-seed/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type HabitRepository interface {
	FindActiveHabits(ctx context.Context) ([]model.Habit, error)
}

type mongoHabitRepo struct {
	col *mongo.Collection
}

func NewHabitRepository(db *mongo.Database) HabitRepository {
	return &mongoHabitRepo{col: db.Collection("habits")}
}

func (r *mongoHabitRepo) FindActiveHabits(ctx context.Context) ([]model.Habit, error) {
	cursor, err := r.col.Find(ctx, bson.M{"status": "active"})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var habits []model.Habit
	if err := cursor.All(ctx, &habits); err != nil {
		return nil, err
	}
	return habits, nil
}
