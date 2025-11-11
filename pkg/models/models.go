package models

import "time"

// User represents an authenticated subject.
type User struct {
	ID    string
	Email string
	Role  string
}

// Benchmark describes a benchmark suite definition.
type Benchmark struct {
	ID          string
	Name        string
	Description string
	CreatedAt   time.Time
}

// Submission is a benchmark submission by an agent.
type Submission struct {
	ID           string
	AgentID      string
	BenchmarkID  string
	Payload      string
	SubmittedAt  time.Time
	CompletedAt  *time.Time
	Status       string
	ScoreSummary *ScoreSummary
}

// ScoreSummary captures scoring results.
type ScoreSummary struct {
	Score      float64
	Metrics    map[string]float64
	Calculated time.Time
}

// TraceEvent stores trace logs produced by benchmark runs.
type TraceEvent struct {
	ID           string
	SubmissionID string
	Message      string
	Level        string
	Timestamp    time.Time
}

// LeaderboardEntry is a projection combining benchmark results.
type LeaderboardEntry struct {
	SubmissionID string
	BenchmarkID  string
	AgentID      string
	Score        float64
	Rank         int
}
