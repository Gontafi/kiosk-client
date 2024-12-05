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
	"strconv"
	"strings"
	"sync"
	"time"
)

const urlFilePath = "last_url.txt"

var urlFileMutex sync.Mutex

func runChromium(user string, cfg *config.Config, url string) bool {
	cmd := exec.Command("sudo", "-u", user, "-E", cfg.ChromiumCommand,
		"--kiosk", "--noerrdialogs",
		"--disable-infobars",
		"--no-first-run",
		"--enable-features=OverlayScrollbar",
		"--start-maximized", url)
		
	cmd.Env = append(os.Environ(), "DISPLAY=:0")
	
	err := cmd.Start()

	if err != nil {
		logger.Error("Failed to start Chromium:", err)
		return false
	}

	logger.Info(fmt.Sprintf("Launched Chromium with URL: %s", url))
	return true
}

func ChromiumRunner(cfg *config.Config, uuid *string) {
	for {
		time.Sleep(time.Second * 1)
		out, err := exec.Command("pgrep", "-f", cfg.ChromiumCommand).Output()
		if err != nil || len(out) == 0 {

			currentURL := fetchURL(cfg, uuid)

			if currentURL == "" {
				continue
			}
			user, err := getNonRootUser(cfg)
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
    logger.Info(fmt.Sprintf("Starting Kiosk Controller with initial URL: %s", currentURL))

    for {
        newURL := fetchURL(cfg, uuid)

        out, err := exec.Command("pgrep", "-f", cfg.ChromiumCommand).Output()

        if (err != nil || len(out) == 0) || newURL != currentURL {
            logger.Info(fmt.Sprintf("Updating URL to: %s", newURL))
            currentURL = newURL
            saveURLToFile(currentURL)

            cmd := exec.Command("pkill", "-f", cfg.ChromiumCommand)
            if err := cmd.Run(); err != nil {
                logger.Warn("Failed to kill Chromium process:", err)
            }

            user, err := getNonRootUser(cfg)
            if err != nil {
                logger.Error("Failed to get non-root user:", err)
                return
            }

            if !runChromium(user, cfg, currentURL) {
                logger.Error("Failed to start Chromium.")
            }

        } else {
            logger.Info("No changes in URL or Chromium is already running.")
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

func getNonRootUser(cfg *config.Config) (string, error) {
	if cfg.NonRootUser != "" {
		return cfg.NonRootUser, nil
	}

	file, err := os.Open("/etc/passwd")
	if err != nil {
		return "", fmt.Errorf("unable to open /etc/passwd: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		if len(parts) > 2 {
			uid, err := strconv.Atoi(parts[2])
			if err != nil {
				continue
			}
			if uid >= 1000 && uid < 65534 { // Range for non-system users
				return parts[0], nil
			}
		}
	}
	return "", fmt.Errorf("non-root user not found")
}

func saveURLToFile(url string) {
	urlFileMutex.Lock()
	defer urlFileMutex.Unlock()

	err := os.WriteFile(urlFilePath, []byte(url), 0644)
	if err != nil {
		logger.Error("Failed to save URL to file:", err)
	}
}

func loadURLFromFile() string {
	urlFileMutex.Lock()
	defer urlFileMutex.Unlock()

	data, err := os.ReadFile(urlFilePath)
	if err != nil {
		logger.Warn("No previous URL found or failed to read file. Defaulting to initial URL")
		return "https://google.com"
	}

	cleanedData := strings.TrimSpace(string(data))
	if cleanedData == "" {
		logger.Warn("URL file is empty. Defaulting to initial URL")
		return "https://google.com"
	}

	return cleanedData
}
