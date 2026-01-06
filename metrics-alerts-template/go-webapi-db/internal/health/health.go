package health

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go-webapi-db/internal/metrics"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	healthStatus = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "health_status",
			Help: "Health status of various components (1 = healthy, 0 = unhealthy)",
		},
		[]string{"component"},
	)
)

type HealthHandler struct {
	db *mongo.Database
}

func NewHealthHandler(db *mongo.Database) *HealthHandler {
	return &HealthHandler{
		db: db,
	}
}

func (h *HealthHandler) HealthCheck(c *gin.Context) {
	ctx := c.Request.Context()
	
	health := gin.H{
		"status": "UP",
		"components": gin.H{},
	}

	// Check database connectivity
	dbHealthy := h.checkDatabase(ctx)
	health["components"].(gin.H)["db"] = map[string]interface{}{
		"status": map[bool]string{true: "UP", false: "DOWN"}[dbHealthy],
	}

	if !dbHealthy {
		health["status"] = "DOWN"
		healthStatus.WithLabelValues("db").Set(0)
		c.JSON(http.StatusServiceUnavailable, health)
		return
	}

	healthStatus.WithLabelValues("db").Set(1)
	c.JSON(http.StatusOK, health)
}

func (h *HealthHandler) checkDatabase(ctx context.Context) bool {
	if h.db == nil {
		return false
	}
	
	start := time.Now()
	err := h.db.Client().Ping(ctx, nil)
	duration := time.Since(start)
	
	// Record ping metrics
	metrics.RecordPing("go-webapi-db", h.db.Name(), duration)
	
	if err != nil {
		metrics.RecordConnectionError("go-webapi-db", h.db.Name(), "ping_failed")
		return false
	}
	
	return true
}

