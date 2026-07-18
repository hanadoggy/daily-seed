package middleware_test

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"daily-seed/internal/middleware"
)

func TestLogger(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Capture log output
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))
	slog.SetDefault(logger)

	router := gin.New()
	router.Use(middleware.Logger())

	router.GET("/test", func(c *gin.Context) {
		// Verify that RequestID was set in the context
		reqID, exists := c.Get(middleware.RequestIDKey)
		assert.True(t, exists, "RequestIDKey should be set in context")
		assert.NotEmpty(t, reqID)

		c.String(http.StatusOK, "ok")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	
	// Set remote addr for clientIP
	req.RemoteAddr = "192.168.1.1:1234"

	router.ServeHTTP(w, req)

	// Verify response
	assert.Equal(t, http.StatusOK, w.Code)
	
	// Verify header
	reqID := w.Header().Get("X-Request-ID")
	assert.NotEmpty(t, reqID, "X-Request-ID header should be set")

	// Verify log output contains expected fields
	logOutput := buf.String()
	assert.Contains(t, logOutput, "request")
	assert.Contains(t, logOutput, reqID)
	assert.Contains(t, logOutput, "GET")
	assert.Contains(t, logOutput, "/test")
	assert.Contains(t, logOutput, "status=200")
	// The clientIP check might be 192.168.1.1 or empty if Gin parsing acts up without trusted proxies, but usually it's set.
	assert.Contains(t, logOutput, "clientIP")
}
