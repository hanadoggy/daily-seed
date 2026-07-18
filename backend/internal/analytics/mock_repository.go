package analytics

import (
	"context"

	"daily-seed/internal/daily"
	"daily-seed/internal/task"

	"github.com/stretchr/testify/mock"
)

type MockDailyRecordRepository struct {
	mock.Mock
}

func (m *MockDailyRecordRepository) FindBetweenDates(ctx context.Context, startDate string, endDate string) ([]*daily.DailyRecord, error) {
	args := m.Called(ctx, startDate, endDate)
	if args.Get(0) != nil {
		return args.Get(0).([]*daily.DailyRecord), args.Error(1)
	}
	return nil, args.Error(1)
}

type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) FindAll(ctx context.Context) ([]task.Task, error) {
	args := m.Called(ctx)
	if args.Get(0) != nil {
		return args.Get(0).([]task.Task), args.Error(1)
	}
	return nil, args.Error(1)
}
