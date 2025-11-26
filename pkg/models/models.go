package models

import "time"

// User represents an authenticated subject or an Agent.
type User struct {
	ID        string
	Name      string
	Email     string
	Role      string
	Provider  string            // New: OpenAI, Anthropic, etc.
	Endpoint  string            // New: API Endpoint
	AuthType  string            // New: Bearer, API Key
	Status    string            // New: active, inactive
	CreatedAt time.Time         // New
	Headers   map[string]string // New: Custom headers
}

// Benchmark describes a benchmark suite definition.
type Benchmark struct {
	ID          string
	Name        string
	Description string
	Domain      string    // New: Customer Support, Coding, etc.
	TasksCount  int       // New
	Tasks       []Task    // New
	CreatedAt   time.Time
}

// Task describes a specific task within a benchmark.
type Task struct {
	ID           string
	Prompt       string
	ExpectedTool string
	Constraints  []string
	MaxTurns     int
}

// Submission is a benchmark submission by an agent (Run).
type Submission struct {
	ID            string
	AgentID       string
	AgentName     string // New: Denormalized for easier UI display
	BenchmarkID   string
	BenchmarkName string // New: Denormalized
	Payload       string
	SubmittedAt   time.Time
	CompletedAt   *time.Time
	Status        string
	Progress      int // New: 0-100
	ScoreSummary  *ScoreSummary
}

// ScoreSummary captures scoring results.
type ScoreSummary struct {
	Score           float64
	SuccessRate     float64 // New
	ToolCorrectness float64 // New
	Violations      int     // New
	AvgTurns        float64 // New
	TotalCost       float64 // New
	AvgLatency      float64 // New
	Metrics         map[string]float64
	Calculated      time.Time
}

// TraceEvent stores trace logs produced by benchmark runs.
type TraceEvent struct {
	ID           string
	SubmissionID string
	TaskID       string // New
	TaskName     string // New
	Type         string // New: user, agent, tool
	Message      string // Content
	ToolName     string // New
	Parameters   map[string]string // New
	Result       map[string]string // New
	Level        string
	Timestamp    time.Time
	Success      bool    // New
	Turns        int     // New
	Cost         float64 // New
	Latency      float64 // New
}

// LeaderboardEntry is a projection combining benchmark results.
type LeaderboardEntry struct {
	SubmissionID    string
	BenchmarkID     string
	AgentID         string
	AgentName       string // New
	Score           float64
	SuccessRate     float64 // New
	ToolCorrectness float64 // New
	Violations      int     // New
	AvgTurns        float64 // New
	TotalCost       float64 // New
	AvgLatency      float64 // New
	Rank            int
}
