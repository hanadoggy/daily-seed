package task_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"daily-seed/internal/task"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
			task: &task.Task{Title: "Learn Go", Section: "dev", Type: "boolean"},
			mockSetup: func(m *task.MockTaskRepository) {
				m.On("Create", mock.Anything, mock.AnythingOfType("*task.Task")).Return(nil)
			},
			expectError: false,
		},
		{
			name: "Pass: successful quantitative task",
			task: &task.Task{Title: "Read", Section: "self_dev", Type: "quantitative", Metrics: task.TaskMetrics{DailyTarget: 10}},
			mockSetup: func(m *task.MockTaskRepository) {
				m.On("Create", mock.Anything, mock.AnythingOfType("*task.Task")).Return(nil)
			},
			expectError: false,
		},
		{
			name: "Fail: missing title",
			task: &task.Task{Title: "", Section: "dev", Type: "boolean"},
			mockSetup: func(m *task.MockTaskRepository) {},
			expectError: true,
			errContains: "title is required",
		},
		{
			name: "Fail: invalid section",
			task: &task.Task{Title: "Test", Section: "unknown", Type: "boolean"},
			mockSetup: func(m *task.MockTaskRepository) {},
			expectError: true,
			errContains: "section must be one of",
		},
		{
			name: "Fail: invalid type",
			task: &task.Task{Title: "Test", Section: "dev", Type: "unknown"},
			mockSetup: func(m *task.MockTaskRepository) {},
			expectError: true,
			errContains: "type must be one of",
		},
		{
			name: "Fail: zero target for quantitative",
			task: &task.Task{Title: "Test", Section: "dev", Type: "quantitative", Metrics: task.TaskMetrics{DailyTarget: 0}},
			mockSetup: func(m *task.MockTaskRepository) {},
			expectError: true,
			errContains: "dailyTarget must be positive",
		},
		{
			name: "Fail: negative total target",
			task: &task.Task{Title: "Test", Section: "dev", Type: "quantitative", Metrics: task.TaskMetrics{DailyTarget: 1, TotalTarget: -10}},
			mockSetup: func(m *task.MockTaskRepository) {},
			expectError: true,
			errContains: "totalTarget cannot be negative",
		},
		{
			name: "Fail: repo error",
			task: &task.Task{Title: "Test", Section: "dev", Type: "boolean"},
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
			svc := task.NewTaskService(repo, new(task.MockTaskProgressAggregator))

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
	existing := &task.Task{ID: "task_1", Status: "active"}

	tests := []struct {
		name        string
		task        *task.Task
		mockSetup   func(*task.MockTaskRepository)
		expectError bool
		errContains string
	}{
		{
			name: "Pass: successful update",
			task: &task.Task{ID: "task_1", Title: "Read updated", Section: "dev", Type: "boolean"},
			mockSetup: func(m *task.MockTaskRepository) {
				m.On("FindByID", mock.Anything, "task_1").Return(existing, nil)
				m.On("Update", mock.Anything, mock.AnythingOfType("*task.Task")).Return(nil)
			},
			expectError: false,
		},
		{
			name: "Fail: validation error",
			task: &task.Task{ID: "task_1", Title: "", Section: "dev", Type: "boolean"},
			mockSetup: func(m *task.MockTaskRepository) {
				m.On("FindByID", mock.Anything, "task_1").Return(existing, nil)
			},
			expectError: true,
			errContains: "title is required",
		},
		{
			name: "Fail: invalid section",
			task: &task.Task{ID: "task_1", Title: "Read", Section: "invalid_section", Type: "boolean"},
			mockSetup: func(m *task.MockTaskRepository) {
				m.On("FindByID", mock.Anything, "task_1").Return(existing, nil)
			},
			expectError: true,
			errContains: "section must be one of",
		},
		{
			name: "Fail: invalid type",
			task: &task.Task{ID: "task_1", Title: "Read", Section: "dev", Type: "invalid_type"},
			mockSetup: func(m *task.MockTaskRepository) {
				m.On("FindByID", mock.Anything, "task_1").Return(existing, nil)
			},
			expectError: true,
			errContains: "type must be one of",
		},
		{
			name: "Fail: not found",
			task: &task.Task{ID: "task_2", Title: "Read", Section: "dev", Type: "boolean"},
			mockSetup: func(m *task.MockTaskRepository) {
				m.On("FindByID", mock.Anything, "task_2").Return(nil, nil)
			},
			expectError: true,
			errContains: "task not found",
		},
		{
			name: "Fail: find error",
			task: &task.Task{ID: "task_1", Title: "Read", Section: "dev", Type: "boolean"},
			mockSetup: func(m *task.MockTaskRepository) {
				m.On("FindByID", mock.Anything, "task_1").Return(nil, errors.New("db error"))
			},
			expectError: true,
			errContains: "finding task",
		},
		{
			name: "Fail: update error",
			task: &task.Task{ID: "task_1", Title: "Read", Section: "dev", Type: "boolean"},
			mockSetup: func(m *task.MockTaskRepository) {
				m.On("FindByID", mock.Anything, "task_1").Return(existing, nil)
				m.On("Update", mock.Anything, mock.Anything).Return(errors.New("db error"))
			},
			expectError: true,
			errContains: "updating task",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(task.MockTaskRepository)
			tt.mockSetup(repo)
			svc := task.NewTaskService(repo, new(task.MockTaskProgressAggregator))

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
	existing := &task.Task{ID: "task_1", Status: "active"}

	tests := []struct {
		name        string
		id          string
		mockSetup   func(*task.MockTaskRepository)
		expectError bool
		errContains string
	}{
		{
			name: "Pass: successful archive",
			id:   "task_1",
			mockSetup: func(m *task.MockTaskRepository) {
				m.On("FindByID", mock.Anything, "task_1").Return(existing, nil)
				m.On("Update", mock.Anything, mock.MatchedBy(func(task *task.Task) bool {
					return task.ID == "task_1" && task.Status == "archived"
				})).Return(nil)
			},
			expectError: false,
		},
		{
			name: "Fail: not found",
			id:   "task_2",
			mockSetup: func(m *task.MockTaskRepository) {
				m.On("FindByID", mock.Anything, "task_2").Return(nil, nil)
			},
			expectError: true,
			errContains: "task not found",
		},
		{
			name: "Fail: find error",
			id:   "task_1",
			mockSetup: func(m *task.MockTaskRepository) {
				m.On("FindByID", mock.Anything, "task_1").Return(nil, errors.New("db error"))
			},
			expectError: true,
			errContains: "finding task",
		},
		{
			name: "Fail: update error",
			id:   "task_1",
			mockSetup: func(m *task.MockTaskRepository) {
				m.On("FindByID", mock.Anything, "task_1").Return(existing, nil)
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
			svc := task.NewTaskService(repo, new(task.MockTaskProgressAggregator))

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
	existing := &task.Task{ID: "task_1", Title: "Read"}

	tests := []struct {
		name        string
		id          string
		mockSetup   func(*task.MockTaskRepository)
		expectError bool
		errContains string
	}{
		{
			name: "Pass: successful get",
			id:   "task_1",
			mockSetup: func(m *task.MockTaskRepository) {
				m.On("FindByID", mock.Anything, "task_1").Return(existing, nil)
			},
			expectError: false,
		},
		{
			name: "Fail: not found",
			id:   "task_2",
			mockSetup: func(m *task.MockTaskRepository) {
				m.On("FindByID", mock.Anything, "task_2").Return(nil, nil)
			},
			expectError: true,
			errContains: "task not found",
		},
		{
			name: "Fail: find error",
			id:   "task_1",
			mockSetup: func(m *task.MockTaskRepository) {
				m.On("FindByID", mock.Anything, "task_1").Return(nil, errors.New("db error"))
			},
			expectError: true,
			errContains: "finding task",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(task.MockTaskRepository)
			tt.mockSetup(repo)
			svc := task.NewTaskService(repo, new(task.MockTaskProgressAggregator))

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
				m.On("FindAll", mock.Anything).Return([]task.Task{{ID: "1"}}, nil)
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
			svc := task.NewTaskService(repo, new(task.MockTaskProgressAggregator))

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
			name: "Pass: returns progress for quantitative tasks with totalTarget",
			mockSetup: func(repo *task.MockTaskRepository, agg *task.MockTaskProgressAggregator) {
				repo.On("FindActiveTasks", mock.Anything).Return([]task.Task{
					{ID: "task_1", Title: "Kanji", Type: "quantitative", Metrics: task.TaskMetrics{DailyTarget: 10, TotalTarget: 500}},
					{ID: "task_2", Title: "Read News", Type: "boolean", Metrics: task.TaskMetrics{DailyTarget: 1}},
					{ID: "task_3", Title: "LeetCode", Type: "quantitative", Metrics: task.TaskMetrics{DailyTarget: 3, TotalTarget: 100}},
				}, nil)
				agg.On("SumTaskProgress", mock.Anything, "task_1").Return(250, nil)
				agg.On("SumTaskProgress", mock.Anything, "task_3").Return(50, nil)
			},
			expectedResults: 2,
		},
		{
			name: "Pass: skips quantitative tasks with zero totalTarget",
			mockSetup: func(repo *task.MockTaskRepository, agg *task.MockTaskProgressAggregator) {
				repo.On("FindActiveTasks", mock.Anything).Return([]task.Task{
					{ID: "task_1", Title: "Kanji", Type: "quantitative", Metrics: task.TaskMetrics{DailyTarget: 10, TotalTarget: 0}},
				}, nil)
			},
			expectedResults: 0,
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
					{ID: "task_1", Type: "quantitative", Metrics: task.TaskMetrics{DailyTarget: 10, TotalTarget: 500}},
				}, nil)
				agg.On("SumTaskProgress", mock.Anything, "task_1").Return(0, errors.New("agg error"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(task.MockTaskRepository)
			agg := new(task.MockTaskProgressAggregator)
			tt.mockSetup(repo, agg)
			svc := task.NewTaskService(repo, agg)

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
			ID:         "task_1",
			Section:    "dev",
			Title:      "LeetCode",
			Type:       "quantitative",
			Status:     "active",
			Metrics:    task.TaskMetrics{DailyTarget: 3, TotalTarget: 100},
			Conditions: task.TaskConditions{Weather: "any", Mode: "any"},
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
			id:   "task_1",
			mockSetup: func(repo *task.MockTaskRepository, agg *task.MockTaskProgressAggregator) {
				repo.On("FindByID", mock.Anything, "task_1").Return(newActiveTask(), nil)
				agg.On("SumTaskProgress", mock.Anything, "task_1").Return(100, nil)
				repo.On("Update", mock.Anything, mock.MatchedBy(func(t *task.Task) bool {
					return t.ID == "task_1" && t.Status == "archived"
				})).Return(nil)
				repo.On("Create", mock.Anything, mock.MatchedBy(func(t *task.Task) bool {
					return t.Status == "active" && t.Title == "LeetCode" && t.Section == "dev"
				})).Return(nil)
			},
		},
		{
			name: "Fail: task not found",
			id:   "task_999",
			mockSetup: func(repo *task.MockTaskRepository, agg *task.MockTaskProgressAggregator) {
				repo.On("FindByID", mock.Anything, "task_999").Return(nil, nil)
			},
			expectError: true,
			errContains: "task not found",
		},
		{
			name: "Fail: non-active task",
			id:   "task_2",
			mockSetup: func(repo *task.MockTaskRepository, agg *task.MockTaskProgressAggregator) {
				repo.On("FindByID", mock.Anything, "task_2").Return(&task.Task{ID: "task_2", Status: "archived"}, nil)
			},
			expectError: true,
			errContains: "non-active",
		},
		{
			name: "Fail: progress not reached",
			id:   "task_1",
			mockSetup: func(repo *task.MockTaskRepository, agg *task.MockTaskProgressAggregator) {
				repo.On("FindByID", mock.Anything, "task_1").Return(newActiveTask(), nil)
				agg.On("SumTaskProgress", mock.Anything, "task_1").Return(50, nil)
			},
			expectError: true,
			errContains: "has not reached the target",
		},
		{
			name: "Pass: allows overshooting (>100%)",
			id:   "task_1",
			mockSetup: func(repo *task.MockTaskRepository, agg *task.MockTaskProgressAggregator) {
				repo.On("FindByID", mock.Anything, "task_1").Return(newActiveTask(), nil)
				agg.On("SumTaskProgress", mock.Anything, "task_1").Return(150, nil)
				repo.On("Update", mock.Anything, mock.Anything).Return(nil)
				repo.On("Create", mock.Anything, mock.Anything).Return(nil)
			},
		},
		{
			name: "Fail: aggregator error",
			id:   "task_1",
			mockSetup: func(repo *task.MockTaskRepository, agg *task.MockTaskProgressAggregator) {
				repo.On("FindByID", mock.Anything, "task_1").Return(newActiveTask(), nil)
				agg.On("SumTaskProgress", mock.Anything, "task_1").Return(0, errors.New("agg error"))
			},
			expectError: true,
			errContains: "checking task progress",
		},
		{
			name: "Fail: Update (archive old) fails",
			id:   "task_1",
			mockSetup: func(repo *task.MockTaskRepository, agg *task.MockTaskProgressAggregator) {
				repo.On("FindByID", mock.Anything, "task_1").Return(newActiveTask(), nil)
				agg.On("SumTaskProgress", mock.Anything, "task_1").Return(100, nil)
				repo.On("Update", mock.Anything, mock.Anything).Return(errors.New("db error"))
			},
			expectError: true,
			errContains: "archiving task for migration",
		},
		{
			name: "Fail: Create (new task) fails",
			id:   "task_1",
			mockSetup: func(repo *task.MockTaskRepository, agg *task.MockTaskProgressAggregator) {
				repo.On("FindByID", mock.Anything, "task_1").Return(newActiveTask(), nil)
				agg.On("SumTaskProgress", mock.Anything, "task_1").Return(100, nil)
				repo.On("Update", mock.Anything, mock.Anything).Return(nil)
				repo.On("Create", mock.Anything, mock.Anything).Return(errors.New("db error"))
			},
			expectError: true,
			errContains: "creating migrated task",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(task.MockTaskRepository)
			agg := new(task.MockTaskProgressAggregator)
			tt.mockSetup(repo, agg)
			svc := task.NewTaskService(repo, agg)

			result, err := svc.MigrateTask(context.Background(), tt.id)
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

