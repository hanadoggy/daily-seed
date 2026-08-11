package habit

import (
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Habit represents a habit in the master Habits collection.
type Habit struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	Title    string             `json:"title" bson:"title"`
	Category string             `json:"category" bson:"category"`
	Status   string             `json:"status" bson:"status"` // active, archived
}

func (h *Habit) Validate() error {
	if strings.TrimSpace(h.Title) == "" {
		return fmt.Errorf("title is required")
	}
	if len(h.Title) > 200 {
		return fmt.Errorf("title must not exceed 200 characters")
	}
	if strings.TrimSpace(h.Category) == "" {
		return fmt.Errorf("category is required")
	}
	return nil
}

