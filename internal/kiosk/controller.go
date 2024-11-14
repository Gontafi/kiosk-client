package kiosk

import (
	"encoding/json"
	"fmt"
	"kiosk-client/config"
	"kiosk-client/pkg/logger"
	"kiosk-client/pkg/models"
	"kiosk-client/pkg/utils"
	"os"
	"os/exec"
	"strings"
	"time"
)

const urlFilePath = "last_url.txt"

func StartKioskController(cfg *config.Config, uuid *string) {
	currentURL := loadURLFromFile()

	for {
		newURL := fetchURL(cfg, uuid)
		
		logger.Info(currentURL)
		logger.Info(newURL)
		
		out, err := exec.Command("pgrep", "-f", "chromium-browser").Output()
		if (err != nil || len(out) == 0) || newURL != currentURL {
			currentURL = newURL
			saveURLToFile(currentURL)

			// Kill any running instance of Chromium
			cmd := exec.Command("pkill", "-f", "chromium")
			_ = cmd.Run() // Ignore error if Chromium wasn't running
			cmd = exec.Command("export DISPLAY=:0")
			_ = cmd.Run() // Change to main display
			// Start Chromium in kiosk mode with the new URL
			cmd = exec.Command("chromium", "--kiosk", currentURL)
			output, err := cmd.CombinedOutput()
			logger.Info(fmt.Sprintf("Chromium output: %s", string(output)))
			if err != nil {
				logger.Error("Failed to start Chromium:", err)
				continue
			}
			logger.Info(fmt.Sprintf("Launched Chromium with URL: %s", currentURL))
		}

		time.Sleep(cfg.PollInterval)
	}
}

func fetchURL(cfg *config.Config, uuid *string) string {
	return "https://example.com"
	url := fmt.Sprintf("%s%s/%s", cfg.ServerURL, cfg.GetLinkPath, *uuid)
	resp, _, err := utils.MakeGETRequest(url)

	if err != nil {
		logger.Error("Retrieving URL failed:", err)
		return loadURLFromFile()
	}

	var urlResponse models.URLResponse
	if err := json.Unmarshal(resp, &urlResponse); err != nil {
		logger.Error("Failed to parse JSON response:", err)
		return loadURLFromFile()
	}

	return urlResponse.URL
}

func saveURLToFile(url string) {
	err := os.WriteFile(urlFilePath, []byte(url), 0644)
	if err != nil {
		logger.Error("Failed to save URL to file:", err)
	}
}

func loadURLFromFile() string {
	data, err := os.ReadFile(urlFilePath)
	if err != nil {
		logger.Warn("No previous URL found, defaulting to initial URL")
		return "https://google.com"
	}
	return strings.TrimSpace(string(data))
}
