package analytics

type SummaryResponse struct {
	Period           string                 `json:"period"`           // "weekly" | "monthly"
	StartDate        string                 `json:"startDate"`        // YYYY-MM-DD
	EndDate          string                 `json:"endDate"`          // YYYY-MM-DD
	TotalDays        int                    `json:"totalDays"`        // 총 일수 (주: 7, 월: 28~31)
	RecordedDays     int                    `json:"recordedDays"`     // 기록이 존재하는 일수
	TaskCompletion   TaskCompletionStats    `json:"taskCompletion"`   // Task 완료율 통계
	HabitCompletion  HabitCompletionStats   `json:"habitCompletion"`  // Habit 완료율 통계
	ModeDistribution map[string]int         `json:"modeDistribution"` // 모드별 사용 일수
	Journals         []JournalEntry         `json:"journals"`         // 기간 내 저널 목록
}

type TaskCompletionStats struct {
	Overall  float64            `json:"overall"`  // 전체 평균 완료율 (%)
	Sections map[string]float64 `json:"sections"` // 섹션별 평균 완료율 (%)
	PerTask  []TaskStat         `json:"perTask"`  // 개별 태스크별 완료율
}

type TaskStat struct {
	TaskID    string  `json:"taskId"`
	Title     string  `json:"title"`
	Section   string  `json:"section"`
	Type      string  `json:"type"`
	Rate      float64 `json:"rate"`      // 완료율 (%)
	Completed int     `json:"completed"` // 달성량 합계
	Target    int     `json:"target"`    // 일일 목표량 합계
}

type HabitCompletionStats struct {
	Overall    float64            `json:"overall"`    // 전체 평균 습관 완료율 (%)
	Categories map[string]float64 `json:"categories"` // 카테고리별 평균 습관 완료율 (%)
	PerHabit   []HabitStat        `json:"perHabit"`   // 개별 습관별 완료율
}

type HabitStat struct {
	HabitID   string  `json:"habitId"`
	Title     string  `json:"title"`
	Category  string  `json:"category"`
	Rate      float64 `json:"rate"`      // 완료율 (%)
	Completed int     `json:"completed"` // 완료 일수
	Total     int     `json:"total"`     // 해당 습관이 기록에 존재하는 총 횟수
}

type JournalEntry struct {
	Date           string `json:"date"`
	OneLineReview  string `json:"oneLineReview"`
	ThreeLineDiary string `json:"threeLineDiary"`
}
