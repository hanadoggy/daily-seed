package handler

import (
	"log/slog"
	"net/http"
	"regexp"

	"daily-seed/internal/model"

	"github.com/gin-gonic/gin"
)

var dateRegex = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

type DailyHandler struct {
	svc model.DailyService
}

func NewDailyHandler(svc model.DailyService) *DailyHandler {
	return &DailyHandler{svc: svc}
}

// RegisterRoutes registers daily record routes on the given router group.
func (h *DailyHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/daily/:date", h.GetDailyRecord)
	rg.PATCH("/daily/:date", h.UpdateDailyRecord)
}

// GetDailyRecord handles GET /api/v1/daily/:date.
// Returns the existing record or auto-generates one from active tasks/habits.
func (h *DailyHandler) GetDailyRecord(c *gin.Context) {
	date := c.Param("date")
	if !dateRegex.MatchString(date) {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Code:    "INVALID_DATE",
			Message: "Date must be in YYYY-MM-DD format",
		})
		return
	}

	record, err := h.svc.GetDailyRecord(c.Request.Context(), date)
	if err != nil {
		slog.Error("failed to get daily record", slog.String("date", date), slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
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
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Code:    "INVALID_DATE",
			Message: "Date must be in YYYY-MM-DD format",
		})
		return
	}

	var patch map[string]interface{}
	if err := c.ShouldBindJSON(&patch); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Code:    "INVALID_BODY",
			Message: "Request body must be valid JSON",
		})
		return
	}

	record, err := h.svc.UpdateDailyRecord(c.Request.Context(), date, patch)
	if err != nil {
		slog.Error("failed to update daily record", slog.String("date", date), slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to update daily record",
		})
		return
	}

	c.JSON(http.StatusOK, record)
}
