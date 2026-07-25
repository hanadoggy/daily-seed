package analytics_test

import (
	"context"
	"daily-seed/internal/analytics"
	"daily-seed/internal/daily"
	"daily-seed/internal/habit"
	"daily-seed/internal/task"
	"daily-seed/internal/testutil"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestAnalyticsHandler_Slice(t *testing.T) {
	ctx := context.Background()

	t.Run("GET /api/v1/analytics/heatmap", func(t *testing.T) {
		testutil.ClearDB(ctx)
		router := testutil.SetupRouter(testutil.DB)

		// Seed task & habit
		taskID := testutil.SeedTask(ctx, testutil.DB, task.Task{
			Title: "Code", Section: "dev", Type: "boolean", Status: "active", StartDate: "2026-01-01",
		})
		habitID := testutil.SeedHabit(ctx, testutil.DB, habit.Habit{
			Title: "Drink Water", Category: "Health", Status: "active",
		})

		// Seed daily record with 1 completed task and 1 completed habit, 1 incomplete task
		incompletTaskID := primitive.NewObjectID()
		testutil.SeedDailyRecord(ctx, testutil.DB, daily.DailyRecord{
			Date: "2026-03-15",
			Tasks: []daily.TaskEntry{
				{TaskID: taskID, IsCompleted: true},
				{TaskID: incompletTaskID, IsCompleted: false}, // should not count
			},
			Habits: []daily.HabitEntry{
				{HabitID: habitID, IsCompleted: true},
			},
		})

		// 1. Heatmap for 2026
		w := testutil.DoRequest(router, "GET", "/api/v1/analytics/heatmap?year=2026", nil)
		assert.Equal(t, http.StatusOK, w.Code)

		var res analytics.HeatmapResponse
		err := json.Unmarshal(w.Body.Bytes(), &res)
		require.NoError(t, err)
		assert.Len(t, res.Days, 365) // 2026 is regular year

		var march15 *analytics.HeatmapDay
		for i, day := range res.Days {
			if day.Date == "2026-03-15" {
				march15 = &res.Days[i]
				break
			}
		}
		require.NotNil(t, march15)
		assert.Equal(t, 2, march15.Total)
		assert.Equal(t, 1, march15.Habits)
		assert.Equal(t, 1, march15.SectionCounts["dev"])

		// 2. Default Year (when year query parameter is empty)
		currentYear := time.Now().Year()
		w = testutil.DoRequest(router, "GET", "/api/v1/analytics/heatmap", nil)
		assert.Equal(t, http.StatusOK, w.Code)
		var resDefault analytics.HeatmapResponse
		_ = json.Unmarshal(w.Body.Bytes(), &resDefault)
		if currentYear%4 == 0 {
			assert.Len(t, resDefault.Days, 366)
		} else {
			assert.Len(t, resDefault.Days, 365)
		}

		// 3. Leap Year (2024 -> 366 days)
		w = testutil.DoRequest(router, "GET", "/api/v1/analytics/heatmap?year=2024", nil)
		assert.Equal(t, http.StatusOK, w.Code)
		var res2024 analytics.HeatmapResponse
		_ = json.Unmarshal(w.Body.Bytes(), &res2024)
		assert.Len(t, res2024.Days, 366)

		// 4. Invalid Year Query
		w = testutil.DoRequest(router, "GET", "/api/v1/analytics/heatmap?year=invalid", nil)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
