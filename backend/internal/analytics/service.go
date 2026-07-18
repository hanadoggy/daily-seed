package analytics

import (
	"context"
	"fmt"
	"time"

	"daily-seed/internal/daily"
	"daily-seed/internal/task"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DailyRecordRepository interface {
	FindBetweenDates(ctx context.Context, startDate string, endDate string) ([]*daily.DailyRecord, error)
}

type TaskRepository interface {
	FindAll(ctx context.Context) ([]task.Task, error)
}

type AnalyticsService struct {
	dailyRepo DailyRecordRepository
	taskRepo  TaskRepository
}

func NewAnalyticsService(dailyRepo DailyRecordRepository, taskRepo TaskRepository) *AnalyticsService {
	return &AnalyticsService{
		dailyRepo: dailyRepo,
		taskRepo:  taskRepo,
	}
}

func (s *AnalyticsService) GetHeatmapData(ctx context.Context, year int) (*HeatmapResponse, error) {
	startDate := fmt.Sprintf("%04d-01-01", year)
	endDate := fmt.Sprintf("%04d-12-31", year)

	// Fetch all records for the year
	records, err := s.dailyRepo.FindBetweenDates(ctx, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch daily records: %w", err)
	}

	// Fetch all tasks to map ID to Section
	tasks, err := s.taskRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tasks: %w", err)
	}

	taskSectionMap := make(map[primitive.ObjectID]string)
	for _, t := range tasks {
		taskSectionMap[t.ID] = t.Section
	}

	// Initialize days map for O(1) lookup
	recordMap := make(map[string]*daily.DailyRecord)
	for _, r := range records {
		recordMap[r.Date] = r
	}

	// Generate all days in the year
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
			// Count habits
			for _, h := range rec.Habits {
				if h.IsCompleted {
					hDay.Habits++
					hDay.Total++
				}
			}

			// Count tasks
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
