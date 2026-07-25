package habit_test

import (
	"context"
	"daily-seed/internal/habit"
	"daily-seed/internal/testutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestHabitStore(t *testing.T) {
	ctx := context.Background()

	t.Run("FindAll Empty", func(t *testing.T) {
		testutil.ClearDB(ctx)
		store := habit.NewHabitStore(testutil.DB)

		habits, err := store.FindAll(ctx)
		require.NoError(t, err)
		assert.Empty(t, habits)
		assert.NotNil(t, habits)
	})

	t.Run("FindAll With Data", func(t *testing.T) {
		testutil.ClearDB(ctx)
		store := habit.NewHabitStore(testutil.DB)

		testutil.SeedHabit(ctx, testutil.DB, habit.Habit{Title: "Reading", Category: "Life", Status: "active"})
		testutil.SeedHabit(ctx, testutil.DB, habit.Habit{Title: "Running", Category: "Health", Status: "archived"})

		habits, err := store.FindAll(ctx)
		require.NoError(t, err)
		assert.Len(t, habits, 2)
	})

	t.Run("FindActiveHabits", func(t *testing.T) {
		testutil.ClearDB(ctx)
		store := habit.NewHabitStore(testutil.DB)

		testutil.SeedHabit(ctx, testutil.DB, habit.Habit{Title: "Reading", Category: "Life", Status: "active"})
		testutil.SeedHabit(ctx, testutil.DB, habit.Habit{Title: "Running", Category: "Health", Status: "archived"})

		habits, err := store.FindActiveHabits(ctx)
		require.NoError(t, err)
		require.Len(t, habits, 1)
		assert.Equal(t, "Reading", habits[0].Title)
	})

	t.Run("FindByID Exists & NotFound & InvalidID", func(t *testing.T) {
		testutil.ClearDB(ctx)
		store := habit.NewHabitStore(testutil.DB)

		id := testutil.SeedHabit(ctx, testutil.DB, habit.Habit{Title: "Reading", Category: "Life", Status: "active"})

		// Exists
		h, err := store.FindByID(ctx, id.Hex())
		require.NoError(t, err)
		require.NotNil(t, h)
		assert.Equal(t, "Reading", h.Title)

		// NotFound
		nonExistID := primitive.NewObjectID().Hex()
		h, err = store.FindByID(ctx, nonExistID)
		require.NoError(t, err)
		assert.Nil(t, h)

		// InvalidID
		_, err = store.FindByID(ctx, "invalid-hex")
		require.Error(t, err)
	})

	t.Run("Create and Update", func(t *testing.T) {
		testutil.ClearDB(ctx)
		store := habit.NewHabitStore(testutil.DB)

		newHabit := &habit.Habit{
			ID:       primitive.NewObjectID(),
			Title:    "New Habit",
			Category: "Study",
			Status:   "active",
		}
		err := store.Create(ctx, newHabit)
		require.NoError(t, err)

		h, err := store.FindByID(ctx, newHabit.ID.Hex())
		require.NoError(t, err)
		assert.Equal(t, "New Habit", h.Title)

		// Update
		h.Title = "Updated Habit"
		err = store.Update(ctx, h)
		require.NoError(t, err)

		updated, err := store.FindByID(ctx, newHabit.ID.Hex())
		require.NoError(t, err)
		assert.Equal(t, "Updated Habit", updated.Title)
	})

	t.Run("EnsureIndexes Idempotent", func(t *testing.T) {
		testutil.ClearDB(ctx)
		store := habit.NewHabitStore(testutil.DB)

		err := store.EnsureIndexes(ctx)
		require.NoError(t, err)

		err = store.EnsureIndexes(ctx)
		require.NoError(t, err)
	})
}
