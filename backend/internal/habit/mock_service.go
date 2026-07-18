package habit

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockHabitService struct {
	mock.Mock
}

func (m *MockHabitService) List(ctx context.Context) ([]Habit, error) {
	args := m.Called(ctx)
	if args.Get(0) != nil {
		return args.Get(0).([]Habit), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockHabitService) Get(ctx context.Context, id string) (*Habit, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*Habit), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockHabitService) Create(ctx context.Context, habit *Habit) (*Habit, error) {
	args := m.Called(ctx, habit)
	if args.Get(0) != nil {
		return args.Get(0).(*Habit), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockHabitService) Update(ctx context.Context, habit *Habit) (*Habit, error) {
	args := m.Called(ctx, habit)
	if args.Get(0) != nil {
		return args.Get(0).(*Habit), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockHabitService) Archive(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
