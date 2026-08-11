package daily

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DailyRecord represents a single day's transaction data.
// The _id is a YYYY-MM-DD string in JST.
type DailyRecord struct {
	ID      primitive.ObjectID `json:"id" bson:"_id"`
	Date    string             `json:"date" bson:"date"`
	Context DayContext         `json:"context" bson:"context"`
	Tasks   []TaskEntry        `json:"tasks" bson:"tasks"`
	Habits  []HabitEntry       `json:"habits" bson:"habits"`
	Journal Journal            `json:"journal" bson:"journal"`
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

var (
	validWeathers = map[string]bool{"sunny": true, "rainy": true}
	validModesMap = map[ContextMode]bool{
		ModeGrowth: true, ModeRest: true, ModeOffice: true, ModeRemote: true,
	}
)

func (r *UpdateDailyRecordRequest) Validate() error {
	if r.Context != nil {
		if r.Context.Weather != nil && !validWeathers[*r.Context.Weather] {
			return fmt.Errorf("weather must be one of: sunny, rainy")
		}
		if r.Context.Mode != nil && !validModesMap[*r.Context.Mode] {
			return fmt.Errorf("mode must be one of: Growth, Rest, Office, Remote")
		}
	}
	for i, t := range r.Tasks {
		if t.ActualAmount < 0 {
			return fmt.Errorf("tasks[%d].actualAmount cannot be negative", i)
		}
		if t.TargetAmount < 0 {
			return fmt.Errorf("tasks[%d].targetAmount cannot be negative", i)
		}
	}
	return nil
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
