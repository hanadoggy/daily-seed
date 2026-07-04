package habit

import "context"

// Habit represents a habit in the master Habits collection.
type Habit struct {
	ID       string `json:"id" bson:"_id"`
	Title    string `json:"title" bson:"title"`
	Category string `json:"category" bson:"category"`
	Status   string `json:"status" bson:"status"` // active, archived
}

type HabitService interface {
	List(ctx context.Context) ([]Habit, error)
	Get(ctx context.Context, id string) (*Habit, error)
	Create(ctx context.Context, habit *Habit) (*Habit, error)
	Update(ctx context.Context, habit *Habit) (*Habit, error)
	Archive(ctx context.Context, id string) error
}

type HabitRepository interface {
	FindActiveHabits(ctx context.Context) ([]Habit, error)
	FindAll(ctx context.Context) ([]Habit, error)
	FindByID(ctx context.Context, id string) (*Habit, error)
	Create(ctx context.Context, habit *Habit) error
	Update(ctx context.Context, habit *Habit) error
	Delete(ctx context.Context, id string) error
}
