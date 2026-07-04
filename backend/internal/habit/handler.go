package habit

import (
	"daily-seed/internal/common"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type HabitHandler struct {
	svc HabitService
}

func NewHabitHandler(svc HabitService) *HabitHandler {
	return &HabitHandler{svc: svc}
}

func (h *HabitHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/habits", h.List)
	rg.GET("/habits/:id", h.Get)
	rg.POST("/habits", h.Create)
	rg.PUT("/habits/:id", h.Update)
	rg.DELETE("/habits/:id", h.Archive)
}

func (h *HabitHandler) List(c *gin.Context) {
	habits, err := h.svc.List(c.Request.Context())
	if err != nil {
		slog.Error("failed to list habits", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to list habits",
		})
		return
	}

	c.JSON(http.StatusOK, habits)
}

func (h *HabitHandler) Get(c *gin.Context) {
	id := c.Param("id")

	habit, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, common.ErrorResponse{
				Code:    "NOT_FOUND",
				Message: err.Error(),
			})
			return
		}
		slog.Error("failed to get habit", slog.String("id", id), slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to get habit",
		})
		return
	}

	c.JSON(http.StatusOK, habit)
}

func (h *HabitHandler) Create(c *gin.Context) {
	var habit Habit
	if err := c.ShouldBindJSON(&habit); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{
			Code:    "INVALID_BODY",
			Message: "Request body must be valid JSON",
		})
		return
	}

	created, err := h.svc.Create(c.Request.Context(), &habit)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{
			Code:    "VALIDATION_ERROR",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, created)
}

func (h *HabitHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var habit Habit
	if err := c.ShouldBindJSON(&habit); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{
			Code:    "INVALID_BODY",
			Message: "Request body must be valid JSON",
		})
		return
	}

	habit.ID = id
	updated, err := h.svc.Update(c.Request.Context(), &habit)
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

func (h *HabitHandler) Archive(c *gin.Context) {
	id := c.Param("id")

	if err := h.svc.Archive(c.Request.Context(), id); err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, common.ErrorResponse{
				Code:    "NOT_FOUND",
				Message: err.Error(),
			})
			return
		}
		slog.Error("failed to archive habit", slog.String("id", id), slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to archive habit",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "archived"})
}
