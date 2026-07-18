package task

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockTaskService struct {
	mock.Mock
}

func (m *MockTaskService) List(ctx context.Context) ([]Task, error) {
	args := m.Called(ctx)
	if args.Get(0) != nil {
		return args.Get(0).([]Task), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTaskService) Get(ctx context.Context, id string) (*Task, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*Task), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTaskService) Create(ctx context.Context, task *Task) (*Task, error) {
	args := m.Called(ctx, task)
	if args.Get(0) != nil {
		return args.Get(0).(*Task), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTaskService) Update(ctx context.Context, task *Task) (*Task, error) {
	args := m.Called(ctx, task)
	if args.Get(0) != nil {
		return args.Get(0).(*Task), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTaskService) Archive(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockTaskService) GetProgressForActiveTasks(ctx context.Context) ([]TaskProgress, error) {
	args := m.Called(ctx)
	if args.Get(0) != nil {
		return args.Get(0).([]TaskProgress), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTaskService) MigrateTask(ctx context.Context, id string, req MigrateTaskRequest) (*MigrationResult, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) != nil {
		return args.Get(0).(*MigrationResult), args.Error(1)
	}
	return nil, args.Error(1)
}
