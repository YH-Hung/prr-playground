package main

import (
	"flag"
	"time"

	"github.com/yinghanhung/prr-playground/internal/config"
	"github.com/yinghanhung/prr-playground/services/client/internal/worker"
)

func main() {
	// Parse command-line flags
	var cfg worker.Config
	flag.StringVar(&cfg.TargetURL, "target", config.GetString("TARGET_URL", "http://localhost:8080/hello"), "target URL")
	flag.IntVar(&cfg.Total, "count", config.GetInt("CLIENT_COUNT", 20), "total requests to send")
	flag.IntVar(&cfg.Concurrency, "concurrency", config.GetInt("CLIENT_CONCURRENCY", 2), "number of concurrent workers")
	flag.DurationVar(&cfg.Interval, "interval", config.GetDuration("CLIENT_INTERVAL", 500*time.Millisecond), "delay between requests per worker")
	flag.DurationVar(&cfg.Timeout, "timeout", config.GetDuration("CLIENT_TIMEOUT", 3*time.Second), "HTTP client timeout")
	flag.IntVar(&cfg.MaxRetries, "retries", config.GetInt("CLIENT_MAX_RETRIES", 3), "maximum retry attempts for failed requests")
	flag.Parse()

	// Create and run worker pool
	pool := worker.NewPool(cfg)
	pool.Run()
}
