package monitor

import (
	"kiosk-client/config"
	"kiosk-client/pkg/logger"
	"time"
)

func StartHealthReportSender(cfg *config.Config, uuid string) {
	for range time.Tick(cfg.HealthInterval) {
		report := collectHealthData()
		sendHealthReport(cfg, uuid, report)
	}
}

func collectHealthData() map[string]interface{} {
	// Collect system metrics like CPU temp, memory, etc.
	return map[string]interface{}{"cpu_temp": 45.0}
}

func sendHealthReport(cfg *config.Config, uuid string, report map[string]interface{}) {
	// Send report to server
	logger.Info("Health report sent")
}
