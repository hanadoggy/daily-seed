package task

import (
	"context"
	"daily-seed/internal/common"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskHandler struct {
	store      *TaskStore
	aggregator TaskProgressAggregator
	cleaner    DailyCleaner
}

func NewTaskHandler(store *TaskStore, aggregator TaskProgressAggregator, cleaner DailyCleaner) *TaskHandler {
	return &TaskHandler{
		store:      store,
		aggregator: aggregator,
		cleaner:    cleaner,
	}
}

func (h *TaskHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/tasks", h.List)
	rg.GET("/tasks/progress", h.GetProgress)
	rg.GET("/tasks/:id", h.Get)
	rg.POST("/tasks", h.Create)
	rg.PUT("/tasks/:id", h.Update)
	rg.DELETE("/tasks/:id", h.Archive)
	rg.POST("/tasks/:id/migrate", h.Migrate)
}

func (h *TaskHandler) List(c *gin.Context) {
	tasks, err := h.store.FindAll(c.Request.Context())
	if err != nil {
		slog.Error("failed to list tasks", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to list tasks",
		})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func (h *TaskHandler) Get(c *gin.Context) {
	id := c.Param("id")

	task, err := h.store.FindByID(c.Request.Context(), id)
	if err != nil {
		slog.Error("failed to get task", slog.String("id", id), slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to get task",
		})
		return
	}
	if task == nil {
		c.JSON(http.StatusNotFound, common.ErrorResponse{
			Code:    "NOT_FOUND",
			Message: fmt.Sprintf("task not found: %s", id),
		})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) Create(c *gin.Context) {
	var task Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{
			Code:    "INVALID_BODY",
			Message: "Request body must be valid JSON",
		})
		return
	}

	// Default conditions if not provided.
	if len(task.Conditions.Weather) == 0 {
		task.Conditions.Weather = []string{"sunny", "rainy"}
	}
	if len(task.Conditions.Mode) == 0 {
		task.Conditions.Mode = []string{"Growth", "Rest", "Office", "Remote"}
	}
	// Boolean tasks always have dailyTarget = 1.
	if task.Type == "boolean" {
		task.Metrics.DailyTarget = 1
	}

	if err := validateTask(&task); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{
			Code:    "VALIDATION_ERROR",
			Message: err.Error(),
		})
		return
	}

	task.ID = primitive.NewObjectID()
	task.Status = "active"

	if err := h.store.Create(c.Request.Context(), &task); err != nil {
		slog.Error("failed to create task", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to create task",
		})
		return
	}

	c.JSON(http.StatusCreated, task)
}

func (h *TaskHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var task Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{
			Code:    "INVALID_BODY",
			Message: "Request body must be valid JSON",
		})
		return
	}

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{
			Code:    "INVALID_ID",
			Message: "Invalid task ID format",
		})
		return
	}
	task.ID = oid

	existing, err := h.store.FindByID(c.Request.Context(), id)
	if err != nil {
		slog.Error("failed to find task for update", slog.String("id", id), slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to update task",
		})
		return
	}
	if existing == nil {
		c.JSON(http.StatusNotFound, common.ErrorResponse{
			Code:    "NOT_FOUND",
			Message: fmt.Sprintf("task not found: %s", id),
		})
		return
	}
	if existing.Status == "archived" {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{
			Code:    "VALIDATION_ERROR",
			Message: "cannot update an archived task",
		})
		return
	}

	if task.Type != existing.Type {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{
			Code:    "VALIDATION_ERROR",
			Message: "task type cannot be changed after creation",
		})
		return
	}

	if err := validateTask(&task); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{
			Code:    "VALIDATION_ERROR",
			Message: err.Error(),
		})
		return
	}

	task.Status = existing.Status
	task.Type = existing.Type

	// If StartDate is delayed, remove task from past DailyRecords
	if task.StartDate > existing.StartDate && h.cleaner != nil {
		if err := h.cleaner.RemoveTaskFromRecordsBeforeDate(c.Request.Context(), task.ID, task.StartDate); err != nil {
			slog.Error("failed to clean past daily records", slog.String("id", id), slog.String("error", err.Error()))
			c.JSON(http.StatusInternalServerError, common.ErrorResponse{
				Code:    "INTERNAL_ERROR",
				Message: "Failed to clean past daily records",
			})
			return
		}
	}

	if err := h.store.Update(c.Request.Context(), &task); err != nil {
		slog.Error("failed to update task", slog.String("id", id), slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to update task",
		})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) Archive(c *gin.Context) {
	id := c.Param("id")

	existing, err := h.store.FindByID(c.Request.Context(), id)
	if err != nil {
		slog.Error("failed to find task for archive", slog.String("id", id), slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to archive task",
		})
		return
	}
	if existing == nil {
		c.JSON(http.StatusNotFound, common.ErrorResponse{
			Code:    "NOT_FOUND",
			Message: fmt.Sprintf("task not found: %s", id),
		})
		return
	}
	if existing.Status == "archived" {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{
			Code:    "VALIDATION_ERROR",
			Message: "task is already archived",
		})
		return
	}

	existing.Status = "archived"
	existing.EndDate = time.Now().Format("2006-01-02")
	if err := h.store.Update(c.Request.Context(), existing); err != nil {
		slog.Error("failed to archive task", slog.String("id", id), slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to archive task",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "archived"})
}

func (h *TaskHandler) GetProgress(c *gin.Context) {
	progress, err := h.getProgressForActiveTasks(c.Request.Context())
	if err != nil {
		slog.Error("failed to get task progress", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to get task progress",
		})
		return
	}

	c.JSON(http.StatusOK, progress)
}

func (h *TaskHandler) Migrate(c *gin.Context) {
	id := c.Param("id")

	var req MigrateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{
			Code:    "INVALID_BODY",
			Message: "Request body must be valid JSON",
		})
		return
	}

	result, err := h.migrateTask(c.Request.Context(), id, req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, common.ErrorResponse{
				Code:    "NOT_FOUND",
				Message: err.Error(),
			})
			return
		}
		if strings.Contains(err.Error(), "has not reached") || strings.Contains(err.Error(), "non-active") || strings.Contains(err.Error(), "completionDate") {
			c.JSON(http.StatusBadRequest, common.ErrorResponse{
				Code:    "MIGRATION_NOT_ALLOWED",
				Message: err.Error(),
			})
			return
		}
		slog.Error("failed to migrate task", slog.String("id", id), slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to migrate task",
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// Helper methods

var validSections = map[string]bool{
	"japanese": true,
	"dev":      true,
	"self_dev": true,
	"exercise": true,
}

var validTaskTypes = map[string]bool{
	"quantitative": true,
	"boolean":      true,
}

var validWeathers = map[string]bool{"sunny": true, "rainy": true}
var validModes = map[string]bool{"Growth": true, "Rest": true, "Office": true, "Remote": true}

func validateTask(task *Task) error {
	if strings.TrimSpace(task.Title) == "" {
		return fmt.Errorf("title is required")
	}
	if !validSections[task.Section] {
		return fmt.Errorf("section must be one of: japanese, dev, self_dev, exercise")
	}
	if !validTaskTypes[task.Type] {
		return fmt.Errorf("type must be one of: quantitative, boolean")
	}
	if task.Type == "quantitative" && task.Metrics.DailyTarget <= 0 {
		return fmt.Errorf("dailyTarget must be positive for quantitative tasks")
	}
	if task.Metrics.TotalTarget < 0 {
		return fmt.Errorf("totalTarget cannot be negative")
	}
	if task.StartDate == "" {
		return fmt.Errorf("startDate is required")
	}

	if len(task.Conditions.Weather) == 0 {
		return fmt.Errorf("at least one weather condition is required")
	}
	if len(task.Conditions.Mode) == 0 {
		return fmt.Errorf("at least one mode condition is required")
	}
	for _, w := range task.Conditions.Weather {
		if !validWeathers[w] {
			return fmt.Errorf("invalid weather value: %s", w)
		}
	}
	for _, m := range task.Conditions.Mode {
		if !validModes[m] {
			return fmt.Errorf("invalid mode value: %s", m)
		}
	}

	return nil
}

func (h *TaskHandler) getProgressForActiveTasks(ctx context.Context) ([]TaskProgress, error) {
	tasks, err := h.store.FindActiveTasks(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetching active tasks: %w", err)
	}

	taskIDs := make([]primitive.ObjectID, 0, len(tasks))
	for _, t := range tasks {
		taskIDs = append(taskIDs, t.ID)
	}

	progressMap, err := h.aggregator.SumTaskProgressByIDs(ctx, taskIDs)
	if err != nil {
		return nil, fmt.Errorf("fetching cumulative progress: %w", err)
	}

	progress := make([]TaskProgress, 0, len(taskIDs))
	for _, t := range tasks {
		completed := progressMap[t.ID]

		pct := 0.0
		if t.Metrics.TotalTarget > 0 {
			pct = float64(completed) / float64(t.Metrics.TotalTarget) * 100
		}

		progress = append(progress, TaskProgress{
			TaskID:         t.ID,
			Title:          t.Title,
			Type:           t.Type,
			TotalTarget:    t.Metrics.TotalTarget,
			TotalCompleted: completed,
			Percentage:     pct,
		})
	}

	return progress, nil
}

func (h *TaskHandler) migrateTask(ctx context.Context, id string, req MigrateTaskRequest) (*MigrationResult, error) {
	if req.CompletionDate == "" {
		return nil, fmt.Errorf("completionDate is required")
	}
	parsedDate, err := time.Parse("2006-01-02", req.CompletionDate)
	if err != nil {
		return nil, fmt.Errorf("invalid completionDate format")
	}

	existing, err := h.store.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("finding task: %w", err)
	}
	if existing == nil {
		return nil, fmt.Errorf("task not found: %s", id)
	}
	if existing.Status != "active" {
		return nil, fmt.Errorf("cannot migrate non-active task: %s", id)
	}

	if existing.Metrics.TotalTarget > 0 {
		progressMap, err := h.aggregator.SumTaskProgressByIDs(ctx, []primitive.ObjectID{existing.ID})
		if err != nil {
			return nil, fmt.Errorf("checking task progress: %w", err)
		}
		completed := progressMap[existing.ID]
		if completed < existing.Metrics.TotalTarget {
			return nil, fmt.Errorf("task progress (%d/%d) has not reached the target", completed, existing.Metrics.TotalTarget)
		}
	}

	existing.Status = "archived"
	existing.EndDate = req.CompletionDate

	newTask := &Task{
		ID:      primitive.NewObjectID(),
		Section: existing.Section,
		Title:   existing.Title,
		Type:    existing.Type,
		Metrics: TaskMetrics{
			DailyTarget: existing.Metrics.DailyTarget,
			TotalTarget: existing.Metrics.TotalTarget,
		},
		Conditions: existing.Conditions,
		Status:     "active",
		StartDate:  parsedDate.AddDate(0, 0, 1).Format("2006-01-02"),
	}

	if err := h.store.MigrateTaskAtomic(ctx, existing, newTask); err != nil {
		return nil, fmt.Errorf("atomic migration failed: %w", err)
	}

	return &MigrationResult{
		ArchivedTask: *existing,
		NewTask:      *newTask,
	}, nil
}
