package daily

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

// DailyRecord represents a single day's transaction data.
// The _id is a YYYY-MM-DD string in JST.
import "go.mongodb.org/mongo-driver/bson/primitive"

type DailyRecord struct {
	ID      primitive.ObjectID `json:"id" bson:"_id"`
	Date    string             `json:"date" bson:"date"`
	Context DayContext         `json:"context" bson:"context"`
	Tasks   []TaskEntry  `json:"tasks" bson:"tasks"`
	Habits  []HabitEntry `json:"habits" bson:"habits"`
	Journal Journal      `json:"journal" bson:"journal"`
}

type ContextMode string

const (
	ModeGrowth ContextMode = "Growth"
	ModeRest   ContextMode = "Rest"
	ModeOffice ContextMode = "Office"
	ModeRemote ContextMode = "Remote"
)

type DayContext struct {
	Mode    ContextMode `json:"mode" bson:"mode"` // "Growth" | "Rest" | "Office" | "Remote"
	Weather string      `json:"weather" bson:"weather"`
}

type UpdateDailyRecordRequest struct {
	Context *DayContextPatch `json:"context,omitempty"`
	Tasks   []TaskEntry      `json:"tasks,omitempty" binding:"omitempty,dive"`
	Habits  []HabitEntry     `json:"habits,omitempty" binding:"omitempty,dive"`
	Journal *JournalPatch    `json:"journal,omitempty"`
}

type DayContextPatch struct {
	Mode    *ContextMode `json:"mode,omitempty" binding:"omitempty,oneof=Growth Rest Office Remote"`
	Weather *string      `json:"weather,omitempty"`
}

type JournalPatch struct {
	OneLineReview  *string `json:"oneLineReview,omitempty"`
	ThreeLineDiary *string `json:"threeLineDiary,omitempty"`
}

type TaskEntry struct {
	TaskID       primitive.ObjectID `json:"taskId" bson:"taskId"`
	TargetAmount int                `json:"targetAmount" bson:"targetAmount"`
	ActualAmount int                `json:"actualAmount" bson:"actualAmount"`
	IsCompleted  bool               `json:"isCompleted" bson:"isCompleted"`
}

type HabitEntry struct {
	HabitID     primitive.ObjectID `json:"habitId" bson:"habitId"`
	IsCompleted bool               `json:"isCompleted" bson:"isCompleted"`
}

type Journal struct {
	OneLineReview  string `json:"oneLineReview" bson:"oneLineReview"`
	ThreeLineDiary string `json:"threeLineDiary" bson:"threeLineDiary"`
}

type DailyService interface {
	GetDailyRecord(ctx context.Context, date string) (*DailyRecord, error)
	UpdateDailyRecord(ctx context.Context, date string, req *UpdateDailyRecordRequest) (*DailyRecord, error)
	GetExistingRecordDates(ctx context.Context, year, month int) ([]string, error)
}

type DailyRecordRepository interface {
	FindByDate(ctx context.Context, date string) (*DailyRecord, error)
	Upsert(ctx context.Context, record *DailyRecord) error
	PatchByDate(ctx context.Context, date string, setFields bson.M) error
	EnsureIndexes(ctx context.Context) error
	FindBetweenDates(ctx context.Context, startDate string, endDate string) ([]*DailyRecord, error)
}
