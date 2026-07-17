package task_test

import (
	"context"
	"daily-seed/internal/task"
	"daily-seed/internal/testutil"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaskRepository(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	repo := task.NewTaskRepository(testutil.DB)

	t.Run("Create_And_FindByID", func(t *testing.T) {
		testutil.ClearDB(ctx)
		t1ID := primitive.NewObjectID()
		task := &task.Task{
			ID:     t1ID,
			Title:  "Test Task",
			Status: "active",
		}

		err := repo.Create(ctx, task)
		require.NoError(t, err)

		found, err := repo.FindByID(ctx, t1ID.Hex())
		require.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, "Test Task", found.Title)
	})

	t.Run("FindActiveTasks", func(t *testing.T) {
		testutil.ClearDB(ctx)
		t1ID := primitive.NewObjectID()
		require.NoError(t, repo.Create(ctx, &task.Task{ID: t1ID, Status: "active"}))
		require.NoError(t, repo.Create(ctx, &task.Task{ID: primitive.NewObjectID(), Status: "archived"}))

		active, err := repo.FindActiveTasks(ctx)
		require.NoError(t, err)
		assert.Len(t, active, 1)
		assert.Equal(t, t1ID, active[0].ID)
	})

	t.Run("Update_And_MigrateTaskAtomic", func(t *testing.T) {
		testutil.ClearDB(ctx)
		t1ID := primitive.NewObjectID()
		t1 := &task.Task{ID: t1ID, Title: "Old", Status: "active"}
		require.NoError(t, repo.Create(ctx, t1))

		t1.Title = "New"
		require.NoError(t, repo.Update(ctx, t1))

		found, err := repo.FindByID(ctx, t1ID.Hex())
		require.NoError(t, err)
		assert.Equal(t, "New", found.Title)

		newTask := &task.Task{ID: primitive.NewObjectID(), Title: "New Task", Status: "active"}
		require.NoError(t, repo.MigrateTaskAtomic(ctx, found, newTask))
		
		archived, _ := repo.FindByID(ctx, t1ID.Hex())
		assert.Equal(t, "archived", archived.Status)

		created, _ := repo.FindByID(ctx, newTask.ID.Hex())
		assert.Equal(t, "active", created.Status)
	})

	t.Run("Context_Cancellation", func(t *testing.T) {
		cancelCtx, cancelFunc := context.WithCancel(context.Background())
		cancelFunc() // cancel immediately
		err := repo.Create(cancelCtx, &task.Task{ID: primitive.NewObjectID()})
		assert.ErrorIs(t, err, context.Canceled)
		
		_, err = repo.FindByID(cancelCtx, primitive.NewObjectID().Hex())
		assert.ErrorIs(t, err, context.Canceled)
	})
}
