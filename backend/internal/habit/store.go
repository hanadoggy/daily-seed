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
		return nil, fmt.Errorf("find active habits: %w", err)
	}
	defer cursor.Close(ctx)

	habits := make([]Habit, 0)
	if err := cursor.All(ctx, &habits); err != nil {
		return nil, fmt.Errorf("decode active habits: %w", err)
	}
	return habits, nil
}

func (s *HabitStore) FindAll(ctx context.Context) ([]Habit, error) {
	cursor, err := s.col.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("find all habits: %w", err)
	}
	defer cursor.Close(ctx)

	habits := make([]Habit, 0)
	if err := cursor.All(ctx, &habits); err != nil {
		return nil, fmt.Errorf("decode habits: %w", err)
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
		return nil, fmt.Errorf("find habit by id: %w", err)
	}
	return &habit, nil
}

func (s *HabitStore) Create(ctx context.Context, habit *Habit) error {
	if _, err := s.col.InsertOne(ctx, habit); err != nil {
		return fmt.Errorf("create habit: %w", err)
	}
	return nil
}

func (s *HabitStore) Update(ctx context.Context, habit *Habit) error {
	filter := bson.M{"_id": habit.ID}
	update := bson.M{"$set": habit}
	if _, err := s.col.UpdateOne(ctx, filter, update); err != nil {
		return fmt.Errorf("update habit: %w", err)
	}
	return nil
}

func (s *HabitStore) EnsureIndexes(ctx context.Context) error {
	indexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "status", Value: 1}},
	}
	if _, err := s.col.Indexes().CreateOne(ctx, indexModel); err != nil {
		return fmt.Errorf("ensure habit indexes: %w", err)
	}
	return nil
}
