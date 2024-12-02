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
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get("http://clients3.google.com/generate_204")
	if err != nil {
		return false
	}
	resp.Body.Close()
	return true
}

func ReloadChromiumPage(cfg *config.Config) {
	logger.Info("Reloading Chromium page...")
	_, err := exec.Command("xdotool","search", "--onlyvisible", "--class", "chromium", "windowfocus", "key", "F5").Output()
	if err != nil {
		logger.Error("Failed to reload Chromium page:", err)
	}
}

func MonitorConnection(cfg *config.Config) {
	for {
		time.Sleep(5 * time.Second)
		currentState := connected()

		if currentState && !lastConnectionState {
			logger.Info("Internet connection restored")
			ReloadChromiumPage(cfg)
		} else if !currentState && lastConnectionState {
			logger.Warn("Internet connection lost")
		}

		lastConnectionState = currentState
	}
}
