package habit

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Habit represents a habit in the master Habits collection.
type Habit struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	Title    string             `json:"title" bson:"title"`
	Category string             `json:"category" bson:"category"`
	Status   string             `json:"status" bson:"status"` // active, archived
}
