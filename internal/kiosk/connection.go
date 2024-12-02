package kiosk

import (
	"kiosk-client/config"
	"kiosk-client/pkg/logger"
	"net/http"
	"os/exec"
	"time"
)

var lastConnectionState = true

func connected() bool {
	_, err := http.Get("http://clients3.google.com/generate_204")
	return err == nil
}

func KillChromium(cfg *config.Config) {
	logger.Info("Restarting Chromium after reconnection...")
	_, err := exec.Command("pkill", "-f", cfg.ChromiumCommand).Output()
	if err != nil {
		logger.Error("Failed to kill Chromium:", err)
	}
}

func MonitorConnection(cfg *config.Config) {
	for {
		time.Sleep(5 * time.Second)
		currentState := connected()

		if currentState && !lastConnectionState {
			logger.Info("Internet connection restored")
			KillChromium(cfg)
		} else if !currentState && lastConnectionState {
			logger.Warn("Internet connection lost")
		}

		lastConnectionState = currentState
	}
}
