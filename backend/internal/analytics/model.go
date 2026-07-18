package analytics

type HeatmapDay struct {
	Date          string         `json:"date"`
	Total         int            `json:"total"`
	Habits        int            `json:"habits"`
	SectionCounts map[string]int `json:"sectionCounts"`
}

type HeatmapResponse struct {
	Days []HeatmapDay `json:"days"`
}
