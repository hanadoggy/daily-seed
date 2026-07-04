package task

import (
	"context"
	"fmt"
	"strings"
	"time"

	"daily-seed/pkg/jst"
)

type TaskServiceImpl struct {
	repo       TaskRepository
	aggregator TaskProgressAggregator
}

func NewTaskService(repo TaskRepository, aggregator TaskProgressAggregator) TaskService {
	return &TaskServiceImpl{repo: repo, aggregator: aggregator}
}

func (s *TaskServiceImpl) List(ctx context.Context) ([]Task, error) {
	return s.repo.FindActiveTasks(ctx)
}

func (s *TaskServiceImpl) Get(ctx context.Context, id string) (*Task, error) {
	task, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("finding task: %w", err)
	}
	if task == nil {
		return nil, fmt.Errorf("task not found: %s", id)
	}
	return task, nil
}

var validSections = map[string]bool{
	"japanese": true,
	"dev":      true,
	"self_dev": true,
}

var validTaskTypes = map[string]bool{
	"quantitative": true,
	"boolean":      true,
}

func validateTask(task *Task) error {
	if strings.TrimSpace(task.Title) == "" {
		return fmt.Errorf("title is required")
	}
	if !validSections[task.Section] {
		return fmt.Errorf("section must be one of: japanese, dev, self_dev")
	}
	if !validTaskTypes[task.Type] {
		return fmt.Errorf("type must be one of: quantitative, boolean")
	}
	if task.Type == "quantitative" && task.Metrics.DailyTarget <= 0 {
		return fmt.Errorf("dailyTarget must be positive for quantitative tasks")
	}
	return nil
}

func (s *TaskServiceImpl) Create(ctx context.Context, task *Task) (*Task, error) {
	if err := validateTask(task); err != nil {
		return nil, err
	}

	task.ID = fmt.Sprintf("task_%d", time.Now().In(jst.Location()).UnixMilli())
	task.Status = "active"

	// Default conditions if not provided.
	if task.Conditions.Weather == "" {
		task.Conditions.Weather = "any"
	}
	if task.Conditions.Mode == "" {
		task.Conditions.Mode = "any"
	}
	// Boolean tasks always have dailyTarget = 1.
	if task.Type == "boolean" {
		task.Metrics.DailyTarget = 1
	}

	if err := s.repo.Create(ctx, task); err != nil {
		return nil, fmt.Errorf("creating task: %w", err)
	}
	return task, nil
}

func (s *TaskServiceImpl) Update(ctx context.Context, task *Task) (*Task, error) {
	existing, err := s.repo.FindByID(ctx, task.ID)
	if err != nil {
		return nil, fmt.Errorf("finding task: %w", err)
	}
	if existing == nil {
		return nil, fmt.Errorf("task not found: %s", task.ID)
	}

	if err := validateTask(task); err != nil {
		return nil, err
	}

	task.Status = existing.Status

	if err := s.repo.Update(ctx, task); err != nil {
		return nil, fmt.Errorf("updating task: %w", err)
	}
	return task, nil
}

func (s *TaskServiceImpl) Archive(ctx context.Context, id string) error {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("finding task: %w", err)
	}
	if existing == nil {
		return fmt.Errorf("task not found: %s", id)
	}

	existing.Status = "archived"
	if err := s.repo.Update(ctx, existing); err != nil {
		return fmt.Errorf("archiving task: %w", err)
	}
	return nil
}

func (s *TaskServiceImpl) GetProgressForActiveTasks(ctx context.Context) ([]TaskProgress, error) {
	tasks, err := s.repo.FindActiveTasks(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetching active tasks: %w", err)
	}

	progress := make([]TaskProgress, 0, len(tasks))
	for _, t := range tasks {
		// Skip boolean tasks — they don't have meaningful cumulative progress.
		if t.Type != "quantitative" || t.Metrics.TotalTarget <= 0 {
			continue
		}

		completed, err := s.aggregator.SumTaskProgress(ctx, t.ID)
		if err != nil {
			return nil, fmt.Errorf("summing progress for task %s: %w", t.ID, err)
		}

		pct := 0.0
		if t.Metrics.TotalTarget > 0 {
			pct = float64(completed) / float64(t.Metrics.TotalTarget) * 100
		}

		progress = append(progress, TaskProgress{
			TaskID:         t.ID,
			Title:          t.Title,
			TotalTarget:    t.Metrics.TotalTarget,
			TotalCompleted: completed,
			Percentage:     pct,
		})
	}

	return progress, nil
}

func (s *TaskServiceImpl) MigrateTask(ctx context.Context, id string) (*MigrationResult, error) {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("finding task: %w", err)
	}
	if existing == nil {
		return nil, fmt.Errorf("task not found: %s", id)
	}
	if existing.Status != "active" {
		return nil, fmt.Errorf("cannot migrate non-active task: %s", id)
	}

	// Verify cumulative progress meets or exceeds totalTarget.
	if existing.Metrics.TotalTarget > 0 {
		completed, err := s.aggregator.SumTaskProgress(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("checking task progress: %w", err)
		}
		if completed < existing.Metrics.TotalTarget {
			return nil, fmt.Errorf("task progress (%d/%d) has not reached the target", completed, existing.Metrics.TotalTarget)
		}
	}

	// Archive the old task.
	existing.Status = "archived"
	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, fmt.Errorf("archiving task for migration: %w", err)
	}

	// Create new task with the same configuration but reset progress.
	newTask := &Task{
		ID:      fmt.Sprintf("task_%d", time.Now().In(jst.Location()).UnixMilli()),
		Section: existing.Section,
		Title:   existing.Title,
		Type:    existing.Type,
		Metrics: TaskMetrics{
			DailyTarget: existing.Metrics.DailyTarget,
			TotalTarget: existing.Metrics.TotalTarget,
		},
		Conditions: existing.Conditions,
		Status:     "active",
	}

	if err := s.repo.Create(ctx, newTask); err != nil {
		return nil, fmt.Errorf("creating migrated task: %w", err)
	}

	return &MigrationResult{
		ArchivedTask: *existing,
		NewTask:      *newTask,
	}, nil
}
