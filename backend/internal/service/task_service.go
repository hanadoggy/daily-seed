package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"daily-seed/internal/model"
	"daily-seed/internal/repository"
	"daily-seed/pkg/jst"
)

type TaskService interface {
	List(ctx context.Context) ([]model.Task, error)
	Get(ctx context.Context, id string) (*model.Task, error)
	Create(ctx context.Context, task *model.Task) (*model.Task, error)
	Update(ctx context.Context, task *model.Task) (*model.Task, error)
	Archive(ctx context.Context, id string) error
}

type taskService struct {
	repo repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) TaskService {
	return &taskService{repo: repo}
}

func (s *taskService) List(ctx context.Context) ([]model.Task, error) {
	return s.repo.FindActiveTasks(ctx)
}

func (s *taskService) Get(ctx context.Context, id string) (*model.Task, error) {
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

func validateTask(task *model.Task) error {
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

func (s *taskService) Create(ctx context.Context, task *model.Task) (*model.Task, error) {
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

func (s *taskService) Update(ctx context.Context, task *model.Task) (*model.Task, error) {
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

func (s *taskService) Archive(ctx context.Context, id string) error {
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
