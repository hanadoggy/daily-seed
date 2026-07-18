package task_test

import (
	"bytes"
	"daily-seed/internal/task"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTaskRouter(svc task.TaskService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	rg := r.Group("/api/v1")
	handler := task.NewTaskHandler(svc)
	handler.RegisterRoutes(rg)
	return r
}

func TestTaskHandler_Create(t *testing.T) {
	tests := []struct {
		name           string
		body           interface{}
		mockSetup      func(*task.MockTaskService)
		expectedStatus int
		expectedCode   string
	}{
		{
			name: "Pass: successful creation",
			body: map[string]interface{}{
				"title":     "New Task",
				"section":   "dev",
				"type":      "boolean",
				"startDate": "2026-07-18",
			},
			mockSetup: func(m *task.MockTaskService) {
				m.On("Create", mock.Anything, mock.AnythingOfType("*task.Task")).Return(&task.Task{
					ID:    testOID("000000000000000000000001"),
					Title: "New Task",
				}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedCode:   "",
		},
		{
			name:           "Fail: invalid JSON",
			body:           "invalid json",
			mockSetup:      func(m *task.MockTaskService) {},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "INVALID_BODY",
		},
		{
			name: "Fail: validation error",
			body: map[string]interface{}{
				"title": "",
			},
			mockSetup: func(m *task.MockTaskService) {
				m.On("Create", mock.Anything, mock.AnythingOfType("*task.Task")).Return((*task.Task)(nil), errors.New("title is required"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "VALIDATION_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(task.MockTaskService)
			tt.mockSetup(mockSvc)
			r := setupTaskRouter(mockSvc)

			var bodyBytes []byte
			if str, ok := tt.body.(string); ok {
				bodyBytes = []byte(str)
			} else {
				bodyBytes, _ = json.Marshal(tt.body)
			}

			req, _ := http.NewRequest(http.MethodPost, "/api/v1/tasks", bytes.NewBuffer(bodyBytes))
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

func TestTaskHandler_Get(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		mockSetup      func(*task.MockTaskService)
		expectedStatus int
		expectedCode   string
	}{
		{
			name: "Pass: successful get",
			id:   "000000000000000000000001",
			mockSetup: func(m *task.MockTaskService) {
				m.On("Get", mock.Anything, "000000000000000000000001").Return(&task.Task{
					ID:    testOID("000000000000000000000001"),
					Title: "Task 1",
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCode:   "",
		},
		{
			name: "Fail: not found",
			id:   "000000000000000000000002",
			mockSetup: func(m *task.MockTaskService) {
				m.On("Get", mock.Anything, "000000000000000000000002").Return((*task.Task)(nil), errors.New("task not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedCode:   "NOT_FOUND",
		},
		{
			name: "Fail: internal error",
			id:   "000000000000000000000001",
			mockSetup: func(m *task.MockTaskService) {
				m.On("Get", mock.Anything, "000000000000000000000001").Return((*task.Task)(nil), errors.New("db error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedCode:   "INTERNAL_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(task.MockTaskService)
			tt.mockSetup(mockSvc)
			r := setupTaskRouter(mockSvc)

			req, _ := http.NewRequest(http.MethodGet, "/api/v1/tasks/"+tt.id, nil)
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

func TestTaskHandler_List(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*task.MockTaskService)
		expectedStatus int
		expectedCode   string
	}{
		{
			name: "Pass: successful list",
			mockSetup: func(m *task.MockTaskService) {
				m.On("List", mock.Anything).Return([]task.Task{
					{ID: testOID("000000000000000000000001"), Title: "Task 1"},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCode:   "",
		},
		{
			name: "Fail: internal error",
			mockSetup: func(m *task.MockTaskService) {
				m.On("List", mock.Anything).Return(([]task.Task)(nil), errors.New("db error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedCode:   "INTERNAL_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(task.MockTaskService)
			tt.mockSetup(mockSvc)
			r := setupTaskRouter(mockSvc)

			req, _ := http.NewRequest(http.MethodGet, "/api/v1/tasks", nil)
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

func TestTaskHandler_Update(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		body           interface{}
		mockSetup      func(*task.MockTaskService)
		expectedStatus int
		expectedCode   string
	}{
		{
			name: "Pass: successful update",
			id:   "000000000000000000000001",
			body: map[string]interface{}{
				"title": "Updated Task",
			},
			mockSetup: func(m *task.MockTaskService) {
				m.On("Update", mock.Anything, mock.AnythingOfType("*task.Task")).Return(&task.Task{
					ID:    testOID("000000000000000000000001"),
					Title: "Updated Task",
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCode:   "",
		},
		{
			name: "Fail: invalid ID format",
			id:   "invalid-id",
			body: map[string]interface{}{
				"title": "Updated Task",
			},
			mockSetup:      func(m *task.MockTaskService) {},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "INVALID_ID",
		},
		{
			name:           "Fail: invalid JSON",
			id:             "000000000000000000000001",
			body:           "invalid json",
			mockSetup:      func(m *task.MockTaskService) {},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "INVALID_BODY",
		},
		{
			name: "Fail: not found",
			id:   "000000000000000000000001",
			body: map[string]interface{}{
				"title": "Updated Task",
			},
			mockSetup: func(m *task.MockTaskService) {
				m.On("Update", mock.Anything, mock.AnythingOfType("*task.Task")).Return((*task.Task)(nil), errors.New("task not found"))
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
			mockSetup: func(m *task.MockTaskService) {
				m.On("Update", mock.Anything, mock.AnythingOfType("*task.Task")).Return((*task.Task)(nil), errors.New("title is required"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "VALIDATION_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(task.MockTaskService)
			tt.mockSetup(mockSvc)
			r := setupTaskRouter(mockSvc)

			var bodyBytes []byte
			if str, ok := tt.body.(string); ok {
				bodyBytes = []byte(str)
			} else {
				bodyBytes, _ = json.Marshal(tt.body)
			}

			req, _ := http.NewRequest(http.MethodPut, "/api/v1/tasks/"+tt.id, bytes.NewBuffer(bodyBytes))
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

func TestTaskHandler_Archive(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		mockSetup      func(*task.MockTaskService)
		expectedStatus int
		expectedCode   string
	}{
		{
			name: "Pass: successful archive",
			id:   "000000000000000000000001",
			mockSetup: func(m *task.MockTaskService) {
				m.On("Archive", mock.Anything, "000000000000000000000001").Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedCode:   "",
		},
		{
			name: "Fail: not found",
			id:   "000000000000000000000001",
			mockSetup: func(m *task.MockTaskService) {
				m.On("Archive", mock.Anything, "000000000000000000000001").Return(errors.New("task not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedCode:   "NOT_FOUND",
		},
		{
			name: "Fail: internal error",
			id:   "000000000000000000000001",
			mockSetup: func(m *task.MockTaskService) {
				m.On("Archive", mock.Anything, "000000000000000000000001").Return(errors.New("db error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedCode:   "INTERNAL_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(task.MockTaskService)
			tt.mockSetup(mockSvc)
			r := setupTaskRouter(mockSvc)

			req, _ := http.NewRequest(http.MethodDelete, "/api/v1/tasks/"+tt.id, nil)
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

func TestTaskHandler_GetProgress(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*task.MockTaskService)
		expectedStatus int
		expectedCode   string
	}{
		{
			name: "Pass: successful get progress",
			mockSetup: func(m *task.MockTaskService) {
				m.On("GetProgressForActiveTasks", mock.Anything).Return([]task.TaskProgress{
					{TaskID: testOID("000000000000000000000001"), Title: "Task 1", TotalTarget: 10, TotalCompleted: 5, Percentage: 50},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCode:   "",
		},
		{
			name: "Fail: internal error",
			mockSetup: func(m *task.MockTaskService) {
				m.On("GetProgressForActiveTasks", mock.Anything).Return(([]task.TaskProgress)(nil), errors.New("db error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedCode:   "INTERNAL_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(task.MockTaskService)
			tt.mockSetup(mockSvc)
			r := setupTaskRouter(mockSvc)

			req, _ := http.NewRequest(http.MethodGet, "/api/v1/tasks/progress", nil)
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

func TestTaskHandler_Migrate(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		body           interface{}
		mockSetup      func(*task.MockTaskService)
		expectedStatus int
		expectedCode   string
	}{
		{
			name: "Pass: successful migrate",
			id:   "000000000000000000000001",
			body: map[string]interface{}{
				"completionDate": "2026-07-18",
			},
			mockSetup: func(m *task.MockTaskService) {
				m.On("MigrateTask", mock.Anything, "000000000000000000000001", mock.AnythingOfType("task.MigrateTaskRequest")).Return(&task.MigrationResult{
					ArchivedTask: task.Task{ID: testOID("000000000000000000000001"), Status: "archived"},
					NewTask:      task.Task{ID: testOID("000000000000000000000002"), Status: "active"},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCode:   "",
		},
		{
			name:           "Fail: invalid JSON",
			id:             "000000000000000000000001",
			body:           "invalid json",
			mockSetup:      func(m *task.MockTaskService) {},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "INVALID_BODY",
		},
		{
			name: "Fail: not found",
			id:   "000000000000000000000001",
			body: map[string]interface{}{
				"completionDate": "2026-07-18",
			},
			mockSetup: func(m *task.MockTaskService) {
				m.On("MigrateTask", mock.Anything, "000000000000000000000001", mock.Anything).Return((*task.MigrationResult)(nil), errors.New("task not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedCode:   "NOT_FOUND",
		},
		{
			name: "Fail: migration not allowed (progress not reached)",
			id:   "000000000000000000000001",
			body: map[string]interface{}{
				"completionDate": "2026-07-18",
			},
			mockSetup: func(m *task.MockTaskService) {
				m.On("MigrateTask", mock.Anything, "000000000000000000000001", mock.Anything).Return((*task.MigrationResult)(nil), errors.New("task progress (5/10) has not reached the target"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "MIGRATION_NOT_ALLOWED",
		},
		{
			name: "Fail: internal error",
			id:   "000000000000000000000001",
			body: map[string]interface{}{
				"completionDate": "2026-07-18",
			},
			mockSetup: func(m *task.MockTaskService) {
				m.On("MigrateTask", mock.Anything, "000000000000000000000001", mock.Anything).Return((*task.MigrationResult)(nil), errors.New("db error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedCode:   "INTERNAL_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(task.MockTaskService)
			tt.mockSetup(mockSvc)
			r := setupTaskRouter(mockSvc)

			var bodyBytes []byte
			if str, ok := tt.body.(string); ok {
				bodyBytes = []byte(str)
			} else {
				bodyBytes, _ = json.Marshal(tt.body)
			}

			req, _ := http.NewRequest(http.MethodPost, "/api/v1/tasks/"+tt.id+"/migrate", bytes.NewBuffer(bodyBytes))
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
