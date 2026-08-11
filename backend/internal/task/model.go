package task

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Task represents a task in the master Tasks collection.
type Task struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	Section    string             `json:"section" bson:"section"` // "japanese" | "dev" | "self_dev" | "exercise"
	Title      string             `json:"title" bson:"title"`
	Type       string             `json:"type" bson:"type"` // "quantitative" | "boolean"
	Unit       string             `json:"unit" bson:"unit"`
	Metrics    TaskMetrics        `json:"metrics" bson:"metrics"`
	Conditions TaskConditions     `json:"conditions" bson:"conditions"`
	Status     string             `json:"status" bson:"status"` // "active" | "archived"
	StartDate  string             `json:"startDate" bson:"startDate"`
	EndDate    string             `json:"endDate,omitempty" bson:"endDate,omitempty"`
}

type TaskMetrics struct {
	DailyTarget int `json:"dailyTarget" bson:"dailyTarget"` // e.g. 10 (pages), 1 (boolean)
	TotalTarget int `json:"totalTarget" bson:"totalTarget"` // lifetime goal (0 = no limit)
}

type TaskConditions struct {
	Weather []string `json:"weather" bson:"weather"` // ["sunny", "rainy"]
	Mode    []string `json:"mode" bson:"mode"`       // ["Growth", "Rest", "Office", "Remote"]
}

var (
	validSections  = map[string]bool{"japanese": true, "dev": true, "self_dev": true, "exercise": true}
	validTaskTypes = map[string]bool{"quantitative": true, "boolean": true}
	validWeathers  = map[string]bool{"sunny": true, "rainy": true}
	validModes     = map[string]bool{"Growth": true, "Rest": true, "Office": true, "Remote": true}
)

func (t *Task) Validate() error {
	if strings.TrimSpace(t.Title) == "" {
		return fmt.Errorf("title is required")
	}
	if len(t.Title) > 200 {
		return fmt.Errorf("title must not exceed 200 characters")
	}
	if !validSections[t.Section] {
		return fmt.Errorf("section must be one of: japanese, dev, self_dev, exercise")
	}
	if !validTaskTypes[t.Type] {
		return fmt.Errorf("type must be one of: quantitative, boolean")
	}
	if strings.TrimSpace(t.Unit) == "" {
		return fmt.Errorf("unit is required")
	}
	if len(t.Unit) > 50 {
		return fmt.Errorf("unit must not exceed 50 characters")
	}
	if t.Type == "quantitative" && t.Metrics.DailyTarget <= 0 {
		return fmt.Errorf("dailyTarget must be positive for quantitative tasks")
	}
	if t.Metrics.TotalTarget < 0 {
		return fmt.Errorf("totalTarget cannot be negative")
	}
	if t.StartDate == "" {
		return fmt.Errorf("startDate is required")
	}
	if _, err := time.Parse("2006-01-02", t.StartDate); err != nil {
		return fmt.Errorf("startDate must be in YYYY-MM-DD format")
	}

	if len(t.Conditions.Weather) == 0 {
		return fmt.Errorf("at least one weather condition is required")
	}
	if len(t.Conditions.Mode) == 0 {
		return fmt.Errorf("at least one mode condition is required")
	}
	for _, w := range t.Conditions.Weather {
		if !validWeathers[w] {
			return fmt.Errorf("invalid weather value: %s", w)
		}
	}
	for _, m := range t.Conditions.Mode {
		if !validModes[m] {
			return fmt.Errorf("invalid mode value: %s", m)
		}
	}

	return nil
}

// TaskProgress represents a task's cumulative progress across all daily records.
type TaskProgress struct {
	TaskID         primitive.ObjectID `json:"taskId"`
	Title          string             `json:"title"`
	Type           string             `json:"type"`
	TotalTarget    int                `json:"totalTarget"`
	TotalCompleted int                `json:"totalCompleted"`
	Percentage     float64            `json:"percentage"`
}

// MigrateTaskRequest is the request payload for a task migration operation.
type MigrateTaskRequest struct {
	CompletionDate string `json:"completionDate"`
}

// MigrationResult is the response from a task migration operation.
type MigrationResult struct {
	ArchivedTask Task `json:"archivedTask"`
	NewTask      Task `json:"newTask"`
}

// TaskProgressAggregator provides cumulative progress data from daily records.
// Defined here to avoid circular imports with the daily package.
type TaskProgressAggregator interface {
	SumTaskProgressByIDs(ctx context.Context, taskIDs []primitive.ObjectID) (map[primitive.ObjectID]int, error)
}

// DailyCleaner cleans task records from daily records.
// Defined here to avoid circular imports with the daily package.
type DailyCleaner interface {
	RemoveTaskFromRecordsBeforeDate(ctx context.Context, taskID primitive.ObjectID, date string) error
}
