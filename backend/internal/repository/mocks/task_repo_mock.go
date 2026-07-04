package mocks

import (
	"context"

	"daily-seed/internal/model"

	"github.com/stretchr/testify/mock"
)

type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) FindActiveTasks(ctx context.Context) ([]model.Task, error) {
	args := m.Called(ctx)
	if args.Get(0) != nil {
		return args.Get(0).([]model.Task), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTaskRepository) FindAll(ctx context.Context) ([]model.Task, error) {
	args := m.Called(ctx)
	if args.Get(0) != nil {
		return args.Get(0).([]model.Task), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTaskRepository) FindByID(ctx context.Context, id string) (*model.Task, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*model.Task), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTaskRepository) Create(ctx context.Context, task *model.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *MockTaskRepository) Update(ctx context.Context, task *model.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *MockTaskRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
