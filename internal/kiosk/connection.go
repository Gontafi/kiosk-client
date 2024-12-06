package kiosk

import (
	"kiosk-client/config"
	"kiosk-client/pkg/logger"
	"net/http"
	"os/exec"
	"time"
)

var lastConnectionState = true

func connected(cfg *config.Config) (bool, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get("http://clients3.google.com/generate_204")
	if err != nil {
		return false, err
	}

	resp.Body.Close()

	respServer, err := client.Get(cfg.HealthCheckURL)
	if err != nil {
		return false, err
	}

	respServer.Body.Close()

	if respServer.StatusCode < 200 || respServer.StatusCode >= 300 {
		logger.Warn("Ping failed Status", resp.Status)
		return false, nil
	}

	return true, nil
}

func ReloadChromiumPage(cfg *config.Config) {
	logger.Info("Reloading Chromium page...")
	out, err := exec.Command(
		"env",
		"DISPLAY=:0",
		"xdotool",
		"search",
		"--onlyvisible",
		"--class",
		"chromium",
		"windowfocus",
		"key",
		"F5",
	).Output()
	if err != nil {
		logger.Error("Failed to reload Chromium page:", err, out)
	}
}

func MonitorConnection(cfg *config.Config) {
	for {
		time.Sleep(cfg.PingInterval)
		currentState, _ := connected(cfg)

		if currentState && !lastConnectionState {
			logger.Info("Internet connection restored")
			ReloadChromiumPage(cfg)
		} else if !currentState && lastConnectionState {
			logger.Warn("Internet connection lost")
		}

		lastConnectionState = currentState
	}
}
