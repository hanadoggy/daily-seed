package task_test

import (
	"context"
	"daily-seed/internal/daily"
	"daily-seed/internal/task"
	"daily-seed/internal/testutil"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestTaskHandler_Slice(t *testing.T) {
	ctx := context.Background()

	t.Run("POST /api/v1/tasks - Success and Validation Errors", func(t *testing.T) {
		testutil.ClearDB(ctx)
		router := testutil.SetupRouter(testutil.DB)

		// Quantitative Success with Custom Unit
		quantBody := []byte(`{
			"section": "japanese",
			"title": "Kanji",
			"type": "quantitative",
			"unit": "pages",
			"metrics": {"dailyTarget": 10, "totalTarget": 100},
			"startDate": "2026-01-01"
		}`)
		w := testutil.DoRequest(router, "POST", "/api/v1/tasks", quantBody)
		assert.Equal(t, http.StatusCreated, w.Code)
		var created task.Task
		err := json.Unmarshal(w.Body.Bytes(), &created)
		require.NoError(t, err)
		assert.Equal(t, 10, created.Metrics.DailyTarget)
		assert.Equal(t, "pages", created.Unit)
		assert.Equal(t, []string{"sunny", "rainy"}, created.Conditions.Weather) // defaults applied

		// Boolean Success (dailyTarget override to 1, default unit "units")
		boolBody := []byte(`{
			"section": "exercise",
			"title": "Pushups",
			"type": "boolean",
			"metrics": {"dailyTarget": 50, "totalTarget": 0},
			"startDate": "2026-01-01"
		}`)
		w = testutil.DoRequest(router, "POST", "/api/v1/tasks", boolBody)
		assert.Equal(t, http.StatusCreated, w.Code)
		var boolCreated task.Task
		_ = json.Unmarshal(w.Body.Bytes(), &boolCreated)
		assert.Equal(t, 1, boolCreated.Metrics.DailyTarget)
		assert.Equal(t, "units", boolCreated.Unit)

		// Empty Title
		w = testutil.DoRequest(router, "POST", "/api/v1/tasks", []byte(`{"section":"dev","title":"","type":"boolean","startDate":"2026-01-01"}`))
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// Invalid Section
		w = testutil.DoRequest(router, "POST", "/api/v1/tasks", []byte(`{"section":"cooking","title":"Test","type":"boolean","startDate":"2026-01-01"}`))
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// Invalid Type
		w = testutil.DoRequest(router, "POST", "/api/v1/tasks", []byte(`{"section":"dev","title":"Test","type":"percentage","startDate":"2026-01-01"}`))
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// Negative/Zero Daily Target for Quantitative
		w = testutil.DoRequest(router, "POST", "/api/v1/tasks", []byte(`{"section":"dev","title":"Test","type":"quantitative","metrics":{"dailyTarget":0},"startDate":"2026-01-01"}`))
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// Negative Total Target
		w = testutil.DoRequest(router, "POST", "/api/v1/tasks", []byte(`{"section":"dev","title":"Test","type":"boolean","metrics":{"totalTarget":-1},"startDate":"2026-01-01"}`))
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// Missing StartDate
		w = testutil.DoRequest(router, "POST", "/api/v1/tasks", []byte(`{"section":"dev","title":"Test","type":"boolean"}`))
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// Invalid Weather
		w = testutil.DoRequest(router, "POST", "/api/v1/tasks", []byte(`{"section":"dev","title":"Test","type":"boolean","startDate":"2026-01-01","conditions":{"weather":["cloudy"],"mode":["Growth"]}}`))
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// Invalid Mode
		w = testutil.DoRequest(router, "POST", "/api/v1/tasks", []byte(`{"section":"dev","title":"Test","type":"boolean","startDate":"2026-01-01","conditions":{"weather":["sunny"],"mode":["Holiday"]}}`))
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// Invalid JSON
		w = testutil.DoRequest(router, "POST", "/api/v1/tasks", []byte(`{bad}`))
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("GET /api/v1/tasks & /api/v1/tasks/:id", func(t *testing.T) {
		testutil.ClearDB(ctx)
		router := testutil.SetupRouter(testutil.DB)

		id := testutil.SeedTask(ctx, testutil.DB, task.Task{Section: "dev", Title: "Coding", Type: "boolean", Status: "active", StartDate: "2026-01-01"})

		// List
		w := testutil.DoRequest(router, "GET", "/api/v1/tasks", nil)
		assert.Equal(t, http.StatusOK, w.Code)

		// Get Exists
		w = testutil.DoRequest(router, "GET", "/api/v1/tasks/"+id.Hex(), nil)
		assert.Equal(t, http.StatusOK, w.Code)

		// Get NotFound
		nonExistID := primitive.NewObjectID().Hex()
		w = testutil.DoRequest(router, "GET", "/api/v1/tasks/"+nonExistID, nil)
		assert.Equal(t, http.StatusNotFound, w.Code)

		// Get InvalidID
		w = testutil.DoRequest(router, "GET", "/api/v1/tasks/invalid-hex", nil)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("PUT /api/v1/tasks/:id", func(t *testing.T) {
		testutil.ClearDB(ctx)
		router := testutil.SetupRouter(testutil.DB)

		id := testutil.SeedTask(ctx, testutil.DB, task.Task{
			Section:   "dev",
			Title:     "Old Title",
			Type:      "quantitative",
			Metrics:   task.TaskMetrics{DailyTarget: 10},
			Status:    "active",
			StartDate: "2026-01-01",
			Conditions: task.TaskConditions{
				Weather: []string{"sunny"},
				Mode:    []string{"Growth"},
			},
		})

		// Success
		updateBody := []byte(`{
			"section": "dev",
			"title": "New Title",
			"type": "quantitative",
			"metrics": {"dailyTarget": 20},
			"startDate": "2026-01-01",
			"conditions": {"weather": ["sunny"], "mode": ["Growth"]}
		}`)
		w := testutil.DoRequest(router, "PUT", "/api/v1/tasks/"+id.Hex(), updateBody)
		assert.Equal(t, http.StatusOK, w.Code)

		// Type Change Attempt (Forbidden)
		typeChangeBody := []byte(`{
			"section": "dev",
			"title": "New Title",
			"type": "boolean",
			"startDate": "2026-01-01",
			"conditions": {"weather": ["sunny"], "mode": ["Growth"]}
		}`)
		w = testutil.DoRequest(router, "PUT", "/api/v1/tasks/"+id.Hex(), typeChangeBody)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// Archived Task Update Attempt
		archivedID := testutil.SeedTask(ctx, testutil.DB, task.Task{
			Section: "dev", Title: "Archived", Type: "boolean", Status: "archived", StartDate: "2026-01-01",
			Conditions: task.TaskConditions{Weather: []string{"sunny"}, Mode: []string{"Growth"}},
		})
		w = testutil.DoRequest(router, "PUT", "/api/v1/tasks/"+archivedID.Hex(), updateBody)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// StartDate Delayed -> Cleans past daily records
		testutil.SeedDailyRecord(ctx, testutil.DB, daily.DailyRecord{
			Date: "2026-01-02",
			Tasks: []daily.TaskEntry{
				{TaskID: id, TargetAmount: 10, ActualAmount: 5},
			},
		})
		delayBody := []byte(`{
			"section": "dev",
			"title": "New Title",
			"type": "quantitative",
			"metrics": {"dailyTarget": 20},
			"startDate": "2026-01-05",
			"conditions": {"weather": ["sunny"], "mode": ["Growth"]}
		}`)
		w = testutil.DoRequest(router, "PUT", "/api/v1/tasks/"+id.Hex(), delayBody)
		assert.Equal(t, http.StatusOK, w.Code)

		// Verify task removed from record 2026-01-02
		dailyStore := daily.NewDailyStore(testutil.DB)
		rec, _ := dailyStore.FindByDate(ctx, "2026-01-02")
		assert.Empty(t, rec.Tasks)
	})

	t.Run("DELETE /api/v1/tasks/:id (Archive)", func(t *testing.T) {
		testutil.ClearDB(ctx)
		router := testutil.SetupRouter(testutil.DB)

		id := testutil.SeedTask(ctx, testutil.DB, task.Task{Section: "dev", Title: "Coding", Type: "boolean", Status: "active", StartDate: "2026-01-01"})

		// Success
		w := testutil.DoRequest(router, "DELETE", "/api/v1/tasks/"+id.Hex(), nil)
		assert.Equal(t, http.StatusOK, w.Code)

		// Already Archived
		w = testutil.DoRequest(router, "DELETE", "/api/v1/tasks/"+id.Hex(), nil)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// NotFound
		nonExistID := primitive.NewObjectID().Hex()
		w = testutil.DoRequest(router, "DELETE", "/api/v1/tasks/"+nonExistID, nil)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("GET /api/v1/tasks/progress", func(t *testing.T) {
		testutil.ClearDB(ctx)
		router := testutil.SetupRouter(testutil.DB)

		// Active task 1 (totalTarget 100)
		id1 := testutil.SeedTask(ctx, testutil.DB, task.Task{
			Title: "Read Books", Section: "self_dev", Type: "quantitative",
			Metrics: task.TaskMetrics{DailyTarget: 10, TotalTarget: 100}, Status: "active", StartDate: "2026-01-01",
		})
		// Active task 2 (totalTarget 0 - unlimited)
		id2 := testutil.SeedTask(ctx, testutil.DB, task.Task{
			Title: "Meditate", Section: "self_dev", Type: "boolean",
			Metrics: task.TaskMetrics{DailyTarget: 1, TotalTarget: 0}, Status: "active", StartDate: "2026-01-01",
		})
		// Archived task (should be excluded)
		testutil.SeedTask(ctx, testutil.DB, task.Task{
			Title: "Old Task", Section: "dev", Type: "boolean", Status: "archived", StartDate: "2026-01-01",
		})

		// Seed daily record progress for task 1
		testutil.SeedDailyRecord(ctx, testutil.DB, daily.DailyRecord{
			Date: "2026-01-01",
			Tasks: []daily.TaskEntry{
				{TaskID: id1, ActualAmount: 30},
				{TaskID: id2, ActualAmount: 1},
			},
		})
		testutil.SeedDailyRecord(ctx, testutil.DB, daily.DailyRecord{
			Date: "2026-01-02",
			Tasks: []daily.TaskEntry{
				{TaskID: id1, ActualAmount: 20},
			},
		})

		w := testutil.DoRequest(router, "GET", "/api/v1/tasks/progress", nil)
		assert.Equal(t, http.StatusOK, w.Code)

		var progress []task.TaskProgress
		err := json.Unmarshal(w.Body.Bytes(), &progress)
		require.NoError(t, err)
		require.Len(t, progress, 2)

		for _, p := range progress {
			if p.TaskID == id1 {
				assert.Equal(t, 50, p.TotalCompleted)
				assert.InDelta(t, 50.0, p.Percentage, 0.01)
			} else if p.TaskID == id2 {
				assert.Equal(t, 1, p.TotalCompleted)
				assert.Equal(t, 0.0, p.Percentage) // 0 totalTarget -> 0 percentage
			}
		}
	})

	t.Run("POST /api/v1/tasks/:id/migrate", func(t *testing.T) {
		testutil.ClearDB(ctx)
		router := testutil.SetupRouter(testutil.DB)

		// Target 100 task
		id := testutil.SeedTask(ctx, testutil.DB, task.Task{
			Title: "Finish Course", Section: "dev", Type: "quantitative",
			Metrics: task.TaskMetrics{DailyTarget: 10, TotalTarget: 100}, Status: "active", StartDate: "2026-01-01",
			Conditions: task.TaskConditions{Weather: []string{"sunny"}, Mode: []string{"Growth"}},
		})

		// 1. Attempt migrate before reaching target -> Fail
		migrateReq := []byte(`{"completionDate":"2026-01-10"}`)
		w := testutil.DoRequest(router, "POST", fmt.Sprintf("/api/v1/tasks/%s/migrate", id.Hex()), migrateReq)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// 2. Add progress (100)
		testutil.SeedDailyRecord(ctx, testutil.DB, daily.DailyRecord{
			Date: "2026-01-10",
			Tasks: []daily.TaskEntry{
				{TaskID: id, ActualAmount: 100},
			},
		})

		// 3. Invalid Completion Date Format -> Fail
		w = testutil.DoRequest(router, "POST", fmt.Sprintf("/api/v1/tasks/%s/migrate", id.Hex()), []byte(`{"completionDate":"2026/01/10"}`))
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// 4. Missing Completion Date -> Fail
		w = testutil.DoRequest(router, "POST", fmt.Sprintf("/api/v1/tasks/%s/migrate", id.Hex()), []byte(`{}`))
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// 5. Successful Migrate
		w = testutil.DoRequest(router, "POST", fmt.Sprintf("/api/v1/tasks/%s/migrate", id.Hex()), migrateReq)
		assert.Equal(t, http.StatusOK, w.Code)

		var result task.MigrationResult
		err := json.Unmarshal(w.Body.Bytes(), &result)
		require.NoError(t, err)

		assert.Equal(t, "archived", result.ArchivedTask.Status)
		assert.Equal(t, "2026-01-10", result.ArchivedTask.EndDate)

		assert.Equal(t, "active", result.NewTask.Status)
		assert.Equal(t, "2026-01-11", result.NewTask.StartDate) // completionDate + 1 day
		assert.Equal(t, "Finish Course", result.NewTask.Title)
	})
}
