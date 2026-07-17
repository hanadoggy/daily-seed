package daily_test

import (
	"context"
	"errors"
	"testing"

	"daily-seed/internal/daily"
	"daily-seed/internal/habit"
	"daily-seed/internal/task"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func testOID(hex string) primitive.ObjectID {
	id, _ := primitive.ObjectIDFromHex(hex)
	return id
}

func TestDailyService_GetDailyRecord(t *testing.T) {
	ctx := context.Background()
	date := "2023-10-10"
	existingRecord := &daily.DailyRecord{ID: testOID("000000000000000000000001"), Date: date}

	tests := []struct {
		name        string
		date        string
		mockSetup   func(*daily.MockDailyRecordRepository, *task.MockTaskRepository, *habit.MockHabitRepository)
		expectError bool
		errContains string
		validate    func(*testing.T, *daily.DailyRecord)
	}{
		{
			name: "Pass: existing record",
			date: date,
			mockSetup: func(dRepo *daily.MockDailyRecordRepository, tRepo *task.MockTaskRepository, hRepo *habit.MockHabitRepository) {
				dRepo.On("FindByDate", ctx, date).Return(existingRecord, nil)
				tRepo.On("FindActiveTasks", ctx).Return([]task.Task{}, nil)
				hRepo.On("FindActiveHabits", ctx).Return([]habit.Habit{}, nil)
			},
			expectError: false,
			validate: func(t *testing.T, record *daily.DailyRecord) {
				assert.Equal(t, existingRecord, record)
			},
		},
		{
			name: "Fail: invalid date format",
			date: "invalid-date",
			mockSetup: func(dRepo *daily.MockDailyRecordRepository, tRepo *task.MockTaskRepository, hRepo *habit.MockHabitRepository) {},
			expectError: true,
			errContains: "invalid date format",
		},
		{
			name: "Fail: FindByDate error",
			date: date,
			mockSetup: func(dRepo *daily.MockDailyRecordRepository, tRepo *task.MockTaskRepository, hRepo *habit.MockHabitRepository) {
				dRepo.On("FindByDate", ctx, date).Return(nil, errors.New("db error"))
			},
			expectError: true,
			errContains: "finding daily record",
		},
		{
			name: "Pass: generate new record (no tasks/habits)",
			date: date,
			mockSetup: func(dRepo *daily.MockDailyRecordRepository, tRepo *task.MockTaskRepository, hRepo *habit.MockHabitRepository) {
				dRepo.On("FindByDate", ctx, date).Return(nil, nil)
				tRepo.On("FindActiveTasks", ctx).Return([]task.Task{}, nil)
				hRepo.On("FindActiveHabits", ctx).Return([]habit.Habit{}, nil)
				dRepo.On("Upsert", ctx, mock.AnythingOfType("*daily.DailyRecord")).Return(nil)
			},
			expectError: false,
			validate: func(t *testing.T, record *daily.DailyRecord) {
				assert.Equal(t, date, record.Date)
				assert.Empty(t, record.Tasks)
				assert.Empty(t, record.Habits)
			},
		},
		{
			name: "Pass: generate new record (with tasks/habits)",
			date: date,
			mockSetup: func(dRepo *daily.MockDailyRecordRepository, tRepo *task.MockTaskRepository, hRepo *habit.MockHabitRepository) {
				dRepo.On("FindByDate", ctx, date).Return(nil, nil)
				tRepo.On("FindActiveTasks", ctx).Return([]task.Task{{ID: testOID("000000000000000000000001"), Metrics: task.TaskMetrics{DailyTarget: 2}}}, nil)
				hRepo.On("FindActiveHabits", ctx).Return([]habit.Habit{{ID: testOID("000000000000000000000002")}}, nil)
				dRepo.On("Upsert", ctx, mock.AnythingOfType("*daily.DailyRecord")).Return(nil)
			},
			expectError: false,
			validate: func(t *testing.T, record *daily.DailyRecord) {
				assert.Equal(t, date, record.Date)
				assert.Len(t, record.Tasks, 1)
				assert.Equal(t, testOID("000000000000000000000001"), record.Tasks[0].TaskID)
				assert.Len(t, record.Habits, 1)
				assert.Equal(t, testOID("000000000000000000000002"), record.Habits[0].HabitID)
			},
		},
		{
			name: "Fail: task fetch error during generation",
			date: date,
			mockSetup: func(dRepo *daily.MockDailyRecordRepository, tRepo *task.MockTaskRepository, hRepo *habit.MockHabitRepository) {
				dRepo.On("FindByDate", ctx, date).Return(nil, nil)
				tRepo.On("FindActiveTasks", ctx).Return(nil, errors.New("db error"))
			},
			expectError: true,
			errContains: "generating daily record",
		},
		{
			name: "Fail: habit fetch error during generation",
			date: date,
			mockSetup: func(dRepo *daily.MockDailyRecordRepository, tRepo *task.MockTaskRepository, hRepo *habit.MockHabitRepository) {
				dRepo.On("FindByDate", ctx, date).Return(nil, nil)
				tRepo.On("FindActiveTasks", ctx).Return([]task.Task{}, nil)
				hRepo.On("FindActiveHabits", ctx).Return(nil, errors.New("db error"))
			},
			expectError: true,
			errContains: "generating daily record",
		},
		{
			name: "Fail: Upsert error during generation",
			date: date,
			mockSetup: func(dRepo *daily.MockDailyRecordRepository, tRepo *task.MockTaskRepository, hRepo *habit.MockHabitRepository) {
				dRepo.On("FindByDate", ctx, date).Return(nil, nil)
				tRepo.On("FindActiveTasks", ctx).Return([]task.Task{}, nil)
				hRepo.On("FindActiveHabits", ctx).Return([]habit.Habit{}, nil)
				dRepo.On("Upsert", ctx, mock.Anything).Return(errors.New("db error"))
			},
			expectError: true,
			errContains: "persisting daily record",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDailyRepo := new(daily.MockDailyRecordRepository)
			mockTaskRepo := new(task.MockTaskRepository)
			mockHabitRepo := new(habit.MockHabitRepository)

			tt.mockSetup(mockDailyRepo, mockTaskRepo, mockHabitRepo)
			svc := daily.NewDailyService(mockDailyRepo, mockTaskRepo, mockHabitRepo)

			record, err := svc.GetDailyRecord(ctx, tt.date)
			if tt.expectError {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, record)
			} else {
				assert.NoError(t, err)
				if tt.validate != nil {
					tt.validate(t, record)
				}
			}
			mockDailyRepo.AssertExpectations(t)
			mockTaskRepo.AssertExpectations(t)
			mockHabitRepo.AssertExpectations(t)
		})
	}
}

func TestDailyService_UpdateDailyRecord(t *testing.T) {
	ctx := context.Background()
	date := "2023-10-10"
	existingRecord := &daily.DailyRecord{ID: testOID("000000000000000000000001"), Date: date}

	tests := []struct {
		name        string
		date        string
		patch       *daily.UpdateDailyRecordRequest
		mockSetup   func(*daily.MockDailyRecordRepository, *task.MockTaskRepository, *habit.MockHabitRepository)
		expectError bool
		errContains string
	}{
		{
			name:  "Pass: successful update with fields",
			date:  date,
			patch: &daily.UpdateDailyRecordRequest{Context: &daily.DayContextPatch{Mode: (*daily.ContextMode)(func() *string { s := "Growth"; return &s }())}},
			mockSetup: func(dRepo *daily.MockDailyRecordRepository, tRepo *task.MockTaskRepository, hRepo *habit.MockHabitRepository) {
				dRepo.On("FindByDate", ctx, date).Return(existingRecord, nil).Once() // First GetDailyRecord call
				tRepo.On("FindActiveTasks", ctx).Return([]task.Task{}, nil)
				hRepo.On("FindActiveHabits", ctx).Return([]habit.Habit{}, nil)
				dRepo.On("PatchByDate", ctx, date, mock.Anything).Return(nil)
				dRepo.On("FindByDate", ctx, date).Return(existingRecord, nil).Once() // Final FindByDate
			},
			expectError: false,
		},
		{
			name:  "Pass: empty patch does nothing",
			date:  date,
			patch: &daily.UpdateDailyRecordRequest{},
			mockSetup: func(dRepo *daily.MockDailyRecordRepository, tRepo *task.MockTaskRepository, hRepo *habit.MockHabitRepository) {
				dRepo.On("FindByDate", ctx, date).Return(existingRecord, nil).Once() // Only GetDailyRecord called
				tRepo.On("FindActiveTasks", ctx).Return([]task.Task{}, nil)
				hRepo.On("FindActiveHabits", ctx).Return([]habit.Habit{}, nil)
			},
			expectError: false,
		},
		{
			name:  "Fail: GetDailyRecord error",
			date:  date,
			patch: &daily.UpdateDailyRecordRequest{Context: &daily.DayContextPatch{Mode: (*daily.ContextMode)(func() *string { s := "Growth"; return &s }())}},
			mockSetup: func(dRepo *daily.MockDailyRecordRepository, tRepo *task.MockTaskRepository, hRepo *habit.MockHabitRepository) {
				dRepo.On("FindByDate", ctx, date).Return(nil, errors.New("db error")).Once()
			},
			expectError: true,
			errContains: "finding daily record",
		},
		{
			name:  "Fail: PatchByDate error",
			date:  date,
			patch: &daily.UpdateDailyRecordRequest{Context: &daily.DayContextPatch{Mode: (*daily.ContextMode)(func() *string { s := "Growth"; return &s }())}},
			mockSetup: func(dRepo *daily.MockDailyRecordRepository, tRepo *task.MockTaskRepository, hRepo *habit.MockHabitRepository) {
				dRepo.On("FindByDate", ctx, date).Return(existingRecord, nil).Once()
				tRepo.On("FindActiveTasks", ctx).Return([]task.Task{}, nil)
				hRepo.On("FindActiveHabits", ctx).Return([]habit.Habit{}, nil)
				dRepo.On("PatchByDate", ctx, date, mock.Anything).Return(errors.New("db error"))
			},
			expectError: true,
			errContains: "patching daily record",
		},
		{
			name:  "Fail: invalid date format",
			date:  "invalid-date",
			patch: &daily.UpdateDailyRecordRequest{Context: &daily.DayContextPatch{Mode: (*daily.ContextMode)(func() *string { s := "Growth"; return &s }())}},
			mockSetup: func(dRepo *daily.MockDailyRecordRepository, tRepo *task.MockTaskRepository, hRepo *habit.MockHabitRepository) {},
			expectError: true,
			errContains: "invalid date format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDailyRepo := new(daily.MockDailyRecordRepository)
			mockTaskRepo := new(task.MockTaskRepository)
			mockHabitRepo := new(habit.MockHabitRepository)

			tt.mockSetup(mockDailyRepo, mockTaskRepo, mockHabitRepo)
			svc := daily.NewDailyService(mockDailyRepo, mockTaskRepo, mockHabitRepo)

			record, err := svc.UpdateDailyRecord(ctx, tt.date, tt.patch)
			if tt.expectError {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, record)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, record)
			}
			mockDailyRepo.AssertExpectations(t)
			mockTaskRepo.AssertExpectations(t)
			mockHabitRepo.AssertExpectations(t)
		})
	}
}
