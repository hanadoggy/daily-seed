package habit_test

import (
	"context"
	"daily-seed/internal/common"
	"daily-seed/internal/habit"
	"daily-seed/internal/testutil"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestHabitHandler_Slice(t *testing.T) {
	ctx := context.Background()

	t.Run("GET /api/v1/habits - Empty and With Data", func(t *testing.T) {
		testutil.ClearDB(ctx)
		router := testutil.SetupRouter(testutil.DB)

		// Empty
		w := testutil.DoRequest(router, "GET", "/api/v1/habits", nil)
		assert.Equal(t, http.StatusOK, w.Code)
		var res []habit.Habit
		err := json.Unmarshal(w.Body.Bytes(), &res)
		require.NoError(t, err)
		assert.Empty(t, res)

		// With Data
		testutil.SeedHabit(ctx, testutil.DB, habit.Habit{Title: "Reading", Category: "Life", Status: "active"})
		w = testutil.DoRequest(router, "GET", "/api/v1/habits", nil)
		assert.Equal(t, http.StatusOK, w.Code)
		err = json.Unmarshal(w.Body.Bytes(), &res)
		require.NoError(t, err)
		assert.Len(t, res, 1)
	})

	t.Run("GET /api/v1/habits/:id", func(t *testing.T) {
		testutil.ClearDB(ctx)
		router := testutil.SetupRouter(testutil.DB)

		id := testutil.SeedHabit(ctx, testutil.DB, habit.Habit{Title: "Reading", Category: "Life", Status: "active"})

		// Success
		w := testutil.DoRequest(router, "GET", "/api/v1/habits/"+id.Hex(), nil)
		assert.Equal(t, http.StatusOK, w.Code)

		// NotFound
		nonExistID := primitive.NewObjectID().Hex()
		w = testutil.DoRequest(router, "GET", "/api/v1/habits/"+nonExistID, nil)
		assert.Equal(t, http.StatusNotFound, w.Code)
		var errRes common.ErrorResponse
		_ = json.Unmarshal(w.Body.Bytes(), &errRes)
		assert.Equal(t, "NOT_FOUND", errRes.Code)

		// InvalidID
		w = testutil.DoRequest(router, "GET", "/api/v1/habits/invalid-hex", nil)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("POST /api/v1/habits - Success and Validation Errors", func(t *testing.T) {
		testutil.ClearDB(ctx)
		router := testutil.SetupRouter(testutil.DB)

		// Success (and client status override check)
		validBody := []byte(`{"title":"Exercise","category":"Health","status":"archived"}`)
		w := testutil.DoRequest(router, "POST", "/api/v1/habits", validBody)
		assert.Equal(t, http.StatusCreated, w.Code)
		var created habit.Habit
		err := json.Unmarshal(w.Body.Bytes(), &created)
		require.NoError(t, err)
		assert.False(t, created.ID.IsZero())
		assert.Equal(t, "Exercise", created.Title)
		assert.Equal(t, "active", created.Status) // should override client status

		// Empty Title
		w = testutil.DoRequest(router, "POST", "/api/v1/habits", []byte(`{"title":"","category":"Health"}`))
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// Whitespace Title
		w = testutil.DoRequest(router, "POST", "/api/v1/habits", []byte(`{"title":"   ","category":"Health"}`))
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// Empty Category
		w = testutil.DoRequest(router, "POST", "/api/v1/habits", []byte(`{"title":"Exercise","category":""}`))
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// Invalid JSON Body
		w = testutil.DoRequest(router, "POST", "/api/v1/habits", []byte(`{invalid-json}`))
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// Empty Body
		w = testutil.DoRequest(router, "POST", "/api/v1/habits", nil)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("PUT /api/v1/habits/:id", func(t *testing.T) {
		testutil.ClearDB(ctx)
		router := testutil.SetupRouter(testutil.DB)

		id := testutil.SeedHabit(ctx, testutil.DB, habit.Habit{Title: "Reading", Category: "Life", Status: "archived"})

		// Success & Status Preservation
		w := testutil.DoRequest(router, "PUT", "/api/v1/habits/"+id.Hex(), []byte(`{"title":"Updated Reading","category":"Growth"}`))
		assert.Equal(t, http.StatusOK, w.Code)
		var updated habit.Habit
		_ = json.Unmarshal(w.Body.Bytes(), &updated)
		assert.Equal(t, "Updated Reading", updated.Title)
		assert.Equal(t, "archived", updated.Status) // preserved

		// NotFound
		nonExistID := primitive.NewObjectID().Hex()
		w = testutil.DoRequest(router, "PUT", "/api/v1/habits/"+nonExistID, []byte(`{"title":"Updated","category":"Growth"}`))
		assert.Equal(t, http.StatusNotFound, w.Code)

		// InvalidID Format
		w = testutil.DoRequest(router, "PUT", "/api/v1/habits/invalid-hex", []byte(`{"title":"Updated","category":"Growth"}`))
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// Validation Error (Empty Title)
		w = testutil.DoRequest(router, "PUT", "/api/v1/habits/"+id.Hex(), []byte(`{"title":"","category":"Growth"}`))
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// Invalid JSON
		w = testutil.DoRequest(router, "PUT", "/api/v1/habits/"+id.Hex(), []byte(`{bad-json}`))
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("DELETE /api/v1/habits/:id (Archive)", func(t *testing.T) {
		testutil.ClearDB(ctx)
		router := testutil.SetupRouter(testutil.DB)

		id := testutil.SeedHabit(ctx, testutil.DB, habit.Habit{Title: "Reading", Category: "Life", Status: "active"})

		// Success
		w := testutil.DoRequest(router, "DELETE", "/api/v1/habits/"+id.Hex(), nil)
		assert.Equal(t, http.StatusOK, w.Code)

		// Check DB status
		store := habit.NewHabitStore(testutil.DB)
		h, err := store.FindByID(ctx, id.Hex())
		require.NoError(t, err)
		assert.Equal(t, "archived", h.Status)

		// NotFound
		nonExistID := primitive.NewObjectID().Hex()
		w = testutil.DoRequest(router, "DELETE", "/api/v1/habits/"+nonExistID, nil)
		assert.Equal(t, http.StatusNotFound, w.Code)

		// InvalidID Format
		w = testutil.DoRequest(router, "DELETE", "/api/v1/habits/invalid-hex", nil)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
