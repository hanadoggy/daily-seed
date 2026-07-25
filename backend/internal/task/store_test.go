package task_test

import (
	"context"
	"daily-seed/internal/task"
	"daily-seed/internal/testutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestTaskStore(t *testing.T) {
	ctx := context.Background()

	t.Run("FindAll Empty & With Data", func(t *testing.T) {
		testutil.ClearDB(ctx)
		store := task.NewTaskStore(testutil.DB)

		// Empty
		tasks, err := store.FindAll(ctx)
		require.NoError(t, err)
		assert.Empty(t, tasks)

		// Seed
		testutil.SeedTask(ctx, testutil.DB, task.Task{Title: "Study Japanese", Section: "japanese", Type: "quantitative", Status: "active"})
		testutil.SeedTask(ctx, testutil.DB, task.Task{Title: "Workout", Section: "exercise", Type: "boolean", Status: "archived"})

		tasks, err = store.FindAll(ctx)
		require.NoError(t, err)
		assert.Len(t, tasks, 2)
	})

	t.Run("FindActiveTasks", func(t *testing.T) {
		testutil.ClearDB(ctx)
		store := task.NewTaskStore(testutil.DB)

		testutil.SeedTask(ctx, testutil.DB, task.Task{Title: "Study Japanese", Section: "japanese", Type: "quantitative", Status: "active"})
		testutil.SeedTask(ctx, testutil.DB, task.Task{Title: "Workout", Section: "exercise", Type: "boolean", Status: "archived"})

		actives, err := store.FindActiveTasks(ctx)
		require.NoError(t, err)
		require.Len(t, actives, 1)
		assert.Equal(t, "Study Japanese", actives[0].Title)
	})

	t.Run("FindByID Exists & NotFound & InvalidID", func(t *testing.T) {
		testutil.ClearDB(ctx)
		store := task.NewTaskStore(testutil.DB)

		id := testutil.SeedTask(ctx, testutil.DB, task.Task{Title: "Study Japanese", Section: "japanese", Type: "quantitative", Status: "active"})

		// Exists
		tk, err := store.FindByID(ctx, id.Hex())
		require.NoError(t, err)
		require.NotNil(t, tk)
		assert.Equal(t, "Study Japanese", tk.Title)

		// NotFound
		nonExistID := primitive.NewObjectID().Hex()
		tk, err = store.FindByID(ctx, nonExistID)
		require.NoError(t, err)
		assert.Nil(t, tk)

		// InvalidID
		_, err = store.FindByID(ctx, "invalid-hex")
		require.Error(t, err)
	})

	t.Run("Create and Update", func(t *testing.T) {
		testutil.ClearDB(ctx)
		store := task.NewTaskStore(testutil.DB)

		newTask := &task.Task{
			ID:        primitive.NewObjectID(),
			Section:   "dev",
			Title:     "Write Code",
			Type:      "quantitative",
			Status:    "active",
			StartDate: "2026-01-01",
		}
		err := store.Create(ctx, newTask)
		require.NoError(t, err)

		tk, err := store.FindByID(ctx, newTask.ID.Hex())
		require.NoError(t, err)
		assert.Equal(t, "Write Code", tk.Title)

		// Update
		tk.Title = "Refactor Code"
		err = store.Update(ctx, tk)
		require.NoError(t, err)

		updated, err := store.FindByID(ctx, newTask.ID.Hex())
		require.NoError(t, err)
		assert.Equal(t, "Refactor Code", updated.Title)
	})

	t.Run("MigrateTaskAtomic", func(t *testing.T) {
		testutil.ClearDB(ctx)
		store := task.NewTaskStore(testutil.DB)

		oldID := testutil.SeedTask(ctx, testutil.DB, task.Task{
			Title:     "Old Goal",
			Section:   "self_dev",
			Type:      "quantitative",
			Status:    "active",
			StartDate: "2026-01-01",
		})

		oldTask, err := store.FindByID(ctx, oldID.Hex())
		require.NoError(t, err)
		oldTask.Status = "archived"
		oldTask.EndDate = "2026-01-10"

		newTask := &task.Task{
			ID:        primitive.NewObjectID(),
			Title:     "New Goal",
			Section:   "self_dev",
			Type:      "quantitative",
			Status:    "active",
			StartDate: "2026-01-11",
		}

		err = store.MigrateTaskAtomic(ctx, oldTask, newTask)
		require.NoError(t, err)

		// Check old task archived
		archived, err := store.FindByID(ctx, oldID.Hex())
		require.NoError(t, err)
		assert.Equal(t, "archived", archived.Status)
		assert.Equal(t, "2026-01-10", archived.EndDate)

		// Check new task created
		created, err := store.FindByID(ctx, newTask.ID.Hex())
		require.NoError(t, err)
		assert.Equal(t, "active", created.Status)
		assert.Equal(t, "2026-01-11", created.StartDate)
	})

	t.Run("EnsureIndexes Idempotent", func(t *testing.T) {
		testutil.ClearDB(ctx)
		store := task.NewTaskStore(testutil.DB)

		err := store.EnsureIndexes(ctx)
		require.NoError(t, err)

		err = store.EnsureIndexes(ctx)
		require.NoError(t, err)
	})
}
