package review

import "time"

// Response contains code review results
type Response struct {
	Findings []Finding     `json:"findings"`
	Summary  string        `json:"summary"`
	Provider string        `json:"provider"`
	Duration time.Duration `json:"duration_ms"` // Will be serialized as milliseconds
	Metadata *Metadata     `json:"metadata,omitempty"`
}

// Finding represents a single code issue
type Finding struct {
	Category    string `json:"category"`               // "bug", "security", "performance", "style", "best-practice"
	Severity    string `json:"severity"`               // "critical", "high", "medium", "low", "info"
	Line        *int   `json:"line,omitempty"`         // Line number (nil for file-level issues)
	FilePath    string `json:"file_path,omitempty"`    // Relative file path (for multi-file diffs)
	Description string `json:"description"`            // Issue explanation
	Suggestion  string `json:"suggestion"`             // Remediation advice
	CodeSnippet string `json:"code_snippet,omitempty"` // Relevant code excerpt
}

// Metadata provides additional context about the review
type Metadata struct {
	SourceType   string `json:"source_type"`
	FileCount    int    `json:"file_count,omitempty"`
	LineCount    int    `json:"line_count,omitempty"`
	LinesAdded   int    `json:"lines_added,omitempty"`
	LinesRemoved int    `json:"lines_removed,omitempty"`
	Model        string `json:"model,omitempty"` // Specific LLM model used
}
