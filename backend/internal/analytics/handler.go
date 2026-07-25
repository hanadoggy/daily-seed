package analytics

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"daily-seed/internal/common"
	"daily-seed/internal/daily"
	"daily-seed/internal/task"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AnalyticsHandler struct {
	dailyStore *daily.DailyStore
	taskStore  *task.TaskStore
}

func NewAnalyticsHandler(dailyStore *daily.DailyStore, taskStore *task.TaskStore) *AnalyticsHandler {
	return &AnalyticsHandler{
		dailyStore: dailyStore,
		taskStore:  taskStore,
	}
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

	res, err := h.getHeatmapData(c.Request.Context(), year)
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

func (h *AnalyticsHandler) getHeatmapData(ctx context.Context, year int) (*HeatmapResponse, error) {
	startDate := fmt.Sprintf("%04d-01-01", year)
	endDate := fmt.Sprintf("%04d-12-31", year)

	records, err := h.dailyStore.FindBetweenDates(ctx, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch daily records: %w", err)
	}

	tasks, err := h.taskStore.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tasks: %w", err)
	}

	taskSectionMap := make(map[primitive.ObjectID]string)
	for _, t := range tasks {
		taskSectionMap[t.ID] = t.Section
	}

	recordMap := make(map[string]*daily.DailyRecord)
	for _, r := range records {
		recordMap[r.Date] = r
	}

	var days []HeatmapDay

	loc := time.UTC
	startT := time.Date(year, 1, 1, 0, 0, 0, 0, loc)
	endT := time.Date(year, 12, 31, 0, 0, 0, 0, loc)

	for d := startT; !d.After(endT); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")

		hDay := HeatmapDay{
			Date:          dateStr,
			Total:         0,
			Habits:        0,
			SectionCounts: make(map[string]int),
		}

		if rec, ok := recordMap[dateStr]; ok {
			for _, hb := range rec.Habits {
				if hb.IsCompleted {
					hDay.Habits++
					hDay.Total++
				}
			}

			for _, t := range rec.Tasks {
				if t.IsCompleted {
					hDay.Total++
					if section, exists := taskSectionMap[t.TaskID]; exists {
						hDay.SectionCounts[section]++
					}
				}
			}
		}

		days = append(days, hDay)
	}

	return &HeatmapResponse{Days: days}, nil
}
