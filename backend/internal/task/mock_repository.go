package task

import (
	"context"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) FindActiveTasks(ctx context.Context) ([]Task, error) {
	args := m.Called(ctx)
	if args.Get(0) != nil {
		return args.Get(0).([]Task), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTaskRepository) FindAll(ctx context.Context) ([]Task, error) {
	args := m.Called(ctx)
	if args.Get(0) != nil {
		return args.Get(0).([]Task), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTaskRepository) FindByID(ctx context.Context, id string) (*Task, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*Task), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTaskRepository) Create(ctx context.Context, task *Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *MockTaskRepository) Update(ctx context.Context, task *Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *MockTaskRepository) MigrateTaskAtomic(ctx context.Context, archivedTask *Task, newTask *Task) error {
	args := m.Called(ctx, archivedTask, newTask)
	return args.Error(0)
}

func (m *MockTaskRepository) EnsureIndexes(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

type MockTaskProgressAggregator struct {
	mock.Mock
}

func (m *MockTaskProgressAggregator) SumTaskProgressByIDs(ctx context.Context, taskIDs []primitive.ObjectID) (map[primitive.ObjectID]int, error) {
	args := m.Called(ctx, taskIDs)
	if args.Get(0) != nil {
		return args.Get(0).(map[primitive.ObjectID]int), args.Error(1)
	}
	return nil, args.Error(1)
}

