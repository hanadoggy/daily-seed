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

type HabitService interface {
	List(ctx context.Context) ([]model.Habit, error)
	Get(ctx context.Context, id string) (*model.Habit, error)
	Create(ctx context.Context, habit *model.Habit) (*model.Habit, error)
	Update(ctx context.Context, habit *model.Habit) (*model.Habit, error)
	Archive(ctx context.Context, id string) error
}

type habitService struct {
	repo repository.HabitRepository
}

func NewHabitService(repo repository.HabitRepository) HabitService {
	return &habitService{repo: repo}
}

func (s *habitService) List(ctx context.Context) ([]model.Habit, error) {
	return s.repo.FindActiveHabits(ctx)
}

func (s *habitService) Get(ctx context.Context, id string) (*model.Habit, error) {
	habit, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("finding habit: %w", err)
	}
	if habit == nil {
		return nil, fmt.Errorf("habit not found: %s", id)
	}
	return habit, nil
}

func (s *habitService) Create(ctx context.Context, habit *model.Habit) (*model.Habit, error) {
	if strings.TrimSpace(habit.Title) == "" {
		return nil, fmt.Errorf("title is required")
	}
	if strings.TrimSpace(habit.Category) == "" {
		return nil, fmt.Errorf("category is required")
	}

	habit.ID = fmt.Sprintf("habit_%d", time.Now().In(jst.Location()).UnixMilli())
	habit.Status = "active"

	if err := s.repo.Create(ctx, habit); err != nil {
		return nil, fmt.Errorf("creating habit: %w", err)
	}
	return habit, nil
}

func (s *habitService) Update(ctx context.Context, habit *model.Habit) (*model.Habit, error) {
	existing, err := s.repo.FindByID(ctx, habit.ID)
	if err != nil {
		return nil, fmt.Errorf("finding habit: %w", err)
	}
	if existing == nil {
		return nil, fmt.Errorf("habit not found: %s", habit.ID)
	}

	if strings.TrimSpace(habit.Title) == "" {
		return nil, fmt.Errorf("title is required")
	}

	habit.Status = existing.Status

	if err := s.repo.Update(ctx, habit); err != nil {
		return nil, fmt.Errorf("updating habit: %w", err)
	}
	return habit, nil
}

func (s *habitService) Archive(ctx context.Context, id string) error {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("finding habit: %w", err)
	}
	if existing == nil {
		return fmt.Errorf("habit not found: %s", id)
	}

	existing.Status = "archived"
	if err := s.repo.Update(ctx, existing); err != nil {
		return fmt.Errorf("archiving habit: %w", err)
	}
	return nil
}
