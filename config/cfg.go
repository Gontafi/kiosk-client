package config

import "time"

type Config struct {
	ServerURL        string
	RegistrationPath string
	HealthReportPath string
	UpdatePath       string
	PollInterval     time.Duration
	HealthInterval   time.Duration
}

func Load() *Config {
	return &Config{
		ServerURL:        "https://example.com/api",
		RegistrationPath: "/device/register",
		HealthReportPath: "/device/health",
		UpdatePath:       "/device/update",
		PollInterval:     24 * time.Hour,
		HealthInterval:   1 * time.Hour,
	}
}
