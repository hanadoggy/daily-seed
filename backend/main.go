package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"daily-seed/internal/config"
	"daily-seed/internal/daily"
	"daily-seed/internal/habit"
	"daily-seed/internal/middleware"
	"daily-seed/internal/task"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Structured JSON logging.
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	cfg := config.Load()

	// Connect to MongoDB.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(cfg.MongoURI)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		slog.Error("failed to connect to MongoDB", slog.String("error", err.Error()))
		os.Exit(1)
	}

	if err := client.Ping(ctx, nil); err != nil {
		slog.Error("failed to ping MongoDB", slog.String("error", err.Error()))
		os.Exit(1)
	}
	slog.Info("connected to MongoDB", slog.String("uri", cfg.MongoURI))

	db := client.Database(cfg.MongoDBName)

	// Wire dependencies.
	habitRepo := habit.NewHabitRepository(db)
	taskRepo := task.NewTaskRepository(db)
	dailyRecordRepo := daily.NewDailyRecordRepository(db)

	dailySvc := daily.NewDailyService(dailyRecordRepo, taskRepo, habitRepo)
	taskSvc := task.NewTaskService(taskRepo, dailyRecordRepo)
	habitSvc := habit.NewHabitService(habitRepo)

	dailyHandler := daily.NewDailyHandler(dailySvc)
	taskHandler := task.NewTaskHandler(taskSvc)
	habitHandler := habit.NewHabitHandler(habitSvc)

	// Gin setup.
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.Logger())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "PATCH", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Health check.
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// API routes.
	v1 := r.Group("/api/v1")
	dailyHandler.RegisterRoutes(v1)
	taskHandler.RegisterRoutes(v1)
	habitHandler.RegisterRoutes(v1)

	// Start server.
	addr := fmt.Sprintf(":%s", cfg.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		slog.Info("server starting", slog.String("addr", addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	// Graceful shutdown.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("server forced to shutdown", slog.String("error", err.Error()))
	}

	if err := client.Disconnect(shutdownCtx); err != nil {
		slog.Error("error disconnecting from MongoDB", slog.String("error", err.Error()))
	}

	slog.Info("server exited")
}
