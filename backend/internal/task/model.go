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
	TotalTarget int `json:"totalTarget" bson:"totalTarget"`  // lifetime goal (0 = no limit)
}

type TaskConditions struct {
	Weather string `json:"weather" bson:"weather"` // "any" | "sunny" | "rainy"
	Mode    string `json:"mode" bson:"mode"`       // "any" | "Growth" | "Rest" | "Office" | "Remote"
}

// TaskProgress represents a task's cumulative progress across all daily records.
type TaskProgress struct {
	TaskID         string  `json:"taskId"`
	Title          string  `json:"title"`
	TotalTarget    int     `json:"totalTarget"`
	TotalCompleted int     `json:"totalCompleted"`
	Percentage     float64 `json:"percentage"`
}

// MigrationResult is the response from a task migration operation.
type MigrationResult struct {
	ArchivedTask Task `json:"archivedTask"`
	NewTask      Task `json:"newTask"`
}

type TaskService interface {
	List(ctx context.Context) ([]Task, error)
	Get(ctx context.Context, id string) (*Task, error)
	Create(ctx context.Context, task *Task) (*Task, error)
	Update(ctx context.Context, task *Task) (*Task, error)
	Archive(ctx context.Context, id string) error
	GetProgressForActiveTasks(ctx context.Context) ([]TaskProgress, error)
	MigrateTask(ctx context.Context, id string) (*MigrationResult, error)
}

type TaskRepository interface {
	FindActiveTasks(ctx context.Context) ([]Task, error)
	FindAll(ctx context.Context) ([]Task, error)
	FindByID(ctx context.Context, id string) (*Task, error)
	Create(ctx context.Context, task *Task) error
	Update(ctx context.Context, task *Task) error
	Delete(ctx context.Context, id string) error
}

// TaskProgressAggregator provides cumulative progress data from daily records.
// Defined here to avoid circular imports with the daily package.
type TaskProgressAggregator interface {
	SumTaskProgress(ctx context.Context, taskID string) (int, error)
}

