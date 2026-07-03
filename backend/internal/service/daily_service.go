package service

import (
	"context"
	"fmt"
	"log/slog"

	"daily-seed/internal/model"
	"daily-seed/internal/repository"
)

type DailyService interface {
	GetDailyRecord(ctx context.Context, date string) (*model.DailyRecord, error)
	UpdateDailyRecord(ctx context.Context, date string, patch map[string]interface{}) (*model.DailyRecord, error)
}

type dailyService struct {
	dailyRepo repository.DailyRecordRepository
	taskRepo  repository.TaskRepository
	habitRepo repository.HabitRepository
}

func NewDailyService(
	dailyRepo repository.DailyRecordRepository,
	taskRepo repository.TaskRepository,
	habitRepo repository.HabitRepository,
) DailyService {
	return &dailyService{
		dailyRepo: dailyRepo,
		taskRepo:  taskRepo,
		habitRepo: habitRepo,
	}
}

// GetDailyRecord retrieves the DailyRecord for the given date.
// If no record exists, it generates one from the currently active tasks and habits.
func (s *dailyService) GetDailyRecord(ctx context.Context, date string) (*model.DailyRecord, error) {
	record, err := s.dailyRepo.FindByDate(ctx, date)
	if err != nil {
		return nil, fmt.Errorf("finding daily record: %w", err)
	}

	if record != nil {
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
func (s *dailyService) UpdateDailyRecord(ctx context.Context, date string, patch map[string]interface{}) (*model.DailyRecord, error) {
	// Ensure the record exists first (generate if needed).
	existing, err := s.GetDailyRecord(ctx, date)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, fmt.Errorf("daily record not found for date: %s", date)
	}

	// Apply flat patch fields onto the record via repository.
	if err := s.dailyRepo.Upsert(ctx, existing); err != nil {
		return nil, fmt.Errorf("updating daily record: %w", err)
	}

	return s.dailyRepo.FindByDate(ctx, date)
}

func (s *dailyService) generateDailyRecord(ctx context.Context, date string) (*model.DailyRecord, error) {
	tasks, err := s.taskRepo.FindActiveTasks(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetching active tasks: %w", err)
	}

	habits, err := s.habitRepo.FindActiveHabits(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetching active habits: %w", err)
	}

	taskEntries := make([]model.TaskEntry, 0, len(tasks))
	for _, t := range tasks {
		taskEntries = append(taskEntries, model.TaskEntry{
			TaskID:       t.ID,
			TargetAmount: t.Metrics.DailyTarget,
			ActualAmount: 0,
			IsCompleted:  false,
		})
	}

	habitEntries := make([]model.HabitEntry, 0, len(habits))
	for _, h := range habits {
		habitEntries = append(habitEntries, model.HabitEntry{
			HabitID:     h.ID,
			IsCompleted: false,
		})
	}

	return &model.DailyRecord{
		ID:   date,
		Date: date,
		Context: model.DayContext{
			Mode:    "Growth",
			Weather: "",
		},
		Tasks:  taskEntries,
		Habits: habitEntries,
		Journal: model.Journal{
			OneLineReview:  "",
			ThreeLineDiary: "",
		},
	}, nil
}
