package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	ServerURL        string
	RegistrationPath string
	GetLinkPath      string
	HealthReportPath string
	PollInterval     time.Duration
	HealthInterval   time.Duration
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using default values")
	}

	getEnv := func(key, fallback string) string {
		if value, exists := os.LookupEnv(key); exists {
			return value
		}
		return fallback
	}

	parseHours := func(envValue string, defaultValue time.Duration) time.Duration {
		if hours, err := strconv.Atoi(envValue); err == nil {
			return time.Duration(hours) * time.Hour
		}
		return defaultValue
	}

	//parseMinutes := func(envValue string, defaultValue time.Duration) time.Duration {
	//	if minutes, err := strconv.Atoi(envValue); err == nil {
	//		return time.Duration(minutes) * time.Minute
	//	}
	//	return defaultValue
	//}

	return &Config{
		ServerURL:        getEnv("SERVER_URL", "https://example.com/api"),
		RegistrationPath: getEnv("REGISTRATION_PATH", "/api/register"),
		GetLinkPath:      getEnv("GET_LINK_PATH", "/api/get_link"),
		HealthReportPath: getEnv("HEALTH_REPORT_PATH", "/api/status_update"),
		PollInterval:     parseHours(getEnv("POLL_INTERVAL_HOUR", ""), 24*time.Hour),
		HealthInterval:   5*time.Second,
	}
}
