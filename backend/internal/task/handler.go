package task

import (
	"daily-seed/internal/common"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskHandler struct {
	svc TaskService
}

func NewTaskHandler(svc TaskService) *TaskHandler {
	return &TaskHandler{svc: svc}
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
	tasks, err := h.svc.List(c.Request.Context())
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

	task, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, common.ErrorResponse{
				Code:    "NOT_FOUND",
				Message: err.Error(),
			})
			return
		}
		slog.Error("failed to get task", slog.String("id", id), slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to get task",
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

	created, err := h.svc.Create(c.Request.Context(), &task)
	if err != nil {
		// Validation errors are returned as-is.
		c.JSON(http.StatusBadRequest, common.ErrorResponse{
			Code:    "VALIDATION_ERROR",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, created)
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
	updated, err := h.svc.Update(c.Request.Context(), &task)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, common.ErrorResponse{
				Code:    "NOT_FOUND",
				Message: err.Error(),
			})
			return
		}
		c.JSON(http.StatusBadRequest, common.ErrorResponse{
			Code:    "VALIDATION_ERROR",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, updated)
}

func (h *TaskHandler) Archive(c *gin.Context) {
	id := c.Param("id")

	if err := h.svc.Archive(c.Request.Context(), id); err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, common.ErrorResponse{
				Code:    "NOT_FOUND",
				Message: err.Error(),
			})
			return
		}
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
	progress, err := h.svc.GetProgressForActiveTasks(c.Request.Context())
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

	result, err := h.svc.MigrateTask(c.Request.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, common.ErrorResponse{
				Code:    "NOT_FOUND",
				Message: err.Error(),
			})
			return
		}
		if strings.Contains(err.Error(), "has not reached") || strings.Contains(err.Error(), "non-active") {
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
