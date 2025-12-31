// Package worker provides HTTP client worker pool functionality for load testing.
package worker

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/yinghanhung/prr-playground/internal/retry"
	"github.com/yinghanhung/prr-playground/internal/trace"
)

// Config holds the worker pool configuration.
type Config struct {
	TargetURL   string
	Total       int
	Concurrency int
	Interval    time.Duration
	Timeout     time.Duration
	MaxRetries  int
}

// Pool manages a pool of HTTP client workers.
type Pool struct {
	config Config
	client *http.Client
}

// NewPool creates a new worker pool with the given configuration.
func NewPool(cfg Config) *Pool {
	return &Pool{
		config: cfg,
		client: &http.Client{Timeout: cfg.Timeout},
	}
}

// Run executes the load test with the configured number of workers and requests.
func (p *Pool) Run() {
	log.Printf("starting client target=%s total=%d concurrency=%d interval=%s",
		p.config.TargetURL, p.config.Total, p.config.Concurrency, p.config.Interval)

	jobs := make(chan int, p.config.Total)

	var wg sync.WaitGroup
	for i := 0; i < p.config.Concurrency; i++ {
		wg.Add(1)
		go p.worker(i, jobs, &wg)
	}

	for i := 0; i < p.config.Total; i++ {
		jobs <- i + 1
	}
	close(jobs)

	wg.Wait()
	log.Println("client finished")
}

func (p *Pool) worker(id int, jobs <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		traceID := trace.New()
		success, latency := p.doRequestWithRetry(id, job, traceID)

		if success {
			log.Printf("[worker %d] request %d ok (trace %s) latency=%s", id, job, traceID, latency)
		}

		time.Sleep(p.config.Interval)
	}
}

func (p *Pool) doRequestWithRetry(workerID, jobID int, traceID string) (bool, time.Duration) {
	var lastErr error
	var lastStatusCode int
	var lastLatency time.Duration

	isRetryable := func(err error) bool {
		if err != nil {
			return true // Network errors are retryable
		}
		// 5xx errors are retryable, 4xx (except 429) are not
		return lastStatusCode >= 500 || lastStatusCode == 429
	}

	attempt := 0
	err := retry.Do(context.Background(), p.config.MaxRetries, func() error {
		req, err := http.NewRequest(http.MethodGet, p.config.TargetURL, nil)
		if err != nil {
			log.Printf("[worker %d] request %d build error (trace %s): %v", workerID, jobID, traceID, err)
			return err
		}
		req.Header.Set(trace.HeaderName, traceID)

		start := time.Now()
		resp, err := p.client.Do(req)
		lastLatency = time.Since(start)

		if err != nil {
			lastErr = err
			lastStatusCode = 0
			if attempt > 0 {
				log.Printf("[worker %d] request %d failed (trace %s) attempt %d/%d: %v",
					workerID, jobID, traceID, attempt+1, p.config.MaxRetries+1, err)
			}
			attempt++
			return err
		}

		lastStatusCode = resp.StatusCode
		_ = resp.Body.Close()

		// Success case
		if lastStatusCode < 400 {
			if attempt > 0 {
				log.Printf("[worker %d] request %d succeeded on retry %d (trace %s) status=%d latency=%s",
					workerID, jobID, attempt, traceID, lastStatusCode, lastLatency)
			}
			return nil
		}

		// Failed with status code >= 400
		if !isRetryable(nil) {
			log.Printf("[worker %d] request %d failed non-retryable (trace %s) status=%d",
				workerID, jobID, traceID, lastStatusCode)
			return nil // Don't retry
		}

		if attempt > 0 {
			log.Printf("[worker %d] request %d failed (trace %s) attempt %d/%d status=%d",
				workerID, jobID, traceID, attempt+1, p.config.MaxRetries+1, lastStatusCode)
		}
		attempt++
		lastErr = http.ErrServerClosed // Dummy error to indicate failure
		return lastErr
	}, isRetryable)

	if err != nil {
		log.Printf("[worker %d] request %d failed after %d retries (trace %s) status=%d: %v",
			workerID, jobID, p.config.MaxRetries, traceID, lastStatusCode, lastErr)
		return false, lastLatency
	}

	return lastStatusCode < 400, lastLatency
}
