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

func TestHabitRepository(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	repo := repository.NewHabitRepository(testDB)

	t.Run("Create_And_FindByID", func(t *testing.T) {
		clearDB(ctx)
		habit := &model.Habit{
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
		clearDB(ctx)
		require.NoError(t, repo.Create(ctx, &model.Habit{ID: "h1", Status: "active"}))
		require.NoError(t, repo.Create(ctx, &model.Habit{ID: "h2", Status: "archived"}))

		active, err := repo.FindActiveHabits(ctx)
		require.NoError(t, err)
		assert.Len(t, active, 1)
		assert.Equal(t, "h1", active[0].ID)
	})
}
