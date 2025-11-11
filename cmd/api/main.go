package main

import (
	"fmt"
	"net/http"

	"github.com/example/back-end-tcc/pkg/config"
	"github.com/example/back-end-tcc/pkg/logger"
	"github.com/example/back-end-tcc/pkg/models"
	"github.com/example/back-end-tcc/pkg/observability/metrics"
	"github.com/example/back-end-tcc/pkg/queue"
	"github.com/example/back-end-tcc/pkg/storage"
	agenthandlers "github.com/example/back-end-tcc/services/agent/handlers"
	agentrepository "github.com/example/back-end-tcc/services/agent/repository"
	agentservice "github.com/example/back-end-tcc/services/agent/service"
	authhandlers "github.com/example/back-end-tcc/services/auth/handlers"
	authrepository "github.com/example/back-end-tcc/services/auth/repository"
	authservice "github.com/example/back-end-tcc/services/auth/service"
	benchmarkhandlers "github.com/example/back-end-tcc/services/benchmark/handlers"
	benchmarkrepository "github.com/example/back-end-tcc/services/benchmark/repository"
	benchmarkservice "github.com/example/back-end-tcc/services/benchmark/service"
	leaderboardhandlers "github.com/example/back-end-tcc/services/leaderboard/handlers"
	leaderboardrepository "github.com/example/back-end-tcc/services/leaderboard/repository"
	leaderboardservice "github.com/example/back-end-tcc/services/leaderboard/service"
	orchestratorhandlers "github.com/example/back-end-tcc/services/orchestrator/handlers"
	orchestratorrepository "github.com/example/back-end-tcc/services/orchestrator/repository"
	orchestratorservice "github.com/example/back-end-tcc/services/orchestrator/service"
	runnerhandlers "github.com/example/back-end-tcc/services/runner/handlers"
	runnerrepository "github.com/example/back-end-tcc/services/runner/repository"
	runnerservice "github.com/example/back-end-tcc/services/runner/service"
	scoringhandlers "github.com/example/back-end-tcc/services/scoring/handlers"
	scoringrepository "github.com/example/back-end-tcc/services/scoring/repository"
	scoringservice "github.com/example/back-end-tcc/services/scoring/service"
	tracehandlers "github.com/example/back-end-tcc/services/trace/handlers"
	tracerepository "github.com/example/back-end-tcc/services/trace/repository"
	traceservice "github.com/example/back-end-tcc/services/trace/service"
)

func main() {
	cfg, err := config.FromEnv()
	if err != nil {
		panic(err)
	}
	apiLog := logger.New(logger.WithPrefix("api "))
	meter := metrics.NewInMemory()

	bus := queue.NewBus(queue.WithLogger(apiLog), queue.WithMetrics(meter))

	newServiceLogger := func(prefix string) logger.Logger {
		return logger.New(logger.WithPrefix(prefix + " "))
	}

	submissionRepo := storage.NewMemoryRepository[models.Submission]()
	resultRepo := storage.NewMemoryRepository[models.Submission]()
	scoreRepo := storage.NewMemoryRepository[models.ScoreSummary]()
	traceRepo := storage.NewMemoryRepository[models.TraceEvent]()
	leaderboardRepo := storage.NewMemoryRepository[models.LeaderboardEntry]()
	benchmarkRepo := storage.NewMemoryRepository[models.Benchmark]()
	agentRepoStore := storage.NewMemoryRepository[models.User]()

	authRepoStore := storage.NewMemoryRepository[models.User]()
	authRepo := authrepository.NewUserRepository(authRepoStore)
	authRepo.Seed(models.User{ID: "admin", Email: "admin@example.com", Role: "admin"})
	authSrv := authservice.NewAuthService(
		authRepo,
		authservice.WithLogger(newServiceLogger("auth")),
		authservice.WithMetrics(meter),
	)
	authHTTP := authhandlers.NewHTTP(authSrv)

	agentRepo := agentrepository.NewAgentRepository(agentRepoStore)
	agentSrv := agentservice.NewAgentService(
		agentRepo,
		agentservice.WithLogger(newServiceLogger("agent")),
		agentservice.WithMetrics(meter),
	)
	agentHTTP := agenthandlers.NewHTTP(agentSrv)

	benchmarkRepoImpl := benchmarkrepository.New(benchmarkRepo)
	benchmarkSrv := benchmarkservice.New(
		benchmarkRepoImpl,
		benchmarkservice.WithLogger(newServiceLogger("benchmark")),
		benchmarkservice.WithMetrics(meter),
	)
	benchmarkHTTP := benchmarkhandlers.New(benchmarkSrv)

	orchestratorRepo := orchestratorrepository.New(submissionRepo)
	orchestratorSrv := orchestratorservice.New(
		orchestratorRepo,
		bus,
		orchestratorservice.WithLogger(newServiceLogger("orchestrator")),
		orchestratorservice.WithMetrics(meter),
	)
	orchestratorHTTP := orchestratorhandlers.New(orchestratorSrv)

	runnerRepo := runnerrepository.New(resultRepo)
	runnerSrv := runnerservice.New(
		runnerRepo,
		bus,
		bus,
		runnerservice.WithLogger(newServiceLogger("runner")),
		runnerservice.WithMetrics(meter),
	)
	runnerSrv.Start()
	runnerHTTP := runnerhandlers.New(runnerSrv)

	scoringRepo := scoringrepository.New(scoreRepo)
	scoringSrv := scoringservice.New(
		scoringRepo,
		bus,
		bus,
		scoringservice.WithLogger(newServiceLogger("scoring")),
		scoringservice.WithMetrics(meter),
	)
	scoringSrv.Start()
	scoringHTTP := scoringhandlers.New(scoringSrv)

	traceRepoImpl := tracerepository.New(traceRepo)
	traceSrv := traceservice.New(
		traceRepoImpl,
		bus,
		traceservice.WithLogger(newServiceLogger("trace")),
		traceservice.WithMetrics(meter),
	)
	traceHTTP := tracehandlers.New(traceSrv)

	leaderboardRepoImpl := leaderboardrepository.New(leaderboardRepo)
	leaderboardSrv := leaderboardservice.New(
		leaderboardRepoImpl,
		bus,
		leaderboardservice.WithLogger(newServiceLogger("leaderboard")),
		leaderboardservice.WithMetrics(meter),
	)
	leaderboardSrv.Start()
	leaderboardHTTP := leaderboardhandlers.New(leaderboardSrv)

	mux := http.NewServeMux()
	mux.HandleFunc("/auth", authHTTP.Authenticate)
	mux.HandleFunc("/agents", withMethod(map[string]http.HandlerFunc{
		http.MethodPost: agentHTTP.Register,
		http.MethodGet:  agentHTTP.List,
	}))
	mux.HandleFunc("/benchmarks", withMethod(map[string]http.HandlerFunc{
		http.MethodPost: benchmarkHTTP.Create,
		http.MethodGet:  benchmarkHTTP.List,
	}))
	mux.HandleFunc("/submissions", withMethod(map[string]http.HandlerFunc{
		http.MethodPost: orchestratorHTTP.Submit,
		http.MethodGet:  orchestratorHTTP.List,
	}))
	mux.HandleFunc("/results", runnerHTTP.Results)
	mux.HandleFunc("/scores", scoringHTTP.List)
	mux.HandleFunc("/traces", withMethod(map[string]http.HandlerFunc{
		http.MethodPost: traceHTTP.Record,
		http.MethodGet:  traceHTTP.List,
	}))
	mux.HandleFunc("/leaderboard", leaderboardHTTP.List)

	apiLog.Printf("API gateway listening on :%d", cfg.HTTPPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.HTTPPort), mux); err != nil {
		apiLog.Println("server error:", err)
	}
}

func withMethod(handlers map[string]http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if h, ok := handlers[r.Method]; ok {
			h(w, r)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
