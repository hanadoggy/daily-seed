package daily_test

import (
	"bytes"
	"daily-seed/internal/daily"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupDailyRouter(svc daily.DailyService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	rg := r.Group("/api/v1")
	handler := daily.NewDailyHandler(svc)
	handler.RegisterRoutes(rg)
	return r
}

func TestDailyHandler_GetDailyRecord(t *testing.T) {
	tests := []struct {
		name           string
		date           string
		mockSetup      func(*daily.MockDailyService)
		expectedStatus int
		expectedCode   string
	}{
		{
			name: "Pass: successful get",
			date: "2026-07-18",
			mockSetup: func(m *daily.MockDailyService) {
				m.On("GetDailyRecord", mock.Anything, "2026-07-18").Return(&daily.DailyRecord{
					Date: "2026-07-18",
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCode:   "",
		},
		{
			name:           "Fail: invalid date format",
			date:           "2026-07-1",
			mockSetup:      func(m *daily.MockDailyService) {},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "INVALID_DATE",
		},
		{
			name: "Fail: internal error",
			date: "2026-07-18",
			mockSetup: func(m *daily.MockDailyService) {
				m.On("GetDailyRecord", mock.Anything, "2026-07-18").Return((*daily.DailyRecord)(nil), errors.New("db error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedCode:   "INTERNAL_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(daily.MockDailyService)
			tt.mockSetup(mockSvc)
			r := setupDailyRouter(mockSvc)

			req, _ := http.NewRequest(http.MethodGet, "/api/v1/daily/"+tt.date, nil)
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

func TestDailyHandler_UpdateDailyRecord(t *testing.T) {
	tests := []struct {
		name           string
		date           string
		body           interface{}
		mockSetup      func(*daily.MockDailyService)
		expectedStatus int
		expectedCode   string
	}{
		{
			name: "Pass: successful update",
			date: "2026-07-18",
			body: map[string]interface{}{
				"context": map[string]string{
					"mode": "Growth",
				},
			},
			mockSetup: func(m *daily.MockDailyService) {
				m.On("UpdateDailyRecord", mock.Anything, "2026-07-18", mock.AnythingOfType("*daily.UpdateDailyRecordRequest")).Return(&daily.DailyRecord{
					Date: "2026-07-18",
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCode:   "",
		},
		{
			name:           "Fail: invalid date format",
			date:           "2026-07-1",
			body:           map[string]interface{}{},
			mockSetup:      func(m *daily.MockDailyService) {},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "INVALID_DATE",
		},
		{
			name:           "Fail: invalid JSON",
			date:           "2026-07-18",
			body:           "invalid json",
			mockSetup:      func(m *daily.MockDailyService) {},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "INVALID_BODY",
		},
		{
			name: "Fail: internal error",
			date: "2026-07-18",
			body: map[string]interface{}{},
			mockSetup: func(m *daily.MockDailyService) {
				m.On("UpdateDailyRecord", mock.Anything, "2026-07-18", mock.AnythingOfType("*daily.UpdateDailyRecordRequest")).Return((*daily.DailyRecord)(nil), errors.New("db error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedCode:   "INTERNAL_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(daily.MockDailyService)
			tt.mockSetup(mockSvc)
			r := setupDailyRouter(mockSvc)

			var bodyBytes []byte
			if str, ok := tt.body.(string); ok {
				bodyBytes = []byte(str)
			} else {
				bodyBytes, _ = json.Marshal(tt.body)
			}

			req, _ := http.NewRequest(http.MethodPatch, "/api/v1/daily/"+tt.date, bytes.NewBuffer(bodyBytes))
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
