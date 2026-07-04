package service_test

import (
	"context"
	"testing"

	"daily-seed/internal/model"
	"daily-seed/internal/repository/mocks"
	"daily-seed/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDailyService_GetDailyRecord(t *testing.T) {
	ctx := context.Background()
	date := "2023-10-10"

	t.Run("existing_record", func(t *testing.T) {
		mockDailyRepo := new(mocks.MockDailyRecordRepository)
		mockTaskRepo := new(mocks.MockTaskRepository)
		mockHabitRepo := new(mocks.MockHabitRepository)

		svc := service.NewDailyService(mockDailyRepo, mockTaskRepo, mockHabitRepo)

		expected := &model.DailyRecord{ID: date, Date: date}
		mockDailyRepo.On("FindByDate", ctx, date).Return(expected, nil)

		record, err := svc.GetDailyRecord(ctx, date)
		assert.NoError(t, err)
		assert.Equal(t, expected, record)
		mockDailyRepo.AssertExpectations(t)
	})

	t.Run("generate_new_record", func(t *testing.T) {
		mockDailyRepo := new(mocks.MockDailyRecordRepository)
		mockTaskRepo := new(mocks.MockTaskRepository)
		mockHabitRepo := new(mocks.MockHabitRepository)

		svc := service.NewDailyService(mockDailyRepo, mockTaskRepo, mockHabitRepo)

		// FindByDate returns nil to trigger generation
		mockDailyRepo.On("FindByDate", ctx, date).Return((*model.DailyRecord)(nil), nil)

		// Mock fetching tasks and habits
		tasks := []model.Task{{ID: "t1", Metrics: model.TaskMetrics{DailyTarget: 2}}}
		mockTaskRepo.On("FindActiveTasks", ctx).Return(tasks, nil)

		habits := []model.Habit{{ID: "h1"}}
		mockHabitRepo.On("FindActiveHabits", ctx).Return(habits, nil)

		// Mock Upsert
		mockDailyRepo.On("Upsert", ctx, mock.AnythingOfType("*model.DailyRecord")).Return(nil)

		record, err := svc.GetDailyRecord(ctx, date)
		assert.NoError(t, err)
		assert.NotNil(t, record)
		assert.Equal(t, date, record.Date)
		assert.Len(t, record.Tasks, 1)
		assert.Equal(t, "t1", record.Tasks[0].TaskID)
		assert.Len(t, record.Habits, 1)

		mockDailyRepo.AssertExpectations(t)
		mockTaskRepo.AssertExpectations(t)
		mockHabitRepo.AssertExpectations(t)
	})
}
