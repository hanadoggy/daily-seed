package analytics

type HabitStreak struct {
	HabitID       string `json:"habitId"`
	Title         string `json:"title"`
	Category      string `json:"category"`
	CurrentStreak int    `json:"currentStreak"`
	LongestStreak int    `json:"longestStreak"`
	TotalDays     int    `json:"totalDays"`
	LastCompleted string `json:"lastCompleted"` // YYYY-MM-DD or ""
	Milestones    []int  `json:"milestones"`    // e.g. [7, 30, 100, 365]
}

type StreakResponse struct {
	Habits []HabitStreak `json:"habits"`
}
