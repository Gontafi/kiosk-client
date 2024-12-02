package kiosk

import (
	"bufio"
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

func runChromium(user string, cfg *config.Config, url string) bool {
	cmd := exec.Command("sudo", "-u", user, "-E", cfg.ChromiumCommand,
		"--kiosk", "--noerrdialogs",
		"--disable-infobars",
		"--no-first-run",
		"--enable-features=OverlayScrollbar",
		"--start-maximized", url)
	cmd.Env = append(os.Environ(), "DISPLAY=:0")

	_, err := cmd.CombinedOutput()

	//logger.Info(fmt.Sprintf("Chromium output: %s", string(output)))
	if err != nil {
		logger.Error("Failed to start Chromium:", err)
		return false
	}

	logger.Info(fmt.Sprintf("Launched Chromium with URL: %s", url))
	return true
}

func ChromiumRunner(cfg *config.Config) {
	for {
		time.Sleep(time.Millisecond * 100)
		out, err := exec.Command("pgrep", "-f", cfg.ChromiumCommand).Output()
		if err != nil || len(out) == 0 {

			currentURL := loadURLFromFile()

			if currentURL == "" {
				continue
			}
			user, err := getNonRootUser()
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			runChromium(user, cfg, currentURL)
		}
	}
}

func StartKioskController(cfg *config.Config, uuid *string) {
	currentURL := loadURLFromFile()

	for {
		newURL := fetchURL(cfg, uuid)

		out, err := exec.Command("pgrep", "-f", cfg.ChromiumCommand).Output()
		if (err != nil || len(out) == 0) || newURL != currentURL {
			currentURL = newURL
			saveURLToFile(currentURL)

			cmd := exec.Command("pkill", "-f", cfg.ChromiumCommand)
			_ = cmd.Run()

			user, err := getNonRootUser()
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			runChromium(user, cfg, currentURL)
		}

		time.Sleep(cfg.PollInterval)
	}
}

func fetchURL(cfg *config.Config, uuid *string) string {
	url := fmt.Sprintf("%s%s/%s", cfg.ServerURL, cfg.GetLinkPath, *uuid)
	resp, _, err := utils.MakeGETRequest(strings.TrimSuffix(url, "\n"))

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

func getNonRootUser() (string, error) {
	file, err := os.Open("/etc/passwd")
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		if len(parts) > 2 {
			// Check if UID is 1000 (default for main non-root user)
			if parts[2] == "1000" {
				return parts[0], nil // return the username
			}
		}
	}
	return "", fmt.Errorf("non-root user not found")
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

	cleanedData := strings.ReplaceAll(string(data), "\n", "")
	cleanedData = strings.ReplaceAll(cleanedData, "\r", "")

	return strings.TrimSpace(cleanedData)
}
