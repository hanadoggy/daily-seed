package habit

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoHabitRepo struct {
	col *mongo.Collection
}

func NewHabitRepository(db *mongo.Database) HabitRepository {
	return &MongoHabitRepo{col: db.Collection("habits")}
}

func (r *MongoHabitRepo) FindActiveHabits(ctx context.Context) ([]Habit, error) {
	cursor, err := r.col.Find(ctx, bson.M{"status": "active"})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	habits := make([]Habit, 0)
	if err := cursor.All(ctx, &habits); err != nil {
		return nil, err
	}
	return habits, nil
}

func (r *MongoHabitRepo) FindAll(ctx context.Context) ([]Habit, error) {
	cursor, err := r.col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	habits := make([]Habit, 0)
	if err := cursor.All(ctx, &habits); err != nil {
		return nil, err
	}
	return habits, nil
}

func (r *MongoHabitRepo) FindByID(ctx context.Context, id string) (*Habit, error) {
	var habit Habit
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&habit)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &habit, nil
}

func (r *MongoHabitRepo) Create(ctx context.Context, habit *Habit) error {
	_, err := r.col.InsertOne(ctx, habit)
	return err
}

func (r *MongoHabitRepo) Update(ctx context.Context, habit *Habit) error {
	filter := bson.M{"_id": habit.ID}
	update := bson.M{"$set": habit}
	_, err := r.col.UpdateOne(ctx, filter, update)
	return err
}

func (r *MongoHabitRepo) Delete(ctx context.Context, id string) error {
	_, err := r.col.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
