package analytics

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"daily-seed/internal/common"
	"daily-seed/internal/daily"
	"daily-seed/internal/habit"
	"daily-seed/internal/task"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AnalyticsHandler struct {
	dailyStore *daily.DailyStore
	taskStore  *task.TaskStore
	habitStore *habit.HabitStore
}

func NewAnalyticsHandler(dailyStore *daily.DailyStore, taskStore *task.TaskStore, habitStore *habit.HabitStore) *AnalyticsHandler {
	return &AnalyticsHandler{
		dailyStore: dailyStore,
		taskStore:  taskStore,
		habitStore: habitStore,
	}
}

func (h *AnalyticsHandler) RegisterRoutes(rg *gin.RouterGroup) {
	analyticsGroup := rg.Group("/analytics")
	{
		analyticsGroup.GET("/heatmap", h.GetHeatmap)
		analyticsGroup.GET("/summary", h.GetSummary)
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

func (h *AnalyticsHandler) GetSummary(c *gin.Context) {
	period := c.DefaultQuery("period", "weekly")
	dateStr := c.Query("date")
	if dateStr == "" {
		dateStr = time.Now().Format("2006-01-02")
	}

	if period != "weekly" && period != "monthly" {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{
			Code:    "INVALID_QUERY",
			Message: "period must be 'weekly' or 'monthly'",
		})
		return
	}

	res, err := h.getSummaryData(c.Request.Context(), period, dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{
			Code:    "INVALID_DATE",
			Message: "failed to get summary data",
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

func calculatePeriodRange(period, dateStr string) (string, string, int, error) {
	var parsedTime time.Time
	var err error

	if len(dateStr) == 7 {
		parsedTime, err = time.Parse("2006-01", dateStr)
	} else {
		parsedTime, err = time.Parse("2006-01-02", dateStr)
	}
	if err != nil {
		return "", "", 0, fmt.Errorf("invalid date format: %w", err)
	}

	if period == "weekly" {
		weekday := int(parsedTime.Weekday())
		start := parsedTime.AddDate(0, 0, -weekday)
		end := start.AddDate(0, 0, 6)
		return start.Format("2006-01-02"), end.Format("2006-01-02"), 7, nil
	}

	year, month, _ := parsedTime.Date()
	start := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, -1)
	totalDays := end.Day()
	return start.Format("2006-01-02"), end.Format("2006-01-02"), totalDays, nil
}

func round1(val float64) float64 {
	return math.Round(val*10) / 10
}

func (h *AnalyticsHandler) getSummaryData(ctx context.Context, period, dateStr string) (*SummaryResponse, error) {
	startDate, endDate, totalDays, err := calculatePeriodRange(period, dateStr)
	if err != nil {
		return nil, err
	}

	records, err := h.dailyStore.FindBetweenDates(ctx, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch daily records: %w", err)
	}

	tasks, err := h.taskStore.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tasks: %w", err)
	}

	habits, err := h.habitStore.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch habits: %w", err)
	}

	taskMap := make(map[primitive.ObjectID]task.Task)
	for _, t := range tasks {
		taskMap[t.ID] = t
	}

	habitMap := make(map[primitive.ObjectID]habit.Habit)
	for _, hb := range habits {
		habitMap[hb.ID] = hb
	}

	recordedDays := len(records)
	modeDistribution := make(map[string]int)
	journals := make([]JournalEntry, 0)

	type taskStatAgg struct {
		Completed int
		Target    int
	}
	taskAggMap := make(map[primitive.ObjectID]*taskStatAgg)

	type habitStatAgg struct {
		Completed int
		Total     int
	}
	habitAggMap := make(map[primitive.ObjectID]*habitStatAgg)

	sectionCompleted := make(map[string]int)
	sectionTarget := make(map[string]int)

	totalTaskCompleted := 0
	totalTaskTarget := 0

	totalHabitCompleted := 0
	totalHabitsTracked := 0

	categoryCompleted := make(map[string]int)
	categoryTotal := make(map[string]int)

	for _, rec := range records {
		if rec.Context.Mode != "" {
			modeDistribution[string(rec.Context.Mode)]++
		}

		if rec.Journal.OneLineReview != "" || rec.Journal.ThreeLineDiary != "" {
			journals = append(journals, JournalEntry{
				Date:           rec.Date,
				OneLineReview:  rec.Journal.OneLineReview,
				ThreeLineDiary: rec.Journal.ThreeLineDiary,
			})
		}

		for _, tEntry := range rec.Tasks {
			tMaster, exists := taskMap[tEntry.TaskID]
			section := "other"
			if exists {
				section = tMaster.Section
			}

			if _, ok := taskAggMap[tEntry.TaskID]; !ok {
				taskAggMap[tEntry.TaskID] = &taskStatAgg{}
			}

			completedAmt := tEntry.ActualAmount
			targetAmt := tEntry.TargetAmount

			taskAggMap[tEntry.TaskID].Completed += completedAmt
			taskAggMap[tEntry.TaskID].Target += targetAmt

			sectionCompleted[section] += completedAmt
			sectionTarget[section] += targetAmt

			totalTaskCompleted += completedAmt
			totalTaskTarget += targetAmt
		}

		for _, hEntry := range rec.Habits {
			hMaster, exists := habitMap[hEntry.HabitID]
			category := "General"
			if exists && hMaster.Category != "" {
				category = hMaster.Category
			}

			if _, ok := habitAggMap[hEntry.HabitID]; !ok {
				habitAggMap[hEntry.HabitID] = &habitStatAgg{}
			}
			habitAggMap[hEntry.HabitID].Total++
			totalHabitsTracked++
			categoryTotal[category]++

			if hEntry.IsCompleted {
				habitAggMap[hEntry.HabitID].Completed++
				totalHabitCompleted++
				categoryCompleted[category]++
			}
		}
	}

	// Task Stats
	overallTaskRate := 0.0
	if totalTaskTarget > 0 {
		overallTaskRate = round1(float64(totalTaskCompleted) / float64(totalTaskTarget) * 100)
	}

	sectionRates := make(map[string]float64)
	for sec, target := range sectionTarget {
		if target > 0 {
			sectionRates[sec] = round1(float64(sectionCompleted[sec]) / float64(target) * 100)
		} else {
			sectionRates[sec] = 0.0
		}
	}

	perTaskStats := make([]TaskStat, 0)
	for tID, agg := range taskAggMap {
		tMaster, exists := taskMap[tID]
		title := "Unknown Task"
		section := "other"
		tType := "quantitative"
		if exists {
			title = tMaster.Title
			section = tMaster.Section
			tType = tMaster.Type
		}

		rate := 0.0
		if agg.Target > 0 {
			rate = round1(float64(agg.Completed) / float64(agg.Target) * 100)
		}

		perTaskStats = append(perTaskStats, TaskStat{
			TaskID:    tID.Hex(),
			Title:     title,
			Section:   section,
			Type:      tType,
			Rate:      rate,
			Completed: agg.Completed,
			Target:    agg.Target,
		})
	}

	// Habit Stats
	overallHabitRate := 0.0
	if totalHabitsTracked > 0 {
		overallHabitRate = round1(float64(totalHabitCompleted) / float64(totalHabitsTracked) * 100)
	}

	categoryRates := make(map[string]float64)
	for cat, total := range categoryTotal {
		if total > 0 {
			categoryRates[cat] = round1(float64(categoryCompleted[cat]) / float64(total) * 100)
		} else {
			categoryRates[cat] = 0.0
		}
	}

	perHabitStats := make([]HabitStat, 0)
	for hID, agg := range habitAggMap {
		hMaster, exists := habitMap[hID]
		title := "Unknown Habit"
		category := "General"
		if exists {
			title = hMaster.Title
			category = hMaster.Category
		}

		rate := 0.0
		if agg.Total > 0 {
			rate = round1(float64(agg.Completed) / float64(agg.Total) * 100)
		}

		perHabitStats = append(perHabitStats, HabitStat{
			HabitID:   hID.Hex(),
			Title:     title,
			Category:  category,
			Rate:      rate,
			Completed: agg.Completed,
			Total:     agg.Total,
		})
	}

	return &SummaryResponse{
		Period:       period,
		StartDate:    startDate,
		EndDate:      endDate,
		TotalDays:    totalDays,
		RecordedDays: recordedDays,
		TaskCompletion: TaskCompletionStats{
			Overall:  overallTaskRate,
			Sections: sectionRates,
			PerTask:  perTaskStats,
		},
		HabitCompletion: HabitCompletionStats{
			Overall:    overallHabitRate,
			Categories: categoryRates,
			PerHabit:   perHabitStats,
		},
		ModeDistribution: modeDistribution,
		Journals:         journals,
	}, nil
}
