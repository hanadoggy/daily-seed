package analytics

import (
	"net/http"
	"strconv"
	"time"

	"daily-seed/internal/common"

	"github.com/gin-gonic/gin"
)

type AnalyticsHandler struct {
	service *AnalyticsService
}

func NewAnalyticsHandler(service *AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{service: service}
}

func (h *AnalyticsHandler) RegisterRoutes(rg *gin.RouterGroup) {
	analyticsGroup := rg.Group("/analytics")
	{
		analyticsGroup.GET("/heatmap", h.GetHeatmap)
	}
}

func (h *AnalyticsHandler) GetHeatmap(c *gin.Context) {
	yearStr := c.Query("year")
	var year int
	if yearStr == "" {
		year = time.Now().Year()
	} else {
		var err error
		year, err = strconv.Atoi(yearStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, common.ErrorResponse{
				Code:    "INVALID_QUERY",
				Message: "year must be a valid integer",
			})
			return
		}
	}

	res, err := h.service.GetHeatmapData(c.Request.Context(), year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "failed to get heatmap data",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, res)
}
