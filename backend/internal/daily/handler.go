package daily

import (
	"daily-seed/internal/common"
	"fmt"
	"log/slog"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

var dateRegex = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

type DailyHandler struct {
	svc DailyService
}

func NewDailyHandler(svc DailyService) *DailyHandler {
	return &DailyHandler{svc: svc}
}

// RegisterRoutes registers daily record routes on the given router group.
func (h *DailyHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/daily/exists", h.GetExistingRecordDates)
	rg.GET("/daily/:date", h.GetDailyRecord)
	rg.PATCH("/daily/:date", h.UpdateDailyRecord)
}

// GetExistingRecordDates handles GET /api/v1/daily/exists?year=YYYY&month=MM
func (h *DailyHandler) GetExistingRecordDates(c *gin.Context) {
	yearStr := c.Query("year")
	monthStr := c.Query("month")

	var year, month int
	if _, err := fmt.Sscanf(yearStr, "%d", &year); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Code: "INVALID_YEAR", Message: "Invalid year parameter"})
		return
	}
	if _, err := fmt.Sscanf(monthStr, "%d", &month); err != nil || month < 1 || month > 12 {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Code: "INVALID_MONTH", Message: "Invalid month parameter"})
		return
	}

	dates, err := h.svc.GetExistingRecordDates(c.Request.Context(), year, month)
	if err != nil {
		slog.Error("failed to get existing record dates", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{Code: "INTERNAL_ERROR", Message: "Failed to retrieve existing dates"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"dates": dates})
}

// GetDailyRecord handles GET /api/v1/daily/:date.
// Returns the existing record or auto-generates one from active tasks/habits.
func (h *DailyHandler) GetDailyRecord(c *gin.Context) {
	date := c.Param("date")
	if !dateRegex.MatchString(date) {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{
			Code:    "INVALID_DATE",
			Message: "Date must be in YYYY-MM-DD format",
		})
		return
	}

	record, err := h.svc.GetDailyRecord(c.Request.Context(), date)
	if err != nil {
		slog.Error("failed to get daily record", slog.String("date", date), slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to retrieve daily record",
		})
		return
	}

	c.JSON(http.StatusOK, record)
}

// UpdateDailyRecord handles PATCH /api/v1/daily/:date.
// Accepts a partial JSON body and applies updates to the daily record.
func (h *DailyHandler) UpdateDailyRecord(c *gin.Context) {
	date := c.Param("date")
	if !dateRegex.MatchString(date) {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{
			Code:    "INVALID_DATE",
			Message: "Date must be in YYYY-MM-DD format",
		})
		return
	}

	var req UpdateDailyRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{
			Code:    "INVALID_BODY",
			Message: "Request body must be valid JSON",
		})
		return
	}

	record, err := h.svc.UpdateDailyRecord(c.Request.Context(), date, &req)
	if err != nil {
		slog.Error("failed to update daily record", slog.String("date", date), slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to update daily record",
		})
		return
	}

	c.JSON(http.StatusOK, record)
}
