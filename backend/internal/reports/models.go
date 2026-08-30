// Package reports provides types and logic for loading DeploySure report files.
package reports

// Finding represents a single analysis finding from any report category.
type Finding struct {
	ID                string `json:"id"`
	Category          string `json:"category"`
	Severity          string `json:"severity"`
	Title             string `json:"title"`
	File              string `json:"file,omitempty"`
	Line              int    `json:"line,omitempty"`
	Evidence          string `json:"evidence,omitempty"`
	Description       string `json:"description,omitempty"`
	RecommendedAction string `json:"recommended_action,omitempty"`
}

// ComparisonResult holds the before and after findings side-by-side.
type ComparisonResult struct {
	Before []Finding `json:"before"`
	After  []Finding `json:"after"`
	// Summary captures counts by severity for each stage.
	Summary ComparisonSummary `json:"summary"`
}

// ComparisonSummary holds finding counts for before/after stages.
type ComparisonSummary struct {
	Before SeverityCounts `json:"before"`
	After  SeverityCounts `json:"after"`
}

// SeverityCounts holds the number of findings at each severity level.
type SeverityCounts struct {
	Blocker int `json:"blocker"`
	High    int `json:"high"`
	Medium  int `json:"medium"`
	Low     int `json:"low"`
	Total   int `json:"total"`
}

// CountSeverities returns a SeverityCounts populated from the given findings.
func CountSeverities(findings []Finding) SeverityCounts {
	var c SeverityCounts
	for _, f := range findings {
		switch f.Severity {
		case "blocker":
			c.Blocker++
		case "high":
			c.High++
		case "medium":
			c.Medium++
		case "low":
			c.Low++
		}
		c.Total++
	}
	return c
}

// GraniteAnalysisRequest is the request body for the Granite analysis endpoint.
type GraniteAnalysisRequest struct {
	// Text is the plain-text content to analyze (e.g. YAML, JSON, Go source).
	Text string `json:"text"`
}

// GraniteAnalysisResponse wraps the model's generated analysis text.
type GraniteAnalysisResponse struct {
	Analysis string `json:"analysis"`
}
