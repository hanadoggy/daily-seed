package habit

import (
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
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.JSON(http.StatusOK, habits)
}

func (h *HabitHandler) Get(c *gin.Context) {
	id := c.Param("id")

	habit, err := h.store.FindByID(c.Request.Context(), id)
	if err != nil {
		slog.Error("failed to get habit", slog.String("id", id), slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	if habit == nil {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	c.JSON(http.StatusOK, habit)
}

func (h *HabitHandler) Create(c *gin.Context) {
	var habit Habit
	if err := c.ShouldBindJSON(&habit); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	if err := habit.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	habit.ID = primitive.NewObjectID()
	habit.Status = "active"

	if err := h.store.Create(c.Request.Context(), &habit); err != nil {
		slog.Error("failed to create habit", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.JSON(http.StatusCreated, habit)
}

func (h *HabitHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var habit Habit
	if err := c.ShouldBindJSON(&habit); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	habit.ID = oid

	existing, err := h.store.FindByID(c.Request.Context(), id)
	if err != nil {
		slog.Error("failed to find habit for update", slog.String("id", id), slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	if existing == nil {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	if err := habit.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	habit.Status = existing.Status

	if err := h.store.Update(c.Request.Context(), &habit); err != nil {
		slog.Error("failed to update habit", slog.String("id", id), slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.JSON(http.StatusOK, habit)
}

func (h *HabitHandler) Archive(c *gin.Context) {
	id := c.Param("id")

	existing, err := h.store.FindByID(c.Request.Context(), id)
	if err != nil {
		slog.Error("failed to find habit for archive", slog.String("id", id), slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	if existing == nil {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	existing.Status = "archived"
	if err := h.store.Update(c.Request.Context(), existing); err != nil {
		slog.Error("failed to archive habit", slog.String("id", id), slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "archived"})
}
