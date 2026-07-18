package analytics

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"daily-seed/internal/common"
	"daily-seed/internal/daily"
	"daily-seed/internal/task"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupRouter(handler *AnalyticsHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	v1 := r.Group("/api/v1")
	handler.RegisterRoutes(v1)
	return r
}

func TestGetHeatmap_Success(t *testing.T) {
	mockDailyRepo := new(MockDailyRecordRepository)
	mockTaskRepo := new(MockTaskRepository)
	svc := NewAnalyticsService(mockDailyRepo, mockTaskRepo)
	handler := NewAnalyticsHandler(svc)
	router := setupRouter(handler)

	year := 2026
	startDate := fmt.Sprintf("%04d-01-01", year)
	endDate := fmt.Sprintf("%04d-12-31", year)

	mockDailyRepo.On("FindBetweenDates", mock.Anything, startDate, endDate).Return([]*daily.DailyRecord{}, nil)
	mockTaskRepo.On("FindAll", mock.Anything).Return([]task.Task{}, nil)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/analytics/heatmap?year=2026", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var res HeatmapResponse
	err := json.Unmarshal(w.Body.Bytes(), &res)
	assert.NoError(t, err)
	assert.Len(t, res.Days, 365)
}

func TestGetHeatmap_DefaultYear(t *testing.T) {
	mockDailyRepo := new(MockDailyRecordRepository)
	mockTaskRepo := new(MockTaskRepository)
	svc := NewAnalyticsService(mockDailyRepo, mockTaskRepo)
	handler := NewAnalyticsHandler(svc)
	router := setupRouter(handler)

	year := time.Now().Year()
	startDate := fmt.Sprintf("%04d-01-01", year)
	endDate := fmt.Sprintf("%04d-12-31", year)

	mockDailyRepo.On("FindBetweenDates", mock.Anything, startDate, endDate).Return([]*daily.DailyRecord{}, nil)
	mockTaskRepo.On("FindAll", mock.Anything).Return([]task.Task{}, nil)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/analytics/heatmap", nil) // no year param
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetHeatmap_InvalidYear(t *testing.T) {
	mockDailyRepo := new(MockDailyRecordRepository)
	mockTaskRepo := new(MockTaskRepository)
	svc := NewAnalyticsService(mockDailyRepo, mockTaskRepo)
	handler := NewAnalyticsHandler(svc)
	router := setupRouter(handler)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/analytics/heatmap?year=abc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	
	var res common.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &res)
	assert.NoError(t, err)
	assert.Equal(t, "INVALID_QUERY", res.Code)
}

func TestGetHeatmap_ServiceError(t *testing.T) {
	mockDailyRepo := new(MockDailyRecordRepository)
	mockTaskRepo := new(MockTaskRepository)
	svc := NewAnalyticsService(mockDailyRepo, mockTaskRepo)
	handler := NewAnalyticsHandler(svc)
	router := setupRouter(handler)

	year := 2026
	startDate := fmt.Sprintf("%04d-01-01", year)
	endDate := fmt.Sprintf("%04d-12-31", year)

	mockDailyRepo.On("FindBetweenDates", mock.Anything, startDate, endDate).Return(nil, fmt.Errorf("db error"))

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/analytics/heatmap?year=2026", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var res common.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &res)
	assert.NoError(t, err)
	assert.Equal(t, "INTERNAL_ERROR", res.Code)
}
