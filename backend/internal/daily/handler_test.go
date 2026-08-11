package daily_test

import (
	"context"
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
)

func TestDailyHandler_Slice(t *testing.T) {
	ctx := context.Background()

	t.Run("GET /api/v1/daily/:date - AutoGenerate & AppendMissing", func(t *testing.T) {
		testutil.ClearDB(ctx)
		router := testutil.SetupRouter(testutil.DB)

		today := time.Now().In(time.FixedZone("JST", 9*3600)).Format("2006-01-02")

		// Active task & habit
		taskID := testutil.SeedTask(ctx, testutil.DB, task.Task{
			Title: "Kanji", Section: "japanese", Type: "quantitative",
			Metrics: task.TaskMetrics{DailyTarget: 10}, Status: "active", StartDate: "2026-01-01",
		})
		habitID := testutil.SeedHabit(ctx, testutil.DB, habit.Habit{
			Title: "Meditation", Category: "Life", Status: "active",
		})
		// Archived task (should not be generated)
		testutil.SeedTask(ctx, testutil.DB, task.Task{
			Title: "Old Task", Section: "dev", Type: "boolean", Status: "archived", StartDate: "2026-01-01",
		})

		// 1. Auto generate today's record
		w := testutil.DoRequest(router, "GET", "/api/v1/daily/"+today, nil)
		assert.Equal(t, http.StatusOK, w.Code)
		var rec daily.DailyRecord
		err := json.Unmarshal(w.Body.Bytes(), &rec)
		require.NoError(t, err)
		assert.Equal(t, today, rec.Date)
		require.Len(t, rec.Tasks, 1)
		assert.Equal(t, taskID, rec.Tasks[0].TaskID)
		require.Len(t, rec.Habits, 1)
		assert.Equal(t, habitID, rec.Habits[0].HabitID)

		// 2. Add a new active task & re-fetch -> Append missing entry
		newTaskID := testutil.SeedTask(ctx, testutil.DB, task.Task{
			Title: "Pushups", Section: "exercise", Type: "boolean", Status: "active", StartDate: "2026-01-01",
		})
		w = testutil.DoRequest(router, "GET", "/api/v1/daily/"+today, nil)
		assert.Equal(t, http.StatusOK, w.Code)
		var recAppended daily.DailyRecord
		_ = json.Unmarshal(w.Body.Bytes(), &recAppended)
		assert.Len(t, recAppended.Tasks, 2)

		// Check if newTaskID was appended
		foundNewTask := false
		for _, tk := range recAppended.Tasks {
			if tk.TaskID == newTaskID {
				foundNewTask = true
			}
		}
		assert.True(t, foundNewTask)

		// 3. Past Date Not Found (Returns 404 Not Found)
		w = testutil.DoRequest(router, "GET", "/api/v1/daily/2020-01-01", nil)
		assert.Equal(t, http.StatusNotFound, w.Code)

		// 4. Invalid Date Format
		w = testutil.DoRequest(router, "GET", "/api/v1/daily/invalid-date", nil)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("PATCH /api/v1/daily/:date - Partial Updates", func(t *testing.T) {
		testutil.ClearDB(ctx)
		router := testutil.SetupRouter(testutil.DB)

		today := time.Now().In(time.FixedZone("JST", 9*3600)).Format("2006-01-02")
		// Generate record first
		_ = testutil.DoRequest(router, "GET", "/api/v1/daily/"+today, nil)

		// Patch context & journal
		patchBody := []byte(`{
			"context": {"mode": "Rest", "weather": "rainy"},
			"journal": {"oneLineReview": "Good day", "threeLineDiary": "Line 1\nLine 2"}
		}`)
		w := testutil.DoRequest(router, "PATCH", "/api/v1/daily/"+today, patchBody)
		assert.Equal(t, http.StatusOK, w.Code)

		var patched daily.DailyRecord
		err := json.Unmarshal(w.Body.Bytes(), &patched)
		require.NoError(t, err)
		assert.Equal(t, daily.ModeRest, patched.Context.Mode)
		assert.Equal(t, "rainy", patched.Context.Weather)
		assert.Equal(t, "Good day", patched.Journal.OneLineReview)

		// Patch tasks & habits array
		taskID := testutil.SeedTask(ctx, testutil.DB, task.Task{Title: "Run", Section: "exercise", Type: "boolean", Status: "active", StartDate: "2026-01-01"})
		arrayBody := []byte(`{
			"tasks": [{"taskId": "` + taskID.Hex() + `", "targetAmount": 1, "actualAmount": 1, "isCompleted": true}]
		}`)
		w = testutil.DoRequest(router, "PATCH", "/api/v1/daily/"+today, arrayBody)
		assert.Equal(t, http.StatusOK, w.Code)
		_ = json.Unmarshal(w.Body.Bytes(), &patched)
		require.Len(t, patched.Tasks, 1)
		assert.True(t, patched.Tasks[0].IsCompleted)

		// Empty Body
		w = testutil.DoRequest(router, "PATCH", "/api/v1/daily/"+today, []byte(`{}`))
		assert.Equal(t, http.StatusOK, w.Code)

		// Invalid Date
		w = testutil.DoRequest(router, "PATCH", "/api/v1/daily/bad-date", patchBody)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// Invalid JSON
		w = testutil.DoRequest(router, "PATCH", "/api/v1/daily/"+today, []byte(`{invalid}`))
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("GET /api/v1/daily/exists", func(t *testing.T) {
		testutil.ClearDB(ctx)
		router := testutil.SetupRouter(testutil.DB)

		testutil.SeedDailyRecord(ctx, testutil.DB, daily.DailyRecord{Date: "2026-05-10"})
		testutil.SeedDailyRecord(ctx, testutil.DB, daily.DailyRecord{Date: "2026-05-15"})

		// Success
		w := testutil.DoRequest(router, "GET", "/api/v1/daily/exists?year=2026&month=5", nil)
		assert.Equal(t, http.StatusOK, w.Code)
		var res struct {
			Dates []string `json:"dates"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &res)
		require.NoError(t, err)
		assert.Equal(t, []string{"2026-05-10", "2026-05-15"}, res.Dates)

		// Invalid Year
		w = testutil.DoRequest(router, "GET", "/api/v1/daily/exists?year=abc&month=5", nil)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// Invalid Month
		w = testutil.DoRequest(router, "GET", "/api/v1/daily/exists?year=2026&month=13", nil)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
