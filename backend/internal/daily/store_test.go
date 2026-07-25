package daily_test

import (
	"context"
	"daily-seed/internal/daily"
	"daily-seed/internal/testutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestDailyStore(t *testing.T) {
	ctx := context.Background()

	t.Run("FindByDate Exists & NotFound", func(t *testing.T) {
		testutil.ClearDB(ctx)
		store := daily.NewDailyStore(testutil.DB)

		testutil.SeedDailyRecord(ctx, testutil.DB, daily.DailyRecord{Date: "2026-01-01"})

		// Exists
		rec, err := store.FindByDate(ctx, "2026-01-01")
		require.NoError(t, err)
		require.NotNil(t, rec)
		assert.Equal(t, "2026-01-01", rec.Date)

		// NotFound
		rec, err = store.FindByDate(ctx, "2026-01-02")
		require.NoError(t, err)
		assert.Nil(t, rec)
	})

	t.Run("Upsert and PatchByDate", func(t *testing.T) {
		testutil.ClearDB(ctx)
		store := daily.NewDailyStore(testutil.DB)

		id := primitive.NewObjectID()
		rec := &daily.DailyRecord{
			ID:   id,
			Date: "2026-01-01",
			Context: daily.DayContext{
				Mode:    daily.ModeGrowth,
				Weather: "sunny",
			},
		}

		// Insert via Upsert
		err := store.Upsert(ctx, rec)
		require.NoError(t, err)

		found, err := store.FindByDate(ctx, "2026-01-01")
		require.NoError(t, err)
		assert.Equal(t, daily.ModeGrowth, found.Context.Mode)

		// PatchByDate
		err = store.PatchByDate(ctx, "2026-01-01", bson.M{"context.mode": daily.ModeRest})
		require.NoError(t, err)

		patched, err := store.FindByDate(ctx, "2026-01-01")
		require.NoError(t, err)
		assert.Equal(t, daily.ModeRest, patched.Context.Mode)
		assert.Equal(t, "sunny", patched.Context.Weather) // weather preserved

		// PatchByDate NotFound
		err = store.PatchByDate(ctx, "2026-01-02", bson.M{"context.mode": daily.ModeRest})
		require.Error(t, err)
	})

	t.Run("FindBetweenDates & SumTaskProgressByIDs", func(t *testing.T) {
		testutil.ClearDB(ctx)
		store := daily.NewDailyStore(testutil.DB)

		taskID := primitive.NewObjectID()

		testutil.SeedDailyRecord(ctx, testutil.DB, daily.DailyRecord{
			Date: "2026-01-01",
			Tasks: []daily.TaskEntry{
				{TaskID: taskID, ActualAmount: 10},
			},
		})
		testutil.SeedDailyRecord(ctx, testutil.DB, daily.DailyRecord{
			Date: "2026-01-05",
			Tasks: []daily.TaskEntry{
				{TaskID: taskID, ActualAmount: 15},
			},
		})
		testutil.SeedDailyRecord(ctx, testutil.DB, daily.DailyRecord{
			Date: "2026-02-01",
		})

		// FindBetweenDates
		recs, err := store.FindBetweenDates(ctx, "2026-01-01", "2026-01-31")
		require.NoError(t, err)
		assert.Len(t, recs, 2)

		// SumTaskProgressByIDs
		progressMap, err := store.SumTaskProgressByIDs(ctx, []primitive.ObjectID{taskID})
		require.NoError(t, err)
		assert.Equal(t, 25, progressMap[taskID])

		// SumTaskProgressByIDs empty taskIDs
		emptyMap, err := store.SumTaskProgressByIDs(ctx, []primitive.ObjectID{})
		require.NoError(t, err)
		assert.Empty(t, emptyMap)
	})

	t.Run("RemoveTaskFromRecordsBeforeDate", func(t *testing.T) {
		testutil.ClearDB(ctx)
		store := daily.NewDailyStore(testutil.DB)

		taskID := primitive.NewObjectID()

		testutil.SeedDailyRecord(ctx, testutil.DB, daily.DailyRecord{
			Date: "2026-01-01",
			Tasks: []daily.TaskEntry{
				{TaskID: taskID, ActualAmount: 10},
			},
		})
		testutil.SeedDailyRecord(ctx, testutil.DB, daily.DailyRecord{
			Date: "2026-01-10",
			Tasks: []daily.TaskEntry{
				{TaskID: taskID, ActualAmount: 10},
			},
		})

		err := store.RemoveTaskFromRecordsBeforeDate(ctx, taskID, "2026-01-05")
		require.NoError(t, err)

		rec1, _ := store.FindByDate(ctx, "2026-01-01")
		assert.Empty(t, rec1.Tasks) // removed

		rec2, _ := store.FindByDate(ctx, "2026-01-10")
		assert.Len(t, rec2.Tasks, 1) // kept
	})

	t.Run("EnsureIndexes Idempotent", func(t *testing.T) {
		testutil.ClearDB(ctx)
		store := daily.NewDailyStore(testutil.DB)

		err := store.EnsureIndexes(ctx)
		require.NoError(t, err)

		err = store.EnsureIndexes(ctx)
		require.NoError(t, err)
	})
}
