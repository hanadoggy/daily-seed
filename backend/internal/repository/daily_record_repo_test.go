package repository_test

import (
	"context"
	"testing"
	"time"

	"daily-seed/internal/model"
	"daily-seed/internal/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func TestDailyRecordRepository(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	repo := repository.NewDailyRecordRepository(testDB)
	date := "2023-11-11"

	t.Run("FindByDate_NotFound", func(t *testing.T) {
		clearDB(ctx)
		record, err := repo.FindByDate(ctx, date)
		require.NoError(t, err)
		assert.Nil(t, record)
	})

	t.Run("Upsert_And_FindByDate", func(t *testing.T) {
		clearDB(ctx)
		record := &model.DailyRecord{
			ID:   date,
			Date: date,
			Context: model.DayContext{
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
		clearDB(ctx)
		record := &model.DailyRecord{ID: date, Date: date}
		require.NoError(t, repo.Upsert(ctx, record))

		err := repo.PatchByDate(ctx, date, bson.M{"context.mode": "British Green"})
		require.NoError(t, err)

		found, err := repo.FindByDate(ctx, date)
		require.NoError(t, err)
		assert.Equal(t, "British Green", found.Context.Mode)
	})
}
