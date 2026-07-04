package habit

import (
	"context"
	"fmt"
	"strings"
	"time"

	"daily-seed/pkg/jst"
)

type HabitServiceImpl struct {
	repo HabitRepository
}

func NewHabitService(repo HabitRepository) HabitService {
	return &HabitServiceImpl{repo: repo}
}

func (s *HabitServiceImpl) List(ctx context.Context) ([]Habit, error) {
	return s.repo.FindAll(ctx)
}

func (s *HabitServiceImpl) Get(ctx context.Context, id string) (*Habit, error) {
	habit, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("finding habit: %w", err)
	}
	if habit == nil {
		return nil, fmt.Errorf("habit not found: %s", id)
	}
	return habit, nil
}

func (s *HabitServiceImpl) Create(ctx context.Context, habit *Habit) (*Habit, error) {
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

func (s *HabitServiceImpl) Update(ctx context.Context, habit *Habit) (*Habit, error) {
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
	if strings.TrimSpace(habit.Category) == "" {
		return nil, fmt.Errorf("category is required")
	}

	habit.Status = existing.Status

	if err := s.repo.Update(ctx, habit); err != nil {
		return nil, fmt.Errorf("updating habit: %w", err)
	}
	return habit, nil
}

func (s *HabitServiceImpl) Archive(ctx context.Context, id string) error {
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
