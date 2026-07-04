package model

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

// DailyRecord represents a single day's transaction data.
// The _id is a YYYY-MM-DD string in JST.
type DailyRecord struct {
	ID      string       `json:"id" bson:"_id"`
	Date    string       `json:"date" bson:"date"`
	Context DayContext   `json:"context" bson:"context"`
	Tasks   []TaskEntry  `json:"tasks" bson:"tasks"`
	Habits  []HabitEntry `json:"habits" bson:"habits"`
	Journal Journal      `json:"journal" bson:"journal"`
}

type DayContext struct {
	Mode    string `json:"mode" bson:"mode"`       // "Growth" | "Rest" | "Work"
	Weather string `json:"weather" bson:"weather"`
}

type TaskEntry struct {
	TaskID       string `json:"taskId" bson:"taskId"`
	TargetAmount int    `json:"targetAmount" bson:"targetAmount"`
	ActualAmount int    `json:"actualAmount" bson:"actualAmount"`
	IsCompleted  bool   `json:"isCompleted" bson:"isCompleted"`
}

type HabitEntry struct {
	HabitID     string `json:"habitId" bson:"habitId"`
	IsCompleted bool   `json:"isCompleted" bson:"isCompleted"`
}

type Journal struct {
	OneLineReview  string `json:"oneLineReview" bson:"oneLineReview"`
	ThreeLineDiary string `json:"threeLineDiary" bson:"threeLineDiary"`
}

type DailyService interface {
	GetDailyRecord(ctx context.Context, date string) (*DailyRecord, error)
	UpdateDailyRecord(ctx context.Context, date string, patch map[string]interface{}) (*DailyRecord, error)
}

type DailyRecordRepository interface {
	FindByDate(ctx context.Context, date string) (*DailyRecord, error)
	Upsert(ctx context.Context, record *DailyRecord) error
	PatchByDate(ctx context.Context, date string, setFields bson.M) error
}
