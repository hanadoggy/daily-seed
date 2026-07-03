package model

// Task represents a task in the master Tasks collection.
type Task struct {
	ID         string         `json:"id" bson:"_id"`
	Section    string         `json:"section" bson:"section"`       // "japanese" | "dev" | "self_dev"
	Title      string         `json:"title" bson:"title"`
	Type       string         `json:"type" bson:"type"`             // "quantitative" | "boolean"
	Metrics    TaskMetrics    `json:"metrics" bson:"metrics"`
	Conditions TaskConditions `json:"conditions" bson:"conditions"`
	Status     string         `json:"status" bson:"status"`         // "active" | "archived"
}

type TaskMetrics struct {
	DailyTarget int `json:"dailyTarget" bson:"dailyTarget"`
	TotalTarget int `json:"totalTarget" bson:"totalTarget"`
}

type TaskConditions struct {
	Weather string `json:"weather" bson:"weather"` // "any" | "sunny" | "rainy" etc.
	Mode    string `json:"mode" bson:"mode"`       // "any" | "Growth" | "Rest" | "Work"
}
