package daily_test

import (
	"context"
	"daily-seed/internal/daily"
	"daily-seed/internal/testutil"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func TestDailyRecordRepository(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	repo := daily.NewDailyRecordRepository(testutil.DB)
	date := "2023-11-11"

	t.Run("FindByDate_NotFound", func(t *testing.T) {
		testutil.ClearDB(ctx)
		record, err := repo.FindByDate(ctx, date)
		require.NoError(t, err)
		assert.Nil(t, record)
	})

	t.Run("Upsert_And_FindByDate", func(t *testing.T) {
		testutil.ClearDB(ctx)
		record := &daily.DailyRecord{
			ID:   date,
			Date: date,
			Context: daily.DayContext{
				Mode: "Growth",
			},
		}

		err := repo.Upsert(ctx, record)
		require.NoError(t, err)

		found, err := repo.FindByDate(ctx, date)
		require.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, "Growth", found.Context.Mode)
	})

	t.Run("PatchByDate", func(t *testing.T) {
		testutil.ClearDB(ctx)
		record := &daily.DailyRecord{ID: date, Date: date}
		require.NoError(t, repo.Upsert(ctx, record))

		err := repo.PatchByDate(ctx, date, bson.M{"context.mode": "British Green"})
		require.NoError(t, err)

		found, err := repo.FindByDate(ctx, date)
		require.NoError(t, err)
		assert.Equal(t, "British Green", found.Context.Mode)
	})
}
