package service

import (
	"context"
	"time"

	"github.com/example/back-end-tcc/pkg/logger"
	"github.com/example/back-end-tcc/pkg/models"
	"github.com/example/back-end-tcc/pkg/observability/metrics"
	"github.com/example/back-end-tcc/pkg/queue"
	runnerrepo "github.com/example/back-end-tcc/services/runner/repository"
)

// Option allows customizing service dependencies.
type Option func(*Service)

// WithLogger attaches a logger for instrumentation.
func WithLogger(l logger.Logger) Option {
	return func(s *Service) {
		s.log = l
	}
}

// WithMetrics attaches a metrics recorder.
func WithMetrics(rec metrics.Recorder) Option {
	return func(s *Service) {
		s.metrics = rec
	}
}

// Service consumes submissions and produces results.
type Service struct {
	repo       *runnerrepo.ResultRepository
	subscriber queue.Subscriber
	publisher  queue.Publisher
	log        logger.Logger
	metrics    metrics.Recorder
}

// New creates service.
func New(repo *runnerrepo.ResultRepository, subscriber queue.Subscriber, publisher queue.Publisher, opts ...Option) *Service {
	s := &Service{repo: repo, subscriber: subscriber, publisher: publisher, log: logger.New()}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// Start registers queue consumers.
func (s *Service) Start() {
	s.subscriber.Subscribe("submission.created", s.handleSubmission)
	if s.log != nil {
		s.log.Println("runner: subscribed to submission.created")
	}
}

func (s *Service) handleSubmission(ctx context.Context, msg queue.Message) error {
	start := time.Now()
	submission, ok := msg.Data.(models.Submission)
	if !ok {
		s.observeRun(start, "ignored")
		return nil
	}
	if s.log != nil {
		s.log.Printf("runner: processing submission %s", submission.ID)
	}
	now := time.Now()
	submission.Status = "completed"
	submission.CompletedAt = &now
	submission.ScoreSummary = &models.ScoreSummary{
		Score:      1.0,
		Metrics:    map[string]float64{"accuracy": 1.0},
		Calculated: now,
	}
	s.repo.Save(submission)
	if err := s.publisher.Publish(ctx, queue.Message{Type: "score.calculated", Data: submission}); err != nil {
		if s.log != nil {
			s.log.Printf("runner: failed to publish score for submission %s: %v", submission.ID, err)
		}
		s.observeRun(start, "error")
		return err
	}
	if s.log != nil {
		s.log.Printf("runner: completed submission %s", submission.ID)
	}
	s.observeRun(start, "ok")
	return nil
}

// Results returns processed submissions.
func (s *Service) Results() []models.Submission {
	if s.metrics != nil {
		s.metrics.AddCounter("runner_results_total", map[string]string{"result": "ok"}, float64(len(s.repo.List())))
	}
	return s.repo.List()
}

func (s *Service) observeRun(start time.Time, result string) {
	if s.metrics == nil {
		return
	}
	labels := map[string]string{"result": result}
	s.metrics.AddCounter("runner_runs_total", labels, 1)
	s.metrics.ObserveHistogram("runner_duration_ms", labels, float64(time.Since(start).Milliseconds()))
}
