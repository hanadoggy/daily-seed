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
			svc := task.NewTaskService(repo)

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
			svc := task.NewTaskService(repo)

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
			svc := task.NewTaskService(repo)

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
			svc := task.NewTaskService(repo)

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
				m.On("FindActiveTasks", mock.Anything).Return([]task.Task{{ID: "1"}}, nil)
			},
			expectError: false,
		},
		{
			name: "Fail: db error",
			mockSetup: func(m *task.MockTaskRepository) {
				m.On("FindActiveTasks", mock.Anything).Return(nil, errors.New("db error"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(task.MockTaskRepository)
			tt.mockSetup(repo)
			svc := task.NewTaskService(repo)

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
