package daily

import (
	"context"
	"daily-seed/internal/habit"
	"daily-seed/internal/task"
	"fmt"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DailyServiceImpl struct {
	dailyRepo DailyRecordRepository
	taskRepo  task.TaskRepository
	habitRepo habit.HabitRepository
}

func NewDailyService(
	dailyRepo DailyRecordRepository,
	taskRepo task.TaskRepository,
	habitRepo habit.HabitRepository,
) DailyService {
	return &DailyServiceImpl{
		dailyRepo: dailyRepo,
		taskRepo:  taskRepo,
		habitRepo: habitRepo,
	}
}

func validateDate(date string) error {
	if _, err := time.Parse("2006-01-02", date); err != nil {
		return fmt.Errorf("invalid date format, expected YYYY-MM-DD")
	}
	return nil
}

// GetDailyRecord retrieves the DailyRecord for the given date.
// If no record exists, it generates one from the currently active tasks and habits.
func (s *DailyServiceImpl) GetDailyRecord(ctx context.Context, date string) (*DailyRecord, error) {
	if err := validateDate(date); err != nil {
		return nil, err
	}
	record, err := s.dailyRepo.FindByDate(ctx, date)
	if err != nil {
		return nil, fmt.Errorf("finding daily record: %w", err)
	}

	if record != nil {
		updated, err := s.appendMissingEntries(ctx, record)
		if err != nil {
			slog.Warn("failed to append missing entries to daily record", slog.String("error", err.Error()))
		} else if updated {
			if err := s.dailyRepo.Upsert(ctx, record); err != nil {
				return nil, fmt.Errorf("persisting updated daily record: %w", err)
			}
		}
		return record, nil
	}

	slog.Info("generating new daily record", slog.String("date", date))
	record, err = s.generateDailyRecord(ctx, date)
	if err != nil {
		return nil, fmt.Errorf("generating daily record: %w", err)
	}

	if err := s.dailyRepo.Upsert(ctx, record); err != nil {
		return nil, fmt.Errorf("persisting daily record: %w", err)
	}

	return record, nil
}

// UpdateDailyRecord applies a partial update to the daily record, then returns
// the full updated record.
func (s *DailyServiceImpl) UpdateDailyRecord(ctx context.Context, date string, req *UpdateDailyRecordRequest) (*DailyRecord, error) {
	if err := validateDate(date); err != nil {
		return nil, err
	}
	// Ensure the record exists first (generate if needed).
	existing, err := s.GetDailyRecord(ctx, date)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, fmt.Errorf("daily record not found for date: %s", date)
	}

	// Build flat $set fields from the req.
	setFields := buildSetFields(req)
	if len(setFields) == 0 {
		return existing, nil
	}

	if err := s.dailyRepo.PatchByDate(ctx, date, setFields); err != nil {
		return nil, fmt.Errorf("patching daily record: %w", err)
	}

	return s.dailyRepo.FindByDate(ctx, date)
}

// buildSetFields converts a nested patch map into flat dot-notation keys
// suitable for MongoDB $set. Supported top-level keys: context, tasks, habits, journal.
func buildSetFields(req *UpdateDailyRecordRequest) map[string]interface{} {
	set := make(map[string]interface{})

	if req.Context != nil {
		if req.Context.Mode != nil {
			set["context.mode"] = *req.Context.Mode
		}
		if req.Context.Weather != nil {
			set["context.weather"] = *req.Context.Weather
		}
	}
	if req.Tasks != nil {
		set["tasks"] = req.Tasks
	}
	if req.Habits != nil {
		set["habits"] = req.Habits
	}
	if req.Journal != nil {
		if req.Journal.OneLineReview != nil {
			set["journal.oneLineReview"] = *req.Journal.OneLineReview
		}
		if req.Journal.ThreeLineDiary != nil {
			set["journal.threeLineDiary"] = *req.Journal.ThreeLineDiary
		}
	}

	return set
}

func (s *DailyServiceImpl) generateDailyRecord(ctx context.Context, date string) (*DailyRecord, error) {
	tasks, err := s.taskRepo.FindActiveTasks(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetching active tasks: %w", err)
	}

	habits, err := s.habitRepo.FindActiveHabits(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetching active habits: %w", err)
	}

	taskEntries := make([]TaskEntry, 0, len(tasks))
	for _, t := range tasks {
		taskEntries = append(taskEntries, TaskEntry{
			TaskID:       t.ID,
			TargetAmount: t.Metrics.DailyTarget,
			ActualAmount: 0,
			IsCompleted:  false,
		})
	}

	habitEntries := make([]HabitEntry, 0, len(habits))
	for _, h := range habits {
		habitEntries = append(habitEntries, HabitEntry{
			HabitID:     h.ID,
			IsCompleted: false,
		})
	}

	return &DailyRecord{
		ID:   primitive.NewObjectID(),
		Date: date,
		Context: DayContext{
			Mode:    "Growth",
			Weather: "sunny",
		},
		Tasks:  taskEntries,
		Habits: habitEntries,
		Journal: Journal{
			OneLineReview:  "",
			ThreeLineDiary: "",
		},
	}, nil
}

func (s *DailyServiceImpl) appendMissingEntries(ctx context.Context, record *DailyRecord) (bool, error) {
	updated := false

	tasks, err := s.taskRepo.FindActiveTasks(ctx)
	if err != nil {
		return false, fmt.Errorf("fetching active tasks: %w", err)
	}
	existingTasks := make(map[primitive.ObjectID]bool)
	for _, t := range record.Tasks {
		existingTasks[t.TaskID] = true
	}
	for _, t := range tasks {
		if !existingTasks[t.ID] {
			record.Tasks = append(record.Tasks, TaskEntry{
				TaskID:       t.ID,
				TargetAmount: t.Metrics.DailyTarget,
				ActualAmount: 0,
				IsCompleted:  false,
			})
			updated = true
		}
	}

	habits, err := s.habitRepo.FindActiveHabits(ctx)
	if err != nil {
		return false, fmt.Errorf("fetching active habits: %w", err)
	}
	existingHabits := make(map[primitive.ObjectID]bool)
	for _, h := range record.Habits {
		existingHabits[h.HabitID] = true
	}
	for _, h := range habits {
		if !existingHabits[h.ID] {
			record.Habits = append(record.Habits, HabitEntry{
				HabitID:     h.ID,
				IsCompleted: false,
			})
			updated = true
		}
	}

	return updated, nil
}
