package main

import (
	"kiosk-client/config"
	"kiosk-client/internal/kiosk"
	"kiosk-client/internal/monitor"
	"kiosk-client/internal/registration"
	"kiosk-client/internal/updater"
)

func main() {
	cfg := config.Load()

	uuid := registration.RegisterDevice(cfg)

	go monitor.StartHealthReportSender(cfg, uuid)
	go kiosk.StartKioskController(cfg, uuid)
	go updater.CheckForUpdates(cfg, uuid)

	select {}
}
