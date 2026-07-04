package habit

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockHabitRepository struct {
	mock.Mock
}

func (m *MockHabitRepository) FindActiveHabits(ctx context.Context) ([]Habit, error) {
	args := m.Called(ctx)
	if args.Get(0) != nil {
		return args.Get(0).([]Habit), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockHabitRepository) FindAll(ctx context.Context) ([]Habit, error) {
	args := m.Called(ctx)
	if args.Get(0) != nil {
		return args.Get(0).([]Habit), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockHabitRepository) FindByID(ctx context.Context, id string) (*Habit, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*Habit), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockHabitRepository) Create(ctx context.Context, habit *Habit) error {
	args := m.Called(ctx, habit)
	return args.Error(0)
}

func (m *MockHabitRepository) Update(ctx context.Context, habit *Habit) error {
	args := m.Called(ctx, habit)
	return args.Error(0)
}

func (m *MockHabitRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
