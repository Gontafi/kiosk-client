package kiosk

import (
	"kiosk-client/config"
	"kiosk-client/pkg/logger"
	"os/exec"
)

func StartKioskController(cfg *config.Config, uuid string) {
	for {
		url := fetchURL(cfg, uuid)
		cmd := exec.Command("chromium-browser", "--kiosk", url)
		err := cmd.Start()
		if err != nil {
			logger.Error("Failed to start Chromium:", err)
			continue
		}

		err = cmd.Wait()
		if err != nil {
			logger.Warn("Chromium crashed, restarting...")
		}
	}
}

func fetchURL(cfg *config.Config, uuid string) string {
	return "https://default-url.com"
}
