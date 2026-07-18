package task_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"daily-seed/internal/task"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func testOID(hex string) primitive.ObjectID {
	id, _ := primitive.ObjectIDFromHex(hex)
	return id
}

var (
	taskOID1 = testOID("000000000000000000000001")
	taskOID2 = testOID("000000000000000000000002")
	taskOID3 = testOID("000000000000000000000003")
)

func TestTaskService_Create(t *testing.T) {
	tests := []struct {
		name        string
		task        *task.Task
		mockSetup   func(*task.MockTaskRepository)
		expectError bool
		errContains string
	}{
		{
			name: "Pass: successful boolean task",
			task: &task.Task{Title: "Learn Go", Section: "dev", Type: "boolean", StartDate: "2026-07-17"},
			mockSetup: func(m *task.MockTaskRepository) {
				m.On("Create", mock.Anything, mock.AnythingOfType("*task.Task")).Return(nil)
			},
			expectError: false,
		},
		{
			name: "Pass: successful quantitative task",
			task: &task.Task{Title: "Read", Section: "self_dev", Type: "quantitative", StartDate: "2026-07-17", Metrics: task.TaskMetrics{DailyTarget: 10}},
			mockSetup: func(m *task.MockTaskRepository) {
				m.On("Create", mock.Anything, mock.AnythingOfType("*task.Task")).Return(nil)
			},
			expectError: false,
		},
		{
			name: "Fail: missing title",
			task: &task.Task{Title: "", Section: "dev", Type: "boolean", StartDate: "2026-07-17"},
			mockSetup: func(m *task.MockTaskRepository) {},
			expectError: true,
			errContains: "title is required",
		},
		{
			name: "Fail: invalid section",
			task: &task.Task{Title: "Test", Section: "unknown", Type: "boolean", StartDate: "2026-07-17"},
			mockSetup: func(m *task.MockTaskRepository) {},
			expectError: true,
			errContains: "section must be one of",
		},
		{
			name: "Fail: invalid type",
			task: &task.Task{Title: "Test", Section: "dev", Type: "unknown", StartDate: "2026-07-17"},
			mockSetup: func(m *task.MockTaskRepository) {},
			expectError: true,
			errContains: "type must be one of",
		},
		{
			name: "Fail: zero target for quantitative",
			task: &task.Task{Title: "Test", Section: "dev", Type: "quantitative", StartDate: "2026-07-17", Metrics: task.TaskMetrics{DailyTarget: 0}},
			mockSetup: func(m *task.MockTaskRepository) {},
			expectError: true,
			errContains: "dailyTarget must be positive",
		},
		{
			name: "Fail: negative total target",
			task: &task.Task{Title: "Test", Section: "dev", Type: "quantitative", StartDate: "2026-07-17", Metrics: task.TaskMetrics{DailyTarget: 1, TotalTarget: -10}},
			mockSetup: func(m *task.MockTaskRepository) {},
			expectError: true,
			errContains: "totalTarget cannot be negative",
		},
		{
			name: "Fail: repo error",
			task: &task.Task{Title: "Test", Section: "dev", Type: "boolean", StartDate: "2026-07-17"},
			mockSetup: func(m *task.MockTaskRepository) {
				m.On("Create", mock.Anything, mock.Anything).Return(errors.New("db error"))
			},
			expectError: true,
			errContains: "creating task",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(task.MockTaskRepository)
			tt.mockSetup(repo)
			svc := task.NewTaskService(repo, new(task.MockTaskProgressAggregator), nil)

			created, err := svc.Create(context.Background(), tt.task)
			if tt.expectError {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.True(t, strings.Contains(err.Error(), tt.errContains))
				}
				assert.Nil(t, created)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, "active", created.Status)
				if tt.task.Type == "boolean" {
					assert.Equal(t, 1, created.Metrics.DailyTarget)
				}
			}
			repo.AssertExpectations(t)
		})
	}
}

func TestTaskService_Update(t *testing.T) {
	existing := &task.Task{ID: taskOID1, Status: "active", StartDate: "2026-07-17", Type: "boolean"}

	tests := []struct {
		name        string
		task        *task.Task
		mockSetup   func(*task.MockTaskRepository)
		expectError bool
		errContains string
	}{
		{
			name: "Pass: successful update",
			task: &task.Task{ID: taskOID1, Title: "Read updated", Section: "dev", Type: "boolean", StartDate: "2026-07-17", Conditions: task.TaskConditions{Weather: []string{"sunny"}, Mode: []string{"Growth"}}},
			mockSetup: func(m *task.MockTaskRepository) {
				m.On("FindByID", mock.Anything, taskOID1.Hex()).Return(existing, nil)
				m.On("Update", mock.Anything, mock.AnythingOfType("*task.Task")).Return(nil)
			},
			expectError: false,
		},
		{
			name: "Fail: task type changed",
			task: &task.Task{ID: taskOID1, Title: "Read updated", Section: "dev", Type: "quantitative", Metrics: task.TaskMetrics{DailyTarget: 10}, StartDate: "2026-07-17", Conditions: task.TaskConditions{Weather: []string{"sunny"}, Mode: []string{"Growth"}}},
			mockSetup: func(m *task.MockTaskRepository) {
				m.On("FindByID", mock.Anything, taskOID1.Hex()).Return(existing, nil)
			},
			expectError: true,
			errContains: "task type cannot be changed after creation",
		},
		{
			name: "Fail: validation error",
			task: &task.Task{ID: taskOID1, Title: "", Section: "dev", Type: "boolean", StartDate: "2026-07-17"},
			mockSetup: func(m *task.MockTaskRepository) {
				m.On("FindByID", mock.Anything, taskOID1.Hex()).Return(existing, nil)
			},
			expectError: true,
			errContains: "title is required",
		},
		{
			name: "Fail: invalid section",
			task: &task.Task{ID: taskOID1, Title: "Read", Section: "invalid_section", Type: "boolean", StartDate: "2026-07-17", Conditions: task.TaskConditions{Weather: []string{"sunny"}, Mode: []string{"Growth"}}},
			mockSetup: func(m *task.MockTaskRepository) {
				m.On("FindByID", mock.Anything, taskOID1.Hex()).Return(existing, nil)
			},
			expectError: true,
			errContains: "section must be one of",
		},
		{
			name: "Fail: invalid type",
			task: &task.Task{ID: taskOID1, Title: "Read", Section: "dev", Type: "invalid_type", Conditions: task.TaskConditions{Weather: []string{"sunny"}, Mode: []string{"Growth"}}},
			mockSetup: func(m *task.MockTaskRepository) {
				m.On("FindByID", mock.Anything, taskOID1.Hex()).Return(&task.Task{ID: taskOID1, Status: "active", StartDate: "2026-07-17", Type: "invalid_type"}, nil)
			},
			expectError: true,
			errContains: "type must be one of",
		},
		{
			name: "Fail: update archived task",
			task: &task.Task{ID: taskOID1, Title: "Read", Section: "dev", Type: "boolean", StartDate: "2026-07-17", Conditions: task.TaskConditions{Weather: []string{"sunny"}, Mode: []string{"Growth"}}},
			mockSetup: func(m *task.MockTaskRepository) {
				m.On("FindByID", mock.Anything, taskOID1.Hex()).Return(&task.Task{ID: taskOID1, Status: "archived", StartDate: "2026-07-17"}, nil)
			},
			expectError: true,
			errContains: "cannot update an archived task",
		},
		{
			name: "Fail: not found",
			task: &task.Task{ID: taskOID2, Title: "Read", Section: "dev", Type: "boolean", StartDate: "2026-07-17"},
			mockSetup: func(m *task.MockTaskRepository) {
				m.On("FindByID", mock.Anything, taskOID2.Hex()).Return(nil, nil)
			},
			expectError: true,
			errContains: "task not found",
		},
		{
			name: "Fail: find error",
			task: &task.Task{ID: taskOID1, Title: "Read", Section: "dev", Type: "boolean", StartDate: "2026-07-17", Conditions: task.TaskConditions{Weather: []string{"sunny"}, Mode: []string{"Growth"}}},
			mockSetup: func(m *task.MockTaskRepository) {
				m.On("FindByID", mock.Anything, taskOID1.Hex()).Return(nil, errors.New("db error"))
			},
			expectError: true,
			errContains: "finding task",
		},
		{
			name: "Fail: update error",
			task: &task.Task{ID: taskOID1, Title: "Read", Section: "dev", Type: "boolean", StartDate: "2026-07-17", Conditions: task.TaskConditions{Weather: []string{"sunny"}, Mode: []string{"Growth"}}},
			mockSetup: func(m *task.MockTaskRepository) {
				m.On("FindByID", mock.Anything, taskOID1.Hex()).Return(existing, nil)
				m.On("Update", mock.Anything, mock.Anything).Return(errors.New("db error"))
			},
			expectError: true,
			errContains: "updating task",
		},
		{
			name: "Fail: missing weather condition",
			task: &task.Task{ID: taskOID1, Title: "Read", Section: "dev", Type: "boolean", StartDate: "2026-07-17", Conditions: task.TaskConditions{Weather: []string{}, Mode: []string{"Growth"}}},
			mockSetup: func(m *task.MockTaskRepository) {
				m.On("FindByID", mock.Anything, taskOID1.Hex()).Return(existing, nil)
			},
			expectError: true,
			errContains: "at least one weather condition is required",
		},
		{
			name: "Fail: invalid mode value",
			task: &task.Task{ID: taskOID1, Title: "Read", Section: "dev", Type: "boolean", StartDate: "2026-07-17", Conditions: task.TaskConditions{Weather: []string{"sunny"}, Mode: []string{"invalid"}}},
			mockSetup: func(m *task.MockTaskRepository) {
				m.On("FindByID", mock.Anything, taskOID1.Hex()).Return(existing, nil)
			},
			expectError: true,
			errContains: "invalid mode value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(task.MockTaskRepository)
			tt.mockSetup(repo)
			svc := task.NewTaskService(repo, new(task.MockTaskProgressAggregator), nil)

			updated, err := svc.Update(context.Background(), tt.task)
			if tt.expectError {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.True(t, strings.Contains(err.Error(), tt.errContains))
				}
				assert.Nil(t, updated)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, existing.Status, updated.Status)
			}
			repo.AssertExpectations(t)
		})
	}
}

func TestTaskService_Archive(t *testing.T) {
	tests := []struct {
		name        string
		id          string
		mockSetup   func(*task.MockTaskRepository)
		expectError bool
		errContains string
	}{
		{
			name: "Pass: successful archive",
			id:   taskOID1.Hex(),
			mockSetup: func(m *task.MockTaskRepository) {
				m.On("FindByID", mock.Anything, taskOID1.Hex()).Return(&task.Task{ID: taskOID1, Status: "active", StartDate: "2026-07-17"}, nil)
				m.On("Update", mock.Anything, mock.MatchedBy(func(task *task.Task) bool {
					return task.ID == taskOID1 && task.Status == "archived"
				})).Return(nil)
			},
			expectError: false,
		},
		{
			name: "Fail: task is already archived",
			id:   taskOID1.Hex(),
			mockSetup: func(m *task.MockTaskRepository) {
				m.On("FindByID", mock.Anything, taskOID1.Hex()).Return(&task.Task{ID: taskOID1, Status: "archived", StartDate: "2026-07-17"}, nil)
			},
			expectError: true,
			errContains: "task is already archived",
		},
		{
			name: "Fail: not found",
			id:   taskOID2.Hex(),
			mockSetup: func(m *task.MockTaskRepository) {
				m.On("FindByID", mock.Anything, taskOID2.Hex()).Return(nil, nil)
			},
			expectError: true,
			errContains: "task not found",
		},
		{
			name: "Fail: find error",
			id:   taskOID1.Hex(),
			mockSetup: func(m *task.MockTaskRepository) {
				m.On("FindByID", mock.Anything, taskOID1.Hex()).Return(nil, errors.New("db error"))
			},
			expectError: true,
			errContains: "finding task",
		},
		{
			name: "Fail: update error",
			id:   taskOID1.Hex(),
			mockSetup: func(m *task.MockTaskRepository) {
				m.On("FindByID", mock.Anything, taskOID1.Hex()).Return(&task.Task{ID: taskOID1, Status: "active", StartDate: "2026-07-17"}, nil)
				m.On("Update", mock.Anything, mock.Anything).Return(errors.New("db error"))
			},
			expectError: true,
			errContains: "archiving task",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(task.MockTaskRepository)
			tt.mockSetup(repo)
			svc := task.NewTaskService(repo, new(task.MockTaskProgressAggregator), nil)

			err := svc.Archive(context.Background(), tt.id)
			if tt.expectError {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.True(t, strings.Contains(err.Error(), tt.errContains))
				}
			} else {
				assert.NoError(t, err)
			}
			repo.AssertExpectations(t)
		})
	}
}

func TestTaskService_Get(t *testing.T) {
	existing := &task.Task{ID: taskOID1, Title: "Read"}

	tests := []struct {
		name        string
		id          string
		mockSetup   func(*task.MockTaskRepository)
		expectError bool
		errContains string
	}{
		{
			name: "Pass: successful get",
			id:   taskOID1.Hex(),
			mockSetup: func(m *task.MockTaskRepository) {
				m.On("FindByID", mock.Anything, taskOID1.Hex()).Return(existing, nil)
			},
			expectError: false,
		},
		{
			name: "Fail: not found",
			id:   taskOID2.Hex(),
			mockSetup: func(m *task.MockTaskRepository) {
				m.On("FindByID", mock.Anything, taskOID2.Hex()).Return(nil, nil)
			},
			expectError: true,
			errContains: "task not found",
		},
		{
			name: "Fail: find error",
			id:   taskOID1.Hex(),
			mockSetup: func(m *task.MockTaskRepository) {
				m.On("FindByID", mock.Anything, taskOID1.Hex()).Return(nil, errors.New("db error"))
			},
			expectError: true,
			errContains: "finding task",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(task.MockTaskRepository)
			tt.mockSetup(repo)
			svc := task.NewTaskService(repo, new(task.MockTaskProgressAggregator), nil)

			res, err := svc.Get(context.Background(), tt.id)
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, res)
				if tt.errContains != "" {
					assert.True(t, strings.Contains(err.Error(), tt.errContains))
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
			}
			repo.AssertExpectations(t)
		})
	}
}

func TestTaskService_List(t *testing.T) {
	tests := []struct {
		name        string
		mockSetup   func(*task.MockTaskRepository)
		expectError bool
	}{
		{
			name: "Pass: successful list",
			mockSetup: func(m *task.MockTaskRepository) {
				m.On("FindAll", mock.Anything).Return([]task.Task{{ID: taskOID1}}, nil)
			},
			expectError: false,
		},
		{
			name: "Fail: db error",
			mockSetup: func(m *task.MockTaskRepository) {
				m.On("FindAll", mock.Anything).Return(nil, errors.New("db error"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(task.MockTaskRepository)
			tt.mockSetup(repo)
			svc := task.NewTaskService(repo, new(task.MockTaskProgressAggregator), nil)

			res, err := svc.List(context.Background())
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.Len(t, res, 1)
			}
			repo.AssertExpectations(t)
		})
	}
}

func TestTaskService_GetProgressForActiveTasks(t *testing.T) {
	tests := []struct {
		name            string
		mockSetup       func(*task.MockTaskRepository, *task.MockTaskProgressAggregator)
		expectError     bool
		expectedResults int
	}{
		{
			name: "Pass: returns progress for all active tasks",
			mockSetup: func(repo *task.MockTaskRepository, agg *task.MockTaskProgressAggregator) {
				repo.On("FindActiveTasks", mock.Anything).Return([]task.Task{
					{ID: taskOID1, Title: "Kanji", Type: "quantitative", StartDate: "2026-07-17", Metrics: task.TaskMetrics{DailyTarget: 10, TotalTarget: 500}},
					{ID: taskOID2, Title: "Read News", Type: "boolean", StartDate: "2026-07-17", Metrics: task.TaskMetrics{DailyTarget: 1}},
					{ID: taskOID3, Title: "LeetCode", Type: "quantitative", StartDate: "2026-07-17", Metrics: task.TaskMetrics{DailyTarget: 3, TotalTarget: 100}},
				}, nil)
				agg.On("SumTaskProgressByIDs", mock.Anything, []primitive.ObjectID{taskOID1, taskOID2, taskOID3}).Return(map[primitive.ObjectID]int{
					taskOID1: 250,
					taskOID2: 5,
					taskOID3: 50,
				}, nil)
			},
			expectedResults: 3,
		},
		{
			name: "Pass: includes quantitative tasks with zero totalTarget",
			mockSetup: func(repo *task.MockTaskRepository, agg *task.MockTaskProgressAggregator) {
				repo.On("FindActiveTasks", mock.Anything).Return([]task.Task{
					{ID: taskOID1, Title: "Pushups", Type: "quantitative", StartDate: "2026-07-17", Metrics: task.TaskMetrics{DailyTarget: 10, TotalTarget: 0}},
				}, nil)
				agg.On("SumTaskProgressByIDs", mock.Anything, []primitive.ObjectID{taskOID1}).Return(map[primitive.ObjectID]int{}, nil)
			},
			expectedResults: 1,
		},
		{
			name: "Fail: repo error",
			mockSetup: func(repo *task.MockTaskRepository, agg *task.MockTaskProgressAggregator) {
				repo.On("FindActiveTasks", mock.Anything).Return(nil, errors.New("db error"))
			},
			expectError: true,
		},
		{
			name: "Fail: aggregator error",
			mockSetup: func(repo *task.MockTaskRepository, agg *task.MockTaskProgressAggregator) {
				repo.On("FindActiveTasks", mock.Anything).Return([]task.Task{
					{ID: taskOID1, Type: "quantitative", StartDate: "2026-07-17", Metrics: task.TaskMetrics{DailyTarget: 10, TotalTarget: 500}},
				}, nil)
				agg.On("SumTaskProgressByIDs", mock.Anything, []primitive.ObjectID{taskOID1}).Return((map[primitive.ObjectID]int)(nil), errors.New("agg error"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(task.MockTaskRepository)
			agg := new(task.MockTaskProgressAggregator)
			tt.mockSetup(repo, agg)
			svc := task.NewTaskService(repo, agg, nil)

			progress, err := svc.GetProgressForActiveTasks(context.Background())
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, progress)
			} else {
				assert.NoError(t, err)
				assert.Len(t, progress, tt.expectedResults)
			}
			repo.AssertExpectations(t)
			agg.AssertExpectations(t)
		})
	}
}

func TestTaskService_MigrateTask(t *testing.T) {
	newActiveTask := func() *task.Task {
		return &task.Task{
			ID:         taskOID1,
			Section:    "dev",
			Title:      "LeetCode",
			Type:       "quantitative",
			Status:     "active",
			Metrics:    task.TaskMetrics{DailyTarget: 3, TotalTarget: 100},
			Conditions: task.TaskConditions{Weather: []string{"sunny", "rainy"}, Mode: []string{"Growth", "Rest", "Office", "Remote"}},
		}
	}

	tests := []struct {
		name        string
		id          string
		mockSetup   func(*task.MockTaskRepository, *task.MockTaskProgressAggregator)
		expectError bool
		errContains string
	}{
		{
			name: "Pass: successful migration",
			id:   taskOID1.Hex(),
			mockSetup: func(repo *task.MockTaskRepository, agg *task.MockTaskProgressAggregator) {
				repo.On("FindByID", mock.Anything, taskOID1.Hex()).Return(newActiveTask(), nil)
				agg.On("SumTaskProgressByIDs", mock.Anything, []primitive.ObjectID{taskOID1}).Return(map[primitive.ObjectID]int{taskOID1: 100}, nil)
				repo.On("MigrateTaskAtomic", mock.Anything, mock.MatchedBy(func(t *task.Task) bool {
					return t.ID == taskOID1 && t.Status == "archived"
				}), mock.MatchedBy(func(t *task.Task) bool {
					return t.Status == "active" && t.Title == "LeetCode" && t.Section == "dev"
				})).Return(nil)
			},
		},
		{
			name: "Fail: task not found",
			id:   taskOID3.Hex(),
			mockSetup: func(repo *task.MockTaskRepository, agg *task.MockTaskProgressAggregator) {
				repo.On("FindByID", mock.Anything, taskOID3.Hex()).Return(nil, nil)
			},
			expectError: true,
			errContains: "task not found",
		},
		{
			name: "Fail: non-active task",
			id:   taskOID2.Hex(),
			mockSetup: func(repo *task.MockTaskRepository, agg *task.MockTaskProgressAggregator) {
				repo.On("FindByID", mock.Anything, taskOID2.Hex()).Return(&task.Task{ID: taskOID2, Status: "archived"}, nil)
			},
			expectError: true,
			errContains: "non-active",
		},
		{
			name: "Fail: progress not reached",
			id:   taskOID1.Hex(),
			mockSetup: func(repo *task.MockTaskRepository, agg *task.MockTaskProgressAggregator) {
				repo.On("FindByID", mock.Anything, taskOID1.Hex()).Return(newActiveTask(), nil)
				agg.On("SumTaskProgressByIDs", mock.Anything, []primitive.ObjectID{taskOID1}).Return(map[primitive.ObjectID]int{taskOID1: 50}, nil)
			},
			expectError: true,
			errContains: "has not reached the target",
		},
		{
			name: "Pass: allows overshooting (>100%)",
			id:   taskOID1.Hex(),
			mockSetup: func(repo *task.MockTaskRepository, agg *task.MockTaskProgressAggregator) {
				repo.On("FindByID", mock.Anything, taskOID1.Hex()).Return(newActiveTask(), nil)
				agg.On("SumTaskProgressByIDs", mock.Anything, []primitive.ObjectID{taskOID1}).Return(map[primitive.ObjectID]int{taskOID1: 150}, nil)
				repo.On("MigrateTaskAtomic", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
		},
		{
			name: "Fail: aggregator error",
			id:   taskOID1.Hex(),
			mockSetup: func(repo *task.MockTaskRepository, agg *task.MockTaskProgressAggregator) {
				repo.On("FindByID", mock.Anything, taskOID1.Hex()).Return(newActiveTask(), nil)
				agg.On("SumTaskProgressByIDs", mock.Anything, []primitive.ObjectID{taskOID1}).Return((map[primitive.ObjectID]int)(nil), errors.New("agg error"))
			},
			expectError: true,
			errContains: "checking task progress",
		},
		{
			name: "Fail: atomic migration fails",
			id:   taskOID1.Hex(),
			mockSetup: func(repo *task.MockTaskRepository, agg *task.MockTaskProgressAggregator) {
				repo.On("FindByID", mock.Anything, taskOID1.Hex()).Return(newActiveTask(), nil)
				agg.On("SumTaskProgressByIDs", mock.Anything, []primitive.ObjectID{taskOID1}).Return(map[primitive.ObjectID]int{taskOID1: 100}, nil)
				repo.On("MigrateTaskAtomic", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("db error"))
			},
			expectError: true,
			errContains: "atomic migration failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(task.MockTaskRepository)
			agg := new(task.MockTaskProgressAggregator)
			tt.mockSetup(repo, agg)
			svc := task.NewTaskService(repo, agg, nil)

			result, err := svc.MigrateTask(context.Background(), tt.id, task.MigrateTaskRequest{CompletionDate: "2026-07-17"})
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
				if tt.errContains != "" {
					assert.True(t, strings.Contains(err.Error(), tt.errContains))
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, "archived", result.ArchivedTask.Status)
				assert.Equal(t, "active", result.NewTask.Status)
				assert.Equal(t, "LeetCode", result.NewTask.Title)
			}
			repo.AssertExpectations(t)
			agg.AssertExpectations(t)
		})
	}
}

