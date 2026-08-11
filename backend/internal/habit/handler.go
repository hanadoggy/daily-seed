package habit

import (
	"daily-seed/internal/common"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HabitHandler struct {
	store *HabitStore
}

func NewHabitHandler(store *HabitStore) *HabitHandler {
	return &HabitHandler{store: store}
}

func (h *HabitHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/habits", h.List)
	rg.GET("/habits/:id", h.Get)
	rg.POST("/habits", h.Create)
	rg.PUT("/habits/:id", h.Update)
	rg.DELETE("/habits/:id", h.Archive)
}

func (h *HabitHandler) List(c *gin.Context) {
	habits, err := h.store.FindAll(c.Request.Context())
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

	habit, err := h.store.FindByID(c.Request.Context(), id)
	if err != nil {
		slog.Error("failed to get habit", slog.String("id", id), slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to get habit",
		})
		return
	}
	if habit == nil {
		c.JSON(http.StatusNotFound, common.ErrorResponse{
			Code:    "NOT_FOUND",
			Message: fmt.Sprintf("habit not found: %s", id),
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

	if err := habit.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{
			Code:    "VALIDATION_ERROR",
			Message: err.Error(),
		})
		return
	}

	habit.ID = primitive.NewObjectID()
	habit.Status = "active"

	if err := h.store.Create(c.Request.Context(), &habit); err != nil {
		slog.Error("failed to create habit", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to create habit",
		})
		return
	}

	c.JSON(http.StatusCreated, habit)
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

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{
			Code:    "INVALID_ID",
			Message: "Invalid habit ID format",
		})
		return
	}
	habit.ID = oid

	existing, err := h.store.FindByID(c.Request.Context(), id)
	if err != nil {
		slog.Error("failed to find habit for update", slog.String("id", id), slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to update habit",
		})
		return
	}
	if existing == nil {
		c.JSON(http.StatusNotFound, common.ErrorResponse{
			Code:    "NOT_FOUND",
			Message: fmt.Sprintf("habit not found: %s", id),
		})
		return
	}

	if err := habit.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{
			Code:    "VALIDATION_ERROR",
			Message: err.Error(),
		})
		return
	}

	habit.Status = existing.Status

	if err := h.store.Update(c.Request.Context(), &habit); err != nil {
		slog.Error("failed to update habit", slog.String("id", id), slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to update habit",
		})
		return
	}

	c.JSON(http.StatusOK, habit)
}

func (h *HabitHandler) Archive(c *gin.Context) {
	id := c.Param("id")

	existing, err := h.store.FindByID(c.Request.Context(), id)
	if err != nil {
		slog.Error("failed to find habit for archive", slog.String("id", id), slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to archive habit",
		})
		return
	}
	if existing == nil {
		c.JSON(http.StatusNotFound, common.ErrorResponse{
			Code:    "NOT_FOUND",
			Message: fmt.Sprintf("habit not found: %s", id),
		})
		return
	}

	existing.Status = "archived"
	if err := h.store.Update(c.Request.Context(), existing); err != nil {
		slog.Error("failed to archive habit", slog.String("id", id), slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to archive habit",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "archived"})
}
