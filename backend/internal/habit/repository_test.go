package habit_test

import (
	"context"
	"daily-seed/internal/habit"
	"daily-seed/internal/testutil"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHabitRepository(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	repo := habit.NewHabitRepository(testutil.DB)

	t.Run("Create_And_FindByID", func(t *testing.T) {
		testutil.ClearDB(ctx)
		h1ID := primitive.NewObjectID()
		habit := &habit.Habit{
			ID:     h1ID,
			Title:  "Drink Water",
			Status: "active",
		}

		err := repo.Create(ctx, habit)
		require.NoError(t, err)

		found, err := repo.FindByID(ctx, h1ID.Hex())
		require.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, "Drink Water", found.Title)
	})

	t.Run("FindActiveHabits", func(t *testing.T) {
		testutil.ClearDB(ctx)
		h1ID := primitive.NewObjectID()
		require.NoError(t, repo.Create(ctx, &habit.Habit{ID: h1ID, Status: "active"}))
		require.NoError(t, repo.Create(ctx, &habit.Habit{ID: primitive.NewObjectID(), Status: "archived"}))

		active, err := repo.FindActiveHabits(ctx)
		require.NoError(t, err)
		assert.Len(t, active, 1)
		assert.Equal(t, h1ID, active[0].ID)
	})

	t.Run("Context_Cancellation", func(t *testing.T) {
		cancelCtx, cancelFunc := context.WithCancel(context.Background())
		cancelFunc() // cancel immediately
		err := repo.Create(cancelCtx, &habit.Habit{ID: primitive.NewObjectID()})
		assert.ErrorIs(t, err, context.Canceled)
		
		_, err = repo.FindByID(cancelCtx, primitive.NewObjectID().Hex())
		assert.ErrorIs(t, err, context.Canceled)
	})
}
