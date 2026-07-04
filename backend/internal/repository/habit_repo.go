package repository

import (
	"context"

	"daily-seed/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)



type mongoHabitRepo struct {
	col *mongo.Collection
}

func NewHabitRepository(db *mongo.Database) model.HabitRepository {
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

func (r *mongoHabitRepo) FindAll(ctx context.Context) ([]model.Habit, error) {
	cursor, err := r.col.Find(ctx, bson.M{})
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

func (r *mongoHabitRepo) FindByID(ctx context.Context, id string) (*model.Habit, error) {
	var habit model.Habit
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&habit)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &habit, nil
}

func (r *mongoHabitRepo) Create(ctx context.Context, habit *model.Habit) error {
	_, err := r.col.InsertOne(ctx, habit)
	return err
}

func (r *mongoHabitRepo) Update(ctx context.Context, habit *model.Habit) error {
	filter := bson.M{"_id": habit.ID}
	update := bson.M{"$set": habit}
	_, err := r.col.UpdateOne(ctx, filter, update)
	return err
}

func (r *mongoHabitRepo) Delete(ctx context.Context, id string) error {
	_, err := r.col.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
