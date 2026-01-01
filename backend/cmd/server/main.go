package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gitlab.yuki.lan/goodies/gochat/backend/internal/config"
	"gitlab.yuki.lan/goodies/gochat/backend/internal/handler"
	"gitlab.yuki.lan/goodies/gochat/backend/internal/middleware"
	"gitlab.yuki.lan/goodies/gochat/backend/internal/repository"
	"gitlab.yuki.lan/goodies/gochat/backend/internal/service"
)

func main() {
	// Initialize logger
	logger := setupLogger()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to load configuration")
	}

	// Set log level
	setLogLevel(cfg.LogLevel)

	logger.Info().
		Str("port", cfg.Port).
		Str("logLevel", cfg.LogLevel).
		Str("model", cfg.ClaudeModel).
		Msg("starting server")

	// Connect to database
	db, err := connectDB(cfg.DatabaseURL)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to connect to database")
	}
	defer db.Close()

	logger.Info().Msg("connected to database")

	// Initialize repositories
	projectRepo := repository.NewPostgresProjectRepository(db)
	fileRepo := repository.NewPostgresFileRepository(db)
	fileMetadataRepo := repository.NewPostgresFileMetadataRepository(db)
	fileSourceRepo := repository.NewPostgresFileSourceRepository(db)
	discoveryRepo := repository.NewPostgresDiscoveryRepository(db)
	achievementRepo := repository.NewAchievementRepository(db)

	// Initialize Claude service (real or mock)
	var claudeService service.ClaudeMessenger
	var claudeVision service.ClaudeVision
	if os.Getenv("USE_MOCK_CLAUDE") == "true" {
		mockService, err := service.NewMockClaudeService("testdata/discovery")
		if err != nil {
			logger.Fatal().Err(err).Msg("failed to load mock fixtures")
		}
		claudeService = mockService
		claudeVision = mockService // MockClaudeService also implements ClaudeVision
		logger.Info().Msg("using MOCK Claude service (no API calls)")
	} else {
		realClaudeService := service.NewClaudeService(service.ClaudeConfig{
			APIKey:    cfg.ClaudeAPIKey,
			Model:     cfg.ClaudeModel,
			MaxTokens: cfg.ClaudeMaxTokens,
		}, logger)
		claudeService = realClaudeService
		claudeVision = realClaudeService
	}

	// Initialize PRD repository and service
	prdRepo := repository.NewPostgresPRDRepository(db)
	prdService := service.NewPRDService(prdRepo, discoveryRepo, claudeService, logger)

	// Initialize discovery service
	discoveryService := service.NewDiscoveryService(discoveryRepo, projectRepo, logger)
	discoveryService.SetPRDService(prdService)      // Wire PRD generation trigger
	discoveryService.SetMessageCreator(projectRepo) // Wire message creation for welcome messages
	discoveryService.SetClaudeService(claudeService) // Wire Claude for generating welcome messages

	// Initialize agent context service
	agentContextService := service.NewAgentContextService(prdRepo, projectRepo, discoveryRepo, logger)

	// Initialize achievement services (Phase 3: Learning Journey)
	achievementSvc := service.NewAchievementService(achievementRepo, logger)
	nudgeSvc := service.NewNudgeService(achievementRepo, achievementSvc, logger)

	// Initialize chat service
	chatService := service.NewChatService(service.ChatConfig{
		ContextMessageLimit: cfg.ContextMessageLimit,
	}, claudeService, discoveryService, agentContextService, projectRepo, fileRepo, fileMetadataRepo, logger)

	// Initialize completeness checker
	completenessChecker := service.NewCompletenessChecker(fileRepo, logger)

	// Initialize handlers
	healthHandler := handler.NewHealthHandler(db)
	projectHandler := handler.NewProjectHandler(projectRepo)
	fileHandler := handler.NewFileHandler(fileRepo, projectRepo, fileMetadataRepo)
	uploadHandler := handler.NewUploadHandler(projectRepo, fileRepo, fileMetadataRepo, fileSourceRepo, claudeVision, logger)
	discoveryHandler := handler.NewDiscoveryHandler(discoveryService, logger)
	prdHandler := handler.NewPRDHandler(prdService, logger)
	achievementHandler := handler.NewAchievementHandler(achievementSvc, nudgeSvc, logger)
	completenessHandler := handler.NewCompletenessHandler(completenessChecker, logger)
	wsHandler := handler.NewWebSocketHandler(chatService, logger)

	// Set up Gin
	if cfg.LogLevel != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Apply middleware
	router.Use(middleware.Recovery(logger))
	router.Use(middleware.Logging(logger))
	router.Use(middleware.CORS(middleware.CORSConfig{
		AllowOrigins: cfg.CORSOrigins,
	}))

	// Health endpoint
	router.GET("/health", healthHandler.Health)

	// API routes
	api := router.Group("/api")
	{
		projects := api.Group("/projects")
		{
			projects.GET("", projectHandler.List)
			projects.POST("", projectHandler.Create)
			projects.GET("/:id", projectHandler.Get)
			projects.PATCH("/:id", projectHandler.Update)
			projects.DELETE("/:id", projectHandler.Delete)
			projects.GET("/:id/files", fileHandler.ListFiles)
			projects.GET("/:id/download", fileHandler.DownloadProjectZip)
			projects.POST("/:id/upload", uploadHandler.Upload)

			// Discovery routes
			projects.GET("/:id/discovery", discoveryHandler.GetDiscovery)
			projects.PUT("/:id/discovery/stage", discoveryHandler.AdvanceStage)
			projects.PUT("/:id/discovery/data", discoveryHandler.UpdateData)
			projects.POST("/:id/discovery/users", discoveryHandler.AddUser)
			projects.POST("/:id/discovery/features", discoveryHandler.AddFeature)
			projects.POST("/:id/discovery/confirm", discoveryHandler.ConfirmDiscovery)
			projects.POST("/:id/discovery/skip", discoveryHandler.SkipDiscovery)
			projects.DELETE("/:id/discovery", discoveryHandler.ResetDiscovery)

			// PRD routes (project-scoped)
			projects.GET("/:id/prds", prdHandler.ListPRDs)
			projects.GET("/:id/active-prd", prdHandler.GetActivePRD)
			projects.PUT("/:id/active-prd", prdHandler.SetActivePRD)
			projects.DELETE("/:id/active-prd", prdHandler.ClearActivePRD)

			// Completeness check route
			projects.GET("/:id/completeness", completenessHandler.GetCompleteness)
		}
		files := api.Group("/files")
		{
			files.GET("/:id", fileHandler.GetFile)
			files.GET("/:id/download", fileHandler.DownloadFile)
		}

		// PRD routes (direct PRD access)
		prds := api.Group("/prds")
		{
			prds.GET("/:id", prdHandler.GetPRD)
			prds.PUT("/:id/status", prdHandler.UpdatePRDStatus)
			prds.POST("/:id/retry", prdHandler.RetryPRDGeneration)
		}

		// Achievement routes (Phase 3: Learning Journey)
		achievementHandler.RegisterRoutes(api)
	}

	// WebSocket endpoint
	router.GET("/ws/chat", wsHandler.HandleConnection)

	// Create HTTP server
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		logger.Info().Str("addr", server.Addr).Msg("server listening")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("server error")
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info().Msg("shutting down server")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal().Err(err).Msg("server forced to shutdown")
	}

	logger.Info().Msg("server exited")
}

func setupLogger() zerolog.Logger {
	zerolog.TimeFieldFormat = time.RFC3339
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	log.Logger = logger
	return logger
}

func setLogLevel(level string) {
	switch level {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}

func connectDB(databaseURL string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}
