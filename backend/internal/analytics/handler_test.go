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

	t.Run("GET /api/v1/analytics/summary", func(t *testing.T) {
		testutil.ClearDB(ctx)
		router := testutil.SetupRouter(testutil.DB)

		// Seed Task & Habit
		devTaskID := testutil.SeedTask(ctx, testutil.DB, task.Task{
			Title: "Refactor Architecture", Section: "dev", Type: "quantitative", Status: "active", StartDate: "2026-01-01",
		})
		habitID := testutil.SeedHabit(ctx, testutil.DB, habit.Habit{
			Title: "Morning Meditation", Category: "Mindfulness", Status: "active",
		})

		// 2026-07-22 is Wednesday. Week range: 2026-07-19 (Sunday) ~ 2026-07-25 (Saturday). Total days: 7.
		// Seed 2 records in this week
		testutil.SeedDailyRecord(ctx, testutil.DB, daily.DailyRecord{
			Date: "2026-07-20", // Monday
			Context: daily.DayContext{
				Mode:    daily.ModeGrowth,
				Weather: "sunny",
			},
			Tasks: []daily.TaskEntry{
				{TaskID: devTaskID, TargetAmount: 5, ActualAmount: 5, IsCompleted: true},
			},
			Habits: []daily.HabitEntry{
				{HabitID: habitID, IsCompleted: true},
			},
			Journal: daily.Journal{
				OneLineReview:  "Great focus today!",
				ThreeLineDiary: "Finished task 1.\nStarted task 2.\nFelt good.",
			},
		})

		testutil.SeedDailyRecord(ctx, testutil.DB, daily.DailyRecord{
			Date: "2026-07-21", // Tuesday
			Context: daily.DayContext{
				Mode:    daily.ModeGrowth,
				Weather: "rainy",
			},
			Tasks: []daily.TaskEntry{
				{TaskID: devTaskID, TargetAmount: 5, ActualAmount: 2, IsCompleted: false},
			},
			Habits: []daily.HabitEntry{
				{HabitID: habitID, IsCompleted: false},
			},
		})

		// 1. Weekly summary request
		w := testutil.DoRequest(router, "GET", "/api/v1/analytics/summary?period=weekly&date=2026-07-22", nil)
		assert.Equal(t, http.StatusOK, w.Code)

		var summaryRes analytics.SummaryResponse
		err := json.Unmarshal(w.Body.Bytes(), &summaryRes)
		require.NoError(t, err)

		assert.Equal(t, "weekly", summaryRes.Period)
		assert.Equal(t, "2026-07-19", summaryRes.StartDate)
		assert.Equal(t, "2026-07-25", summaryRes.EndDate)
		assert.Equal(t, 7, summaryRes.TotalDays)
		assert.Equal(t, 2, summaryRes.RecordedDays)

		// Tasks: Total completed 7 / target 10 = 70.0%
		assert.Equal(t, 70.0, summaryRes.TaskCompletion.Overall)
		assert.Equal(t, 70.0, summaryRes.TaskCompletion.Sections["dev"])
		require.Len(t, summaryRes.TaskCompletion.PerTask, 1)
		assert.Equal(t, "Refactor Architecture", summaryRes.TaskCompletion.PerTask[0].Title)
		assert.Equal(t, 70.0, summaryRes.TaskCompletion.PerTask[0].Rate)

		// Habits: Total completed 1 / tracked 2 = 50.0%
		assert.Equal(t, 50.0, summaryRes.HabitCompletion.Overall)
		assert.Equal(t, 50.0, summaryRes.HabitCompletion.Categories["Mindfulness"])
		require.Len(t, summaryRes.HabitCompletion.PerHabit, 1)
		assert.Equal(t, "Morning Meditation", summaryRes.HabitCompletion.PerHabit[0].Title)
		assert.Equal(t, 50.0, summaryRes.HabitCompletion.PerHabit[0].Rate)

		// Mode distribution: Growth: 2
		assert.Equal(t, 2, summaryRes.ModeDistribution["Growth"])

		// Journals: 1 entry
		require.Len(t, summaryRes.Journals, 1)
		assert.Equal(t, "2026-07-20", summaryRes.Journals[0].Date)
		assert.Equal(t, "Great focus today!", summaryRes.Journals[0].OneLineReview)

		// 2. Monthly summary request
		wMonth := testutil.DoRequest(router, "GET", "/api/v1/analytics/summary?period=monthly&date=2026-07-15", nil)
		assert.Equal(t, http.StatusOK, wMonth.Code)

		var monthRes analytics.SummaryResponse
		err = json.Unmarshal(wMonth.Body.Bytes(), &monthRes)
		require.NoError(t, err)

		assert.Equal(t, "monthly", monthRes.Period)
		assert.Equal(t, "2026-07-01", monthRes.StartDate)
		assert.Equal(t, "2026-07-31", monthRes.EndDate)
		assert.Equal(t, 31, monthRes.TotalDays)
		assert.Equal(t, 2, monthRes.RecordedDays)

		// 3. Invalid Period Error
		wInvalid := testutil.DoRequest(router, "GET", "/api/v1/analytics/summary?period=yearly", nil)
		assert.Equal(t, http.StatusBadRequest, wInvalid.Code)
	})
}
