package task

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DailyCleaner interface {
	RemoveTaskFromRecordsBeforeDate(ctx context.Context, taskID primitive.ObjectID, date string) error
}

type TaskServiceImpl struct {
	repo       TaskRepository
	aggregator TaskProgressAggregator
	cleaner    DailyCleaner
}

func NewTaskService(repo TaskRepository, aggregator TaskProgressAggregator, cleaner DailyCleaner) TaskService {
	return &TaskServiceImpl{repo: repo, aggregator: aggregator, cleaner: cleaner}
}

func (s *TaskServiceImpl) List(ctx context.Context) ([]Task, error) {
	return s.repo.FindAll(ctx)
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
	if task.Metrics.TotalTarget < 0 {
		return fmt.Errorf("totalTarget cannot be negative")
	}
	if task.StartDate == "" {
		return fmt.Errorf("startDate is required")
	}
	return nil
}

func (s *TaskServiceImpl) Create(ctx context.Context, task *Task) (*Task, error) {
	if err := validateTask(task); err != nil {
		return nil, err
	}

	task.ID = primitive.NewObjectID()
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
	existing, err := s.repo.FindByID(ctx, task.ID.Hex())
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

	// If StartDate is delayed, we need to remove the task from past DailyRecords
	if task.StartDate > existing.StartDate && s.cleaner != nil {
		if err := s.cleaner.RemoveTaskFromRecordsBeforeDate(ctx, task.ID, task.StartDate); err != nil {
			return nil, fmt.Errorf("cleaning past daily records: %w", err)
		}
	}

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
	existing.EndDate = time.Now().Format("2006-01-02")
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

	taskIDs := make([]primitive.ObjectID, 0, len(tasks))
	for _, t := range tasks {
		if t.Type == "quantitative" && t.Metrics.TotalTarget > 0 {
			taskIDs = append(taskIDs, t.ID)
		}
	}

	progressMap, err := s.aggregator.SumTaskProgressByIDs(ctx, taskIDs)
	if err != nil {
		return nil, fmt.Errorf("fetching cumulative progress: %w", err)
	}

	progress := make([]TaskProgress, 0, len(taskIDs))
	for _, t := range tasks {
		if t.Type != "quantitative" || t.Metrics.TotalTarget <= 0 {
			continue
		}

		completed := progressMap[t.ID]

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

func (s *TaskServiceImpl) MigrateTask(ctx context.Context, id string, req MigrateTaskRequest) (*MigrationResult, error) {
	if req.CompletionDate == "" {
		return nil, fmt.Errorf("completionDate is required")
	}
	parsedDate, err := time.Parse("2006-01-02", req.CompletionDate)
	if err != nil {
		return nil, fmt.Errorf("invalid completionDate format")
	}

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
		progressMap, err := s.aggregator.SumTaskProgressByIDs(ctx, []primitive.ObjectID{existing.ID})
		if err != nil {
			return nil, fmt.Errorf("checking task progress: %w", err)
		}
		completed := progressMap[existing.ID]
		if completed < existing.Metrics.TotalTarget {
			return nil, fmt.Errorf("task progress (%d/%d) has not reached the target", completed, existing.Metrics.TotalTarget)
		}
	}

	// Archive the old task in memory so it can be passed to MigrateTaskAtomic
	existing.Status = "archived"
	existing.EndDate = req.CompletionDate

	// Create new task with the same configuration but reset progress.
	newTask := &Task{
		ID:      primitive.NewObjectID(),
		Section: existing.Section,
		Title:   existing.Title,
		Type:    existing.Type,
		Metrics: TaskMetrics{
			DailyTarget: existing.Metrics.DailyTarget,
			TotalTarget: existing.Metrics.TotalTarget,
		},
		Conditions: existing.Conditions,
		Status:     "active",
		StartDate:  parsedDate.AddDate(0, 0, 1).Format("2006-01-02"),
	}

	if err := s.repo.MigrateTaskAtomic(ctx, existing, newTask); err != nil {
		return nil, fmt.Errorf("atomic migration failed: %w", err)
	}

	return &MigrationResult{
		ArchivedTask: *existing,
		NewTask:      *newTask,
	}, nil
}
