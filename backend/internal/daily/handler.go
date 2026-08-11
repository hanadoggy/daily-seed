package daily

import (
	"context"
	"daily-seed/internal/common"
	"daily-seed/internal/habit"
	"daily-seed/internal/task"
	"daily-seed/pkg/jst"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	dateRegex         = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	ErrRecordNotFound = errors.New("daily record not found")
)

type DailyHandler struct {
	store      *DailyStore
	taskStore  *task.TaskStore
	habitStore *habit.HabitStore
}

func NewDailyHandler(store *DailyStore, taskStore *task.TaskStore, habitStore *habit.HabitStore) *DailyHandler {
	return &DailyHandler{
		store:      store,
		taskStore:  taskStore,
		habitStore: habitStore,
	}
}

func (h *DailyHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/daily/exists", h.GetExistingRecordDates)
	rg.GET("/daily/:date", h.GetDailyRecord)
	rg.PATCH("/daily/:date", h.UpdateDailyRecord)
}

func (h *DailyHandler) GetExistingRecordDates(c *gin.Context) {
	yearStr := c.Query("year")
	monthStr := c.Query("month")

	if yearStr == "" || monthStr == "" {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Code: "INVALID_QUERY", Message: "Missing year or month parameter"})
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Code: "INVALID_YEAR", Message: "Invalid year parameter"})
		return
	}
	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Code: "INVALID_MONTH", Message: "Invalid month parameter"})
		return
	}

	startDate := fmt.Sprintf("%04d-%02d-01", year, month)
	t := time.Date(year, time.Month(month+1), 0, 0, 0, 0, 0, time.UTC)
	endDate := t.Format("2006-01-02")

	records, err := h.store.FindBetweenDates(c.Request.Context(), startDate, endDate)
	if err != nil {
		slog.Error("failed to get existing record dates", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{Code: "INTERNAL_ERROR", Message: "Failed to retrieve existing dates"})
		return
	}

	dates := make([]string, 0, len(records))
	for _, r := range records {
		dates = append(dates, r.Date)
	}

	c.JSON(http.StatusOK, gin.H{"dates": dates})
}

func (h *DailyHandler) GetDailyRecord(c *gin.Context) {
	date := c.Param("date")
	if !dateRegex.MatchString(date) {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{
			Code:    "INVALID_DATE",
			Message: "Date must be in YYYY-MM-DD format",
		})
		return
	}

	record, err := h.getDailyRecord(c.Request.Context(), date)
	if err != nil {
		if errors.Is(err, ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, common.ErrorResponse{
				Code:    "NOT_FOUND",
				Message: fmt.Sprintf("Daily record not found for date: %s", date),
			})
			return
		}
		slog.Error("failed to get daily record", slog.String("date", date), slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to retrieve daily record",
		})
		return
	}

	c.JSON(http.StatusOK, record)
}

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

	record, err := h.updateDailyRecord(c.Request.Context(), date, &req)
	if err != nil {
		if errors.Is(err, ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, common.ErrorResponse{
				Code:    "NOT_FOUND",
				Message: fmt.Sprintf("Daily record not found for date: %s", date),
			})
			return
		}
		slog.Error("failed to update daily record", slog.String("date", date), slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to update daily record",
		})
		return
	}

	c.JSON(http.StatusOK, record)
}

// Private Helper Methods

func validateDate(date string) error {
	if _, err := time.Parse("2006-01-02", date); err != nil {
		return fmt.Errorf("invalid date format, expected YYYY-MM-DD")
	}
	return nil
}

func (h *DailyHandler) getDailyRecord(ctx context.Context, date string) (*DailyRecord, error) {
	if err := validateDate(date); err != nil {
		return nil, err
	}
	record, err := h.store.FindByDate(ctx, date)
	if err != nil {
		return nil, fmt.Errorf("finding daily record: %w", err)
	}

	if record != nil {
		updated, err := h.appendMissingEntries(ctx, record)
		if err != nil {
			slog.Warn("failed to append missing entries to daily record", slog.String("error", err.Error()))
		} else if updated {
			if err := h.store.Upsert(ctx, record); err != nil {
				return nil, fmt.Errorf("persisting updated daily record: %w", err)
			}
		}
		return record, nil
	}

	today := jst.Now().Format("2006-01-02")
	if date != today {
		return nil, fmt.Errorf("%w for date: %s", ErrRecordNotFound, date)
	}

	slog.Info("generating new daily record", slog.String("date", date))
	record, err = h.generateDailyRecord(ctx, date)
	if err != nil {
		return nil, fmt.Errorf("generating daily record: %w", err)
	}

	if err := h.store.Upsert(ctx, record); err != nil {
		return nil, fmt.Errorf("persisting daily record: %w", err)
	}

	return record, nil
}

func (h *DailyHandler) updateDailyRecord(ctx context.Context, date string, req *UpdateDailyRecordRequest) (*DailyRecord, error) {
	if err := validateDate(date); err != nil {
		return nil, err
	}
	existing, err := h.getDailyRecord(ctx, date)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, fmt.Errorf("daily record not found for date: %s", date)
	}

	setFields := buildSetFields(req)
	if len(setFields) == 0 {
		return existing, nil
	}

	if err := h.store.PatchByDate(ctx, date, setFields); err != nil {
		return nil, fmt.Errorf("patching daily record: %w", err)
	}

	return h.store.FindByDate(ctx, date)
}

func buildSetFields(req *UpdateDailyRecordRequest) map[string]interface{} {
	set := make(map[string]interface{})

	if req.Context != nil {
		if req.Context.Mode != nil {
			set["context.mode"] = *req.Context.Mode
		}
		if req.Context.Weather != nil {
			set["context.weather"] = *req.Context.Weather
		}
	}
	if req.Tasks != nil {
		set["tasks"] = req.Tasks
	}
	if req.Habits != nil {
		set["habits"] = req.Habits
	}
	if req.Journal != nil {
		if req.Journal.OneLineReview != nil {
			set["journal.oneLineReview"] = *req.Journal.OneLineReview
		}
		if req.Journal.ThreeLineDiary != nil {
			set["journal.threeLineDiary"] = *req.Journal.ThreeLineDiary
		}
	}

	return set
}

func (h *DailyHandler) generateDailyRecord(ctx context.Context, date string) (*DailyRecord, error) {
	tasks, err := h.taskStore.FindActiveTasks(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetching active tasks: %w", err)
	}

	habits, err := h.habitStore.FindActiveHabits(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetching active habits: %w", err)
	}

	taskEntries := make([]TaskEntry, 0, len(tasks))
	for _, t := range tasks {
		if t.StartDate > date {
			continue
		}
		taskEntries = append(taskEntries, TaskEntry{
			TaskID:       t.ID,
			TargetAmount: t.Metrics.DailyTarget,
			ActualAmount: 0,
			IsCompleted:  false,
		})
	}

	habitEntries := make([]HabitEntry, 0, len(habits))
	for _, hb := range habits {
		habitEntries = append(habitEntries, HabitEntry{
			HabitID:     hb.ID,
			IsCompleted: false,
		})
	}

	return &DailyRecord{
		ID:   primitive.NewObjectID(),
		Date: date,
		Context: DayContext{
			Mode:    "Growth",
			Weather: "sunny",
		},
		Tasks:  taskEntries,
		Habits: habitEntries,
		Journal: Journal{
			OneLineReview:  "",
			ThreeLineDiary: "",
		},
	}, nil
}

func (h *DailyHandler) appendMissingEntries(ctx context.Context, record *DailyRecord) (bool, error) {
	updated := false

	tasks, err := h.taskStore.FindActiveTasks(ctx)
	if err != nil {
		return false, fmt.Errorf("fetching active tasks: %w", err)
	}
	existingTasks := make(map[primitive.ObjectID]bool)
	for _, t := range record.Tasks {
		existingTasks[t.TaskID] = true
	}
	for _, t := range tasks {
		if t.StartDate > record.Date {
			continue
		}
		if !existingTasks[t.ID] {
			record.Tasks = append(record.Tasks, TaskEntry{
				TaskID:       t.ID,
				TargetAmount: t.Metrics.DailyTarget,
				ActualAmount: 0,
				IsCompleted:  false,
			})
			updated = true
		}
	}

	habits, err := h.habitStore.FindActiveHabits(ctx)
	if err != nil {
		return false, fmt.Errorf("fetching active habits: %w", err)
	}
	existingHabits := make(map[primitive.ObjectID]bool)
	for _, hb := range record.Habits {
		existingHabits[hb.HabitID] = true
	}
	for _, hb := range habits {
		if !existingHabits[hb.ID] {
			record.Habits = append(record.Habits, HabitEntry{
				HabitID:     hb.ID,
				IsCompleted: false,
			})
			updated = true
		}
	}

	return updated, nil
}
