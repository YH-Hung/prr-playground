package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server   ServerConfig
	MongoDB  MongoDBConfig
	Metrics  MetricsConfig
}

type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type MongoDBConfig struct {
	URI            string
	Database       string
	ConnectTimeout time.Duration
	MaxPoolSize    uint64
	MinPoolSize    uint64
}

type MetricsConfig struct {
	Path string
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8080"),
			ReadTimeout:  getDurationEnv("SERVER_READ_TIMEOUT", 15*time.Second),
			WriteTimeout: getDurationEnv("SERVER_WRITE_TIMEOUT", 15*time.Second),
		},
		MongoDB: MongoDBConfig{
			URI:            getEnv("MONGODB_URI", "mongodb://localhost:27017"),
			Database:       getEnv("MONGODB_DATABASE", "go_webapi_db"),
			ConnectTimeout: getDurationEnv("MONGODB_CONNECT_TIMEOUT", 10*time.Second),
			MaxPoolSize:    getUint64Env("MONGODB_MAX_POOL_SIZE", 10),
			MinPoolSize:    getUint64Env("MONGODB_MIN_POOL_SIZE", 5),
		},
		Metrics: MetricsConfig{
			Path: getEnv("METRICS_PATH", "/metrics"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if d, err := time.ParseDuration(value); err == nil {
			return d
		}
	}
	return defaultValue
}

func getUint64Env(key string, defaultValue uint64) uint64 {
	if value := os.Getenv(key); value != "" {
		if v, err := strconv.ParseUint(value, 10, 64); err == nil {
			return v
		}
	}
	return defaultValue
}

