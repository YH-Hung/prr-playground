package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-webapi-db/internal/config"
	"go-webapi-db/internal/handler"
	"go-webapi-db/internal/health"
	"go-webapi-db/internal/metrics"
	"go-webapi-db/internal/middleware"
	"go-webapi-db/internal/repository"
	"go-webapi-db/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	cfg := config.Load()

	// Initialize MongoDB connection
	mongoClient, err := connectMongoDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoClient.Disconnect(context.Background())

	db := mongoClient.Database(cfg.MongoDB.Database)

	// Initialize MongoDB metrics collector
	mongoMetricsCollector := metrics.NewMongoDBMetricsCollector(mongoClient, cfg.MongoDB.Database, "go-webapi-db")
	mongoMetricsCollector.Start(10 * time.Second) // Collect metrics every 10 seconds
	defer mongoMetricsCollector.Stop()
	
	// Set connection pool configuration metrics
	metrics.SetConnectionPoolConfig("go-webapi-db", cfg.MongoDB.Database, cfg.MongoDB.MaxPoolSize, cfg.MongoDB.MinPoolSize)

	// Initialize services
	metricsService := service.NewMetricsService()
	userRepo := repository.NewUserRepository(db)
	instrumentedRepo := repository.NewInstrumentedUserRepository(userRepo)
	userService := service.NewUserService(instrumentedRepo, metricsService)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService)
	healthHandler := health.NewHealthHandler(db)

	// Register Go runtime metrics
	prometheus.MustRegister(prometheus.NewGoCollector())
	prometheus.MustRegister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))

	// Setup router
	router := setupRouter(cfg, userHandler, healthHandler)

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Server starting on port %s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

func connectMongoDB(cfg *config.Config) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.MongoDB.ConnectTimeout)
	defer cancel()

	opts := options.Client().
		ApplyURI(cfg.MongoDB.URI).
		SetMaxPoolSize(cfg.MongoDB.MaxPoolSize).
		SetMinPoolSize(cfg.MongoDB.MinPoolSize)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	log.Println("Connected to MongoDB successfully")
	return client, nil
}

func setupRouter(cfg *config.Config, userHandler *handler.UserHandler, healthHandler *health.HealthHandler) *gin.Engine {
	router := gin.Default()

	// Middleware
	router.Use(middleware.RecoveryMiddleware())
	router.Use(middleware.MetricsMiddleware())

	// Health check endpoint
	router.GET("/health", healthHandler.HealthCheck)

	// Metrics endpoint
	router.GET(cfg.Metrics.Path, gin.WrapH(promhttp.Handler()))

	// API routes
	api := router.Group("/api/users")
	{
		api.POST("", userHandler.CreateUser)
		api.GET("/:id", userHandler.GetUserByID)
		api.GET("", userHandler.GetAllUsers)
		api.GET("/email/:email", userHandler.GetUserByEmail)
		api.PUT("/:id", userHandler.UpdateUser)
		api.DELETE("/:id", userHandler.DeleteUser)
		api.GET("/status/:status", userHandler.GetUsersByStatus)
		api.GET("/status/:status/count", userHandler.CountUsersByStatus)
		api.GET("/external/:serviceName", userHandler.CallExternalService)
		api.GET("/test/error", userHandler.TriggerError)
		api.GET("/test/slow", userHandler.TriggerSlowResponse)
	}

	return router
}

