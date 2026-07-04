package repository_test

import (
	"context"
	"testing"
	"time"

	"daily-seed/internal/model"
	"daily-seed/internal/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaskRepository(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	repo := repository.NewTaskRepository(testDB)

	t.Run("Create_And_FindByID", func(t *testing.T) {
		clearDB(ctx)
		task := &model.Task{
			ID:     "task_1",
			Title:  "Test Task",
			Status: "active",
		}

		err := repo.Create(ctx, task)
		require.NoError(t, err)

		found, err := repo.FindByID(ctx, "task_1")
		require.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, "Test Task", found.Title)
	})

	t.Run("FindActiveTasks", func(t *testing.T) {
		clearDB(ctx)
		require.NoError(t, repo.Create(ctx, &model.Task{ID: "t1", Status: "active"}))
		require.NoError(t, repo.Create(ctx, &model.Task{ID: "t2", Status: "archived"}))

		active, err := repo.FindActiveTasks(ctx)
		require.NoError(t, err)
		assert.Len(t, active, 1)
		assert.Equal(t, "t1", active[0].ID)
	})

	t.Run("Update_And_Delete", func(t *testing.T) {
		clearDB(ctx)
		task := &model.Task{ID: "t1", Title: "Old", Status: "active"}
		require.NoError(t, repo.Create(ctx, task))

		task.Title = "New"
		require.NoError(t, repo.Update(ctx, task))

		found, err := repo.FindByID(ctx, "t1")
		require.NoError(t, err)
		assert.Equal(t, "New", found.Title)

		require.NoError(t, repo.Delete(ctx, "t1"))
		found, err = repo.FindByID(ctx, "t1")
		require.NoError(t, err)
		assert.Nil(t, found)
	})
}
