package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/yinghanhung/prr-playground/internal/config"
	"github.com/yinghanhung/prr-playground/internal/logger"
	"github.com/yinghanhung/prr-playground/services/server/internal/handlers"
	"github.com/yinghanhung/prr-playground/services/server/internal/metrics"
	"github.com/yinghanhung/prr-playground/services/server/internal/middleware"
)

const (
	defaultLogPath         = "/var/log/app/app.log"
	defaultPort            = "8080"
	defaultShutdownTimeout = 10 * time.Second
)

func ensureLogFile(path string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	return os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
}

func main() {
	// Load configuration
	logPath := config.GetString("LOG_PATH", defaultLogPath)
	port := config.GetString("PORT", defaultPort)
	shutdownTimeout := config.GetDuration("SHUTDOWN_TIMEOUT", defaultShutdownTimeout)

	// Setup logging
	logFile, err := ensureLogFile(logPath)
	if err != nil {
		log.Fatalf("cannot init log file: %v", err)
	}
	defer func() {
		if err := logFile.Sync(); err != nil {
			log.Printf("failed to sync log file: %v", err)
		}
		if err := logFile.Close(); err != nil {
			log.Printf("failed to close log file: %v", err)
		}
	}()

	// Create loggers: stdout with timestamp, file without timestamp for JSON parsing
	stdoutLogger := logger.New(os.Stdout, "")
	fileLogger := logger.New(logFile, "")

	// Setup metrics collector
	collector := metrics.NewCollector()

	// Setup HTTP routes
	mux := http.NewServeMux()
	mux.Handle("/hello", handlers.Hello(stdoutLogger, fileLogger))
	mux.HandleFunc("/health", handlers.Health())
	mux.HandleFunc("/metrics", handlers.Metrics(collector))

	// Wrap with middleware
	handler := middleware.Trace(stdoutLogger, fileLogger, collector, mux)

	// Create HTTP server
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	// Channel to listen for interrupt signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine
	serverErrChan := make(chan error, 1)
	go func() {
		stdoutLogger.Printf(`{"message":"server starting","addr":":%s"}`, port)
		fileLogger.Printf(`{"message":"server starting","addr":":%s"}\n`, port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErrChan <- err
		}
	}()

	// Wait for interrupt signal or server error
	select {
	case err := <-serverErrChan:
		stdoutLogger.Fatalf(`{"message":"server error","error":"%v"}`, err)
	case sig := <-sigChan:
		stdoutLogger.Printf(`{"message":"received signal","signal":"%v","shutting_down":true}`, sig)
		fileLogger.Printf(`{"message":"received signal","signal":"%v","shutting_down":true}\n`, sig)

		// Create shutdown context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		// Graceful shutdown
		if err := server.Shutdown(ctx); err != nil {
			stdoutLogger.Printf(`{"message":"server shutdown error","error":"%v"}`, err)
			fileLogger.Printf(`{"message":"server shutdown error","error":"%v"}\n`, err)
			server.Close()
		} else {
			stdoutLogger.Println(`{"message":"server shutdown gracefully"}`)
			fileLogger.Printf(`{"message":"server shutdown gracefully"}\n`)
		}

		// Final sync of log file
		if err := logFile.Sync(); err != nil {
			stdoutLogger.Printf(`{"message":"failed to sync log file on shutdown","error":"%v"}`, err)
		}
	}
}
