package habit

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type HabitStore struct {
	col *mongo.Collection
}

func NewHabitStore(db *mongo.Database) *HabitStore {
	return &HabitStore{col: db.Collection("habits")}
}

func (s *HabitStore) FindActiveHabits(ctx context.Context) ([]Habit, error) {
	cursor, err := s.col.Find(ctx, bson.M{"status": "active"})
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

func (s *HabitStore) FindAll(ctx context.Context) ([]Habit, error) {
	cursor, err := s.col.Find(ctx, bson.M{})
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

func (s *HabitStore) FindByID(ctx context.Context, id string) (*Habit, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid habit id format: %w", err)
	}
	var habit Habit
	err = s.col.FindOne(ctx, bson.M{"_id": oid}).Decode(&habit)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &habit, nil
}

func (s *HabitStore) Create(ctx context.Context, habit *Habit) error {
	_, err := s.col.InsertOne(ctx, habit)
	return err
}

func (s *HabitStore) Update(ctx context.Context, habit *Habit) error {
	filter := bson.M{"_id": habit.ID}
	update := bson.M{"$set": habit}
	_, err := s.col.UpdateOne(ctx, filter, update)
	return err
}

func (s *HabitStore) EnsureIndexes(ctx context.Context) error {
	indexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "status", Value: 1}},
	}
	_, err := s.col.Indexes().CreateOne(ctx, indexModel)
	return err
}
