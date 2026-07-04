package habit_test

import (
	"context"
	"daily-seed/internal/habit"
	"daily-seed/internal/testutil"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHabitRepository(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	repo := habit.NewHabitRepository(testutil.DB)

	t.Run("Create_And_FindByID", func(t *testing.T) {
		testutil.ClearDB(ctx)
		habit := &habit.Habit{
			ID:     "habit_1",
			Title:  "Drink Water",
			Status: "active",
		}

		err := repo.Create(ctx, habit)
		require.NoError(t, err)

		found, err := repo.FindByID(ctx, "habit_1")
		require.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, "Drink Water", found.Title)
	})

	t.Run("FindActiveHabits", func(t *testing.T) {
		testutil.ClearDB(ctx)
		require.NoError(t, repo.Create(ctx, &habit.Habit{ID: "h1", Status: "active"}))
		require.NoError(t, repo.Create(ctx, &habit.Habit{ID: "h2", Status: "archived"}))

		active, err := repo.FindActiveHabits(ctx)
		require.NoError(t, err)
		assert.Len(t, active, 1)
		assert.Equal(t, "h1", active[0].ID)
	})

	t.Run("Context_Cancellation", func(t *testing.T) {
		cancelCtx, cancelFunc := context.WithCancel(context.Background())
		cancelFunc() // cancel immediately
		err := repo.Create(cancelCtx, &habit.Habit{ID: "ctx_test"})
		assert.ErrorIs(t, err, context.Canceled)
		
		_, err = repo.FindByID(cancelCtx, "h1")
		assert.ErrorIs(t, err, context.Canceled)
	})
}
