package analytics

import (
	"context"
	"fmt"
	"testing"

	"daily-seed/internal/daily"
	"daily-seed/internal/task"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetHeatmapData(t *testing.T) {
	mockDailyRepo := new(MockDailyRecordRepository)
	mockTaskRepo := new(MockTaskRepository)
	svc := NewAnalyticsService(mockDailyRepo, mockTaskRepo)

	ctx := context.Background()
	year := 2026

	taskID1 := primitive.NewObjectID()
	taskID2 := primitive.NewObjectID()

	tasks := []task.Task{
		{ID: taskID1, Section: "dev"},
		{ID: taskID2, Section: "japanese"},
	}

	records := []*daily.DailyRecord{
		{
			Date: "2026-01-01",
			Tasks: []daily.TaskEntry{
				{TaskID: taskID1, IsCompleted: true},
				{TaskID: taskID2, IsCompleted: false},
			},
			Habits: []daily.HabitEntry{
				{HabitID: primitive.NewObjectID(), IsCompleted: true},
			},
		},
		{
			Date: "2026-12-31",
			Tasks: []daily.TaskEntry{
				{TaskID: taskID2, IsCompleted: true},
			},
		},
	}

	startDate := fmt.Sprintf("%04d-01-01", year)
	endDate := fmt.Sprintf("%04d-12-31", year)

	mockDailyRepo.On("FindBetweenDates", ctx, startDate, endDate).Return(records, nil)
	mockTaskRepo.On("FindAll", ctx).Return(tasks, nil)

	res, err := svc.GetHeatmapData(ctx, year)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res.Days, 365) // 2026 is not a leap year

	// Find the day 2026-01-01
	var day1 HeatmapDay
	var dayLast HeatmapDay
	for _, d := range res.Days {
		if d.Date == "2026-01-01" {
			day1 = d
		}
		if d.Date == "2026-12-31" {
			dayLast = d
		}
	}

	assert.Equal(t, "2026-01-01", day1.Date)
	assert.Equal(t, 2, day1.Total) // 1 task + 1 habit
	assert.Equal(t, 1, day1.Habits)
	assert.Equal(t, 1, day1.SectionCounts["dev"])
	assert.Equal(t, 0, day1.SectionCounts["japanese"])

	assert.Equal(t, "2026-12-31", dayLast.Date)
	assert.Equal(t, 1, dayLast.Total)
	assert.Equal(t, 0, dayLast.Habits)
	assert.Equal(t, 0, dayLast.SectionCounts["dev"])
	assert.Equal(t, 1, dayLast.SectionCounts["japanese"])

	mockDailyRepo.AssertExpectations(t)
	mockTaskRepo.AssertExpectations(t)
}

func TestGetHeatmapData_LeapYear(t *testing.T) {
	mockDailyRepo := new(MockDailyRecordRepository)
	mockTaskRepo := new(MockTaskRepository)
	svc := NewAnalyticsService(mockDailyRepo, mockTaskRepo)

	ctx := context.Background()
	year := 2024

	startDate := fmt.Sprintf("%04d-01-01", year)
	endDate := fmt.Sprintf("%04d-12-31", year)

	mockDailyRepo.On("FindBetweenDates", ctx, startDate, endDate).Return([]*daily.DailyRecord{}, nil)
	mockTaskRepo.On("FindAll", ctx).Return([]task.Task{}, nil)

	res, err := svc.GetHeatmapData(ctx, year)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res.Days, 366) // 2024 is a leap year
}

func TestGetHeatmapData_ErrorFetchingDailyRecords(t *testing.T) {
	mockDailyRepo := new(MockDailyRecordRepository)
	mockTaskRepo := new(MockTaskRepository)
	svc := NewAnalyticsService(mockDailyRepo, mockTaskRepo)

	ctx := context.Background()
	year := 2026
	
	startDate := fmt.Sprintf("%04d-01-01", year)
	endDate := fmt.Sprintf("%04d-12-31", year)

	mockDailyRepo.On("FindBetweenDates", ctx, startDate, endDate).Return(nil, fmt.Errorf("db error"))

	_, err := svc.GetHeatmapData(ctx, year)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to fetch daily records")
}

func TestGetHeatmapData_ErrorFetchingTasks(t *testing.T) {
	mockDailyRepo := new(MockDailyRecordRepository)
	mockTaskRepo := new(MockTaskRepository)
	svc := NewAnalyticsService(mockDailyRepo, mockTaskRepo)

	ctx := context.Background()
	year := 2026
	
	startDate := fmt.Sprintf("%04d-01-01", year)
	endDate := fmt.Sprintf("%04d-12-31", year)

	mockDailyRepo.On("FindBetweenDates", ctx, startDate, endDate).Return([]*daily.DailyRecord{}, nil)
	mockTaskRepo.On("FindAll", ctx).Return([]task.Task(nil), fmt.Errorf("db error"))

	_, err := svc.GetHeatmapData(ctx, year)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to fetch tasks")
}

func TestGetHeatmapData_UnmappedTasksAndEmptyRecords(t *testing.T) {
	mockDailyRepo := new(MockDailyRecordRepository)
	mockTaskRepo := new(MockTaskRepository)
	svc := NewAnalyticsService(mockDailyRepo, mockTaskRepo)

	ctx := context.Background()
	year := 2026

	unmappedTaskID := primitive.NewObjectID()
	records := []*daily.DailyRecord{
		{
			Date: "2026-05-15",
			Tasks: []daily.TaskEntry{
				{TaskID: unmappedTaskID, IsCompleted: true},
			},
			Habits: nil, // null habits slice
		},
	}

	startDate := fmt.Sprintf("%04d-01-01", year)
	endDate := fmt.Sprintf("%04d-12-31", year)

	mockDailyRepo.On("FindBetweenDates", ctx, startDate, endDate).Return(records, nil)
	mockTaskRepo.On("FindAll", ctx).Return([]task.Task{}, nil) // No tasks mapped

	res, err := svc.GetHeatmapData(ctx, year)

	assert.NoError(t, err)
	assert.NotNil(t, res)

	var targetDay HeatmapDay
	for _, d := range res.Days {
		if d.Date == "2026-05-15" {
			targetDay = d
			break
		}
	}

	assert.Equal(t, "2026-05-15", targetDay.Date)
	assert.Equal(t, 1, targetDay.Total) // Completed task counted in Total
	assert.Equal(t, 0, targetDay.Habits)
	assert.Empty(t, targetDay.SectionCounts) // Unmapped task section count is not updated
}

