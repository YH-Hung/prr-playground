package health

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestHealthHandler_HealthCheck(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	// Create a test MongoDB client (will fail connection, but tests the handler)
	client, _ := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	defer client.Disconnect(context.Background())
	
	db := client.Database("test")
	handler := NewHealthHandler(db)

	r := gin.New()
	r.GET("/health", handler.HealthCheck)

	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK && w.Code != http.StatusServiceUnavailable {
		t.Errorf("Expected status 200 or 503, got %d", w.Code)
	}
}

func TestHealthHandler_HealthStatusMetric(t *testing.T) {
	prometheus.DefaultRegisterer = prometheus.NewRegistry()
	
	client, _ := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	defer client.Disconnect(context.Background())
	
	db := client.Database("test")
	handler := NewHealthHandler(db)

	// Check database (will likely fail, but sets metric)
	handler.checkDatabase(context.Background())

	reg := prometheus.DefaultRegisterer.(*prometheus.Registry)
	metrics, err := reg.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	var found bool
	for _, mf := range metrics {
		if mf.GetName() == "health_status" {
			found = true
		}
	}

	if !found {
		t.Error("Metric health_status not found")
	}
}

func TestHealthHandler_HealthStatusLabels(t *testing.T) {
	prometheus.DefaultRegisterer = prometheus.NewRegistry()
	
	client, _ := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	defer client.Disconnect(context.Background())
	
	db := client.Database("test")
	handler := NewHealthHandler(db)

	handler.checkDatabase(context.Background())

	reg := prometheus.DefaultRegisterer.(*prometheus.Registry)
	metrics, err := reg.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	var found bool
	for _, mf := range metrics {
		if mf.GetName() == "health_status" {
			found = true
			for _, metric := range mf.GetMetric() {
				labels := metric.GetLabel()
				var component string
				for _, label := range labels {
					if label.GetName() == "component" {
						component = label.GetValue()
					}
				}
				if component != "db" {
					t.Errorf("Expected component label 'db', got '%s'", component)
				}
			}
		}
	}

	if !found {
		t.Error("Metric health_status with component label not found")
	}
}

