package habit_test

import (
	"bytes"
	"daily-seed/internal/habit"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func setupHabitRouter(svc habit.HabitService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	rg := r.Group("/api/v1")
	handler := habit.NewHabitHandler(svc)
	handler.RegisterRoutes(rg)
	return r
}

func testOID(hex string) primitive.ObjectID {
	oid, _ := primitive.ObjectIDFromHex(hex)
	return oid
}

func TestHabitHandler_Create(t *testing.T) {
	tests := []struct {
		name           string
		body           interface{}
		mockSetup      func(*habit.MockHabitService)
		expectedStatus int
		expectedCode   string
	}{
		{
			name: "Pass: successful creation",
			body: map[string]interface{}{
				"title":    "New Habit",
				"category": "Health",
			},
			mockSetup: func(m *habit.MockHabitService) {
				m.On("Create", mock.Anything, mock.AnythingOfType("*habit.Habit")).Return(&habit.Habit{
					ID:    testOID("000000000000000000000001"),
					Title: "New Habit",
				}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedCode:   "",
		},
		{
			name:           "Fail: invalid JSON",
			body:           "invalid json",
			mockSetup:      func(m *habit.MockHabitService) {},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "INVALID_BODY",
		},
		{
			name: "Fail: validation error",
			body: map[string]interface{}{
				"title": "",
			},
			mockSetup: func(m *habit.MockHabitService) {
				m.On("Create", mock.Anything, mock.AnythingOfType("*habit.Habit")).Return((*habit.Habit)(nil), errors.New("title is required"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "VALIDATION_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(habit.MockHabitService)
			tt.mockSetup(mockSvc)
			r := setupHabitRouter(mockSvc)

			var bodyBytes []byte
			if str, ok := tt.body.(string); ok {
				bodyBytes = []byte(str)
			} else {
				bodyBytes, _ = json.Marshal(tt.body)
			}

			req, _ := http.NewRequest(http.MethodPost, "/api/v1/habits", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedCode != "" {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCode, response["code"])
			}
			mockSvc.AssertExpectations(t)
		})
	}
}

func TestHabitHandler_Get(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		mockSetup      func(*habit.MockHabitService)
		expectedStatus int
		expectedCode   string
	}{
		{
			name: "Pass: successful get",
			id:   "000000000000000000000001",
			mockSetup: func(m *habit.MockHabitService) {
				m.On("Get", mock.Anything, "000000000000000000000001").Return(&habit.Habit{
					ID:    testOID("000000000000000000000001"),
					Title: "Habit 1",
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCode:   "",
		},
		{
			name: "Fail: not found",
			id:   "000000000000000000000002",
			mockSetup: func(m *habit.MockHabitService) {
				m.On("Get", mock.Anything, "000000000000000000000002").Return((*habit.Habit)(nil), errors.New("habit not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedCode:   "NOT_FOUND",
		},
		{
			name: "Fail: internal error",
			id:   "000000000000000000000001",
			mockSetup: func(m *habit.MockHabitService) {
				m.On("Get", mock.Anything, "000000000000000000000001").Return((*habit.Habit)(nil), errors.New("db error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedCode:   "INTERNAL_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(habit.MockHabitService)
			tt.mockSetup(mockSvc)
			r := setupHabitRouter(mockSvc)

			req, _ := http.NewRequest(http.MethodGet, "/api/v1/habits/"+tt.id, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedCode != "" {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCode, response["code"])
			}
			mockSvc.AssertExpectations(t)
		})
	}
}

func TestHabitHandler_List(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*habit.MockHabitService)
		expectedStatus int
		expectedCode   string
	}{
		{
			name: "Pass: successful list",
			mockSetup: func(m *habit.MockHabitService) {
				m.On("List", mock.Anything).Return([]habit.Habit{
					{ID: testOID("000000000000000000000001"), Title: "Habit 1"},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCode:   "",
		},
		{
			name: "Fail: internal error",
			mockSetup: func(m *habit.MockHabitService) {
				m.On("List", mock.Anything).Return(([]habit.Habit)(nil), errors.New("db error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedCode:   "INTERNAL_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(habit.MockHabitService)
			tt.mockSetup(mockSvc)
			r := setupHabitRouter(mockSvc)

			req, _ := http.NewRequest(http.MethodGet, "/api/v1/habits", nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedCode != "" {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCode, response["code"])
			}
			mockSvc.AssertExpectations(t)
		})
	}
}

func TestHabitHandler_Update(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		body           interface{}
		mockSetup      func(*habit.MockHabitService)
		expectedStatus int
		expectedCode   string
	}{
		{
			name: "Pass: successful update",
			id:   "000000000000000000000001",
			body: map[string]interface{}{
				"title": "Updated Habit",
			},
			mockSetup: func(m *habit.MockHabitService) {
				m.On("Update", mock.Anything, mock.AnythingOfType("*habit.Habit")).Return(&habit.Habit{
					ID:    testOID("000000000000000000000001"),
					Title: "Updated Habit",
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCode:   "",
		},
		{
			name: "Fail: invalid ID format",
			id:   "invalid-id",
			body: map[string]interface{}{
				"title": "Updated Habit",
			},
			mockSetup:      func(m *habit.MockHabitService) {},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "INVALID_ID",
		},
		{
			name:           "Fail: invalid JSON",
			id:             "000000000000000000000001",
			body:           "invalid json",
			mockSetup:      func(m *habit.MockHabitService) {},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "INVALID_BODY",
		},
		{
			name: "Fail: not found",
			id:   "000000000000000000000001",
			body: map[string]interface{}{
				"title": "Updated Habit",
			},
			mockSetup: func(m *habit.MockHabitService) {
				m.On("Update", mock.Anything, mock.AnythingOfType("*habit.Habit")).Return((*habit.Habit)(nil), errors.New("habit not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedCode:   "NOT_FOUND",
		},
		{
			name: "Fail: validation error",
			id:   "000000000000000000000001",
			body: map[string]interface{}{
				"title": "",
			},
			mockSetup: func(m *habit.MockHabitService) {
				m.On("Update", mock.Anything, mock.AnythingOfType("*habit.Habit")).Return((*habit.Habit)(nil), errors.New("title is required"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "VALIDATION_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(habit.MockHabitService)
			tt.mockSetup(mockSvc)
			r := setupHabitRouter(mockSvc)

			var bodyBytes []byte
			if str, ok := tt.body.(string); ok {
				bodyBytes = []byte(str)
			} else {
				bodyBytes, _ = json.Marshal(tt.body)
			}

			req, _ := http.NewRequest(http.MethodPut, "/api/v1/habits/"+tt.id, bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedCode != "" {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCode, response["code"])
			}
			mockSvc.AssertExpectations(t)
		})
	}
}

func TestHabitHandler_Archive(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		mockSetup      func(*habit.MockHabitService)
		expectedStatus int
		expectedCode   string
	}{
		{
			name: "Pass: successful archive",
			id:   "000000000000000000000001",
			mockSetup: func(m *habit.MockHabitService) {
				m.On("Archive", mock.Anything, "000000000000000000000001").Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedCode:   "",
		},
		{
			name: "Fail: not found",
			id:   "000000000000000000000001",
			mockSetup: func(m *habit.MockHabitService) {
				m.On("Archive", mock.Anything, "000000000000000000000001").Return(errors.New("habit not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedCode:   "NOT_FOUND",
		},
		{
			name: "Fail: internal error",
			id:   "000000000000000000000001",
			mockSetup: func(m *habit.MockHabitService) {
				m.On("Archive", mock.Anything, "000000000000000000000001").Return(errors.New("db error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedCode:   "INTERNAL_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(habit.MockHabitService)
			tt.mockSetup(mockSvc)
			r := setupHabitRouter(mockSvc)

			req, _ := http.NewRequest(http.MethodDelete, "/api/v1/habits/"+tt.id, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedCode != "" {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCode, response["code"])
			}
			mockSvc.AssertExpectations(t)
		})
	}
}
