package task

import "context"

// Task represents a task in the master Tasks collection.
type Task struct {
	ID         string         `json:"id" bson:"_id"`
	Section    string         `json:"section" bson:"section"` // "japanese" | "dev" | "self_dev"
	Title      string         `json:"title" bson:"title"`
	Type       string         `json:"type" bson:"type"` // "quantitative" | "boolean"
	Metrics    TaskMetrics    `json:"metrics" bson:"metrics"`
	Conditions TaskConditions `json:"conditions" bson:"conditions"`
	Status     string         `json:"status" bson:"status"` // "active" | "archived"
}

type TaskMetrics struct {
	DailyTarget int `json:"dailyTarget" bson:"dailyTarget"` // e.g. 10 (pages), 1 (boolean)
	WeeklyGoal  int `json:"weeklyGoal" bson:"weeklyGoal"`   // optional
}

type TaskService interface {
	List(ctx context.Context) ([]Task, error)
	Get(ctx context.Context, id string) (*Task, error)
	Create(ctx context.Context, task *Task) (*Task, error)
	Update(ctx context.Context, task *Task) (*Task, error)
	Archive(ctx context.Context, id string) error
}

type TaskRepository interface {
	FindActiveTasks(ctx context.Context) ([]Task, error)
	FindAll(ctx context.Context) ([]Task, error)
	FindByID(ctx context.Context, id string) (*Task, error)
	Create(ctx context.Context, task *Task) error
	Update(ctx context.Context, task *Task) error
	Delete(ctx context.Context, id string) error
}

type TaskConditions struct {
	Weather string `json:"weather" bson:"weather"` // "any" | "sunny" | "rainy" etc.
	Mode    string `json:"mode" bson:"mode"`       // "any" | "Growth" | "Rest" | "Work"
}
