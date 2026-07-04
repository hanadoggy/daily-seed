package task_test

import (
	"context"
	"daily-seed/internal/task"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTaskService_Create(t *testing.T) {
	mockRepo := new(task.MockTaskRepository)
	svc := task.NewTaskService(mockRepo)
	ctx := context.Background()

	t.Run("invalid_title", func(t *testing.T) {
		task := &task.Task{Section: "dev", Type: "boolean"}
		_, err := svc.Create(ctx, task)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "title is required")
	})

	t.Run("valid_task", func(t *testing.T) {
		task := &task.Task{Title: "Learn Go", Section: "dev", Type: "boolean"}
		mockRepo.On("Create", ctx, mock.AnythingOfType("*task.Task")).Return(nil)

		created, err := svc.Create(ctx, task)
		assert.NoError(t, err)
		assert.Equal(t, "active", created.Status)
		mockRepo.AssertExpectations(t)
	})
}

func TestTaskService_Archive(t *testing.T) {
	mockRepo := new(task.MockTaskRepository)
	svc := task.NewTaskService(mockRepo)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		existing := &task.Task{ID: "task_1", Status: "active"}
		mockRepo.On("FindByID", ctx, "task_1").Return(existing, nil).Once()
		mockRepo.On("Update", ctx, mock.MatchedBy(func(task *task.Task) bool {
			return task.ID == "task_1" && task.Status == "archived"
		})).Return(nil).Once()

		err := svc.Archive(ctx, "task_1")
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not_found", func(t *testing.T) {
		mockRepo.On("FindByID", ctx, "task_2").Return(nil, nil).Once()
		err := svc.Archive(ctx, "task_2")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "task not found")
	})
}
