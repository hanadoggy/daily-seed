package task

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TaskStore struct {
	client *mongo.Client
	col    *mongo.Collection
}

func NewTaskStore(db *mongo.Database) *TaskStore {
	return &TaskStore{
		client: db.Client(),
		col:    db.Collection("tasks"),
	}
}

func (s *TaskStore) EnsureIndexes(ctx context.Context) error {
	indexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "status", Value: 1}},
	}
	if _, err := s.col.Indexes().CreateOne(ctx, indexModel); err != nil {
		return fmt.Errorf("ensure task indexes: %w", err)
	}
	return nil
}

func (s *TaskStore) FindActiveTasks(ctx context.Context) ([]Task, error) {
	cursor, err := s.col.Find(ctx, bson.M{"status": "active"})
	if err != nil {
		return nil, fmt.Errorf("find active tasks: %w", err)
	}
	defer cursor.Close(ctx)

	tasks := make([]Task, 0)
	if err := cursor.All(ctx, &tasks); err != nil {
		return nil, fmt.Errorf("decode active tasks: %w", err)
	}
	return tasks, nil
}

func (s *TaskStore) FindAll(ctx context.Context) ([]Task, error) {
	cursor, err := s.col.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("find all tasks: %w", err)
	}
	defer cursor.Close(ctx)

	tasks := make([]Task, 0)
	if err := cursor.All(ctx, &tasks); err != nil {
		return nil, fmt.Errorf("decode tasks: %w", err)
	}
	return tasks, nil
}

func (s *TaskStore) FindByID(ctx context.Context, id string) (*Task, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid task id format: %w", err)
	}
	var task Task
	err = s.col.FindOne(ctx, bson.M{"_id": oid}).Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("find task by id: %w", err)
	}
	return &task, nil
}

func (s *TaskStore) Create(ctx context.Context, task *Task) error {
	if _, err := s.col.InsertOne(ctx, task); err != nil {
		return fmt.Errorf("create task: %w", err)
	}
	return nil
}

func (s *TaskStore) Update(ctx context.Context, task *Task) error {
	filter := bson.M{"_id": task.ID}
	update := bson.M{
		"$set": bson.M{
			"section":    task.Section,
			"title":      task.Title,
			"type":       task.Type,
			"unit":       task.Unit,
			"metrics":    task.Metrics,
			"conditions": task.Conditions,
			"status":     task.Status,
			"startDate":  task.StartDate,
			"endDate":    task.EndDate,
		},
	}
	if _, err := s.col.UpdateOne(ctx, filter, update); err != nil {
		return fmt.Errorf("update task: %w", err)
	}
	return nil
}

func (s *TaskStore) MigrateTaskAtomic(ctx context.Context, archivedTask *Task, newTask *Task) error {
	session, err := s.client.StartSession()
	if err != nil {
		return fmt.Errorf("starting session: %w", err)
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		// Archive old task
		filter := bson.M{"_id": archivedTask.ID}
		update := bson.M{"$set": bson.M{
			"status":  "archived",
			"endDate": archivedTask.EndDate,
		}}
		if _, err := s.col.UpdateOne(sessCtx, filter, update); err != nil {
			return nil, fmt.Errorf("archiving old task: %w", err)
		}

		// Insert new task
		if _, err := s.col.InsertOne(sessCtx, newTask); err != nil {
			return nil, fmt.Errorf("inserting new task: %w", err)
		}
		return nil, nil
	})
	return err
}
