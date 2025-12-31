// Package metrics provides metrics collection for HTTP requests.
package metrics

import (
	"sync"
	"sync/atomic"
	"time"
)

// Collector collects HTTP request metrics.
type Collector struct {
	requestCount   atomic.Int64
	errorCount     atomic.Int64
	totalLatencyMs atomic.Int64
}

// Stats represents collected metrics statistics.
type Stats struct {
	RequestCount int64
	ErrorCount   int64
	AvgLatencyMs int64
}

// NewCollector creates a new metrics collector.
func NewCollector() *Collector {
	return &Collector{}
}

// RecordRequest increments the request counter.
func (c *Collector) RecordRequest() {
	c.requestCount.Add(1)
}

// RecordError increments the error counter.
func (c *Collector) RecordError() {
	c.errorCount.Add(1)
}

// RecordLatency adds latency to the total.
func (c *Collector) RecordLatency(d time.Duration) {
	c.totalLatencyMs.Add(d.Milliseconds())
}

// GetStats returns the current metrics statistics.
func (c *Collector) GetStats() Stats {
	reqCount := c.requestCount.Load()
	errCount := c.errorCount.Load()
	totalLatency := c.totalLatencyMs.Load()

	var avgLatency int64
	if reqCount > 0 {
		avgLatency = totalLatency / reqCount
	}

	return Stats{
		RequestCount: reqCount,
		ErrorCount:   errCount,
		AvgLatencyMs: avgLatency,
	}
}

// LegacyCollector provides the old metrics collection interface using mutex.
// This is kept for backwards compatibility if needed.
type LegacyCollector struct {
	requestCount   int64
	errorCount     int64
	totalLatencyMs int64
	mu             sync.RWMutex
}

// NewLegacyCollector creates a new legacy metrics collector.
func NewLegacyCollector() *LegacyCollector {
	return &LegacyCollector{}
}

// RecordRequest increments the request counter.
func (c *LegacyCollector) RecordRequest() {
	c.mu.Lock()
	c.requestCount++
	c.mu.Unlock()
}

// RecordError increments the error counter.
func (c *LegacyCollector) RecordError() {
	c.mu.Lock()
	c.errorCount++
	c.mu.Unlock()
}

// RecordLatency adds latency to the total.
func (c *LegacyCollector) RecordLatency(d time.Duration) {
	c.mu.Lock()
	c.totalLatencyMs += d.Milliseconds()
	c.mu.Unlock()
}

// GetStats returns the current metrics statistics.
func (c *LegacyCollector) GetStats() Stats {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var avgLatency int64
	if c.requestCount > 0 {
		avgLatency = c.totalLatencyMs / c.requestCount
	}

	return Stats{
		RequestCount: c.requestCount,
		ErrorCount:   c.errorCount,
		AvgLatencyMs: avgLatency,
	}
}
