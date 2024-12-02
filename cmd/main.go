package main

import (
	"kiosk-client/config"
	"kiosk-client/internal/kiosk"
	"kiosk-client/internal/monitor"
	"kiosk-client/internal/registration"
	"kiosk-client/pkg/logger"
)

func startProgram() {
	defer logger.Close()
	cfg := config.Load()
	uuid := registration.RegisterDevice(cfg)

	go monitor.StartHealthReportSender(cfg, &uuid)
	go kiosk.StartKioskController(cfg, &uuid)
	go kiosk.ChromiumRunner(cfg)

	go kiosk.MonitorConnection(cfg)

	logger.Info("Application started")
	select {} // Block forever
}

func main() {
	startProgram()
}
