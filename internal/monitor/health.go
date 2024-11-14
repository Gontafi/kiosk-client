package monitor

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
	"time"
)

const logFilePath = "application.log" // Path to the log file

func StartHealthReportSender(cfg *config.Config, uuid *string) {
	for range time.Tick(cfg.HealthInterval) {
		report := CollectHealthData(uuid)
		sendHealthReport(cfg, report)
	}
}

func CollectHealthData(uuid *string) *models.HealthRequest {
	temperature := getCPUTemperature()
	cpuLoad := getCPULoad()
	memoryUsage := getMemoryUsage()
	browserStatus := getBrowserStatus()
	logs := getLastLogLines(logFilePath, 10)

	var logsPtr *string
	if logs != "" {
		logsPtr = &logs
	}

	return &models.HealthRequest{
		DeviceID:      *uuid,
		Temperature:   temperature,
		CPULoad:       cpuLoad,
		MemoryUsage:   memoryUsage,
		BrowserStatus: browserStatus,
		Logs:          logsPtr,
	}
}

func sendHealthReport(cfg *config.Config, report *models.HealthRequest) {
	reportJSON, err := json.MarshalIndent(report, "", "    ")
	if err != nil {
		logger.Error("Failed to marshal health report:", err)
		return
	}

	url := fmt.Sprintf("%s%s", cfg.ServerURL, cfg.HealthReportPath)
	fmt.Println(string(reportJSON), url)
	_, _, err = utils.MakePUTRequest(url, reportJSON)
	if err != nil {
	    logger.Error("Failed to send health report:", err)
	    return
	}
}

func getBrowserStatus() string {
	out, err := exec.Command("pgrep", "-f", "chromium-browser").Output()
	if err != nil || len(out) == 0 {
		logger.Warn("Chromium browser is not running")
		return "not running"
	}
	return "working"
}

func getCPUTemperature() int {
	out, err := exec.Command("cat", "/sys/class/thermal/thermal_zone0/temp").Output()
	if err != nil {
		logger.Error("Failed to read CPU temperature:", err)
		return 0
	}

	tempMilli, err := strconv.Atoi(strings.TrimSpace(string(out)))
	if err != nil {
		logger.Error("Failed to parse CPU temperature:", err)
		return 0
	}
	return tempMilli / 1000
}

func getCPULoad() int {
	out, err := exec.Command("top", "-bn1").Output()
	if err != nil {
		logger.Error("Failed to read CPU load:", err)
		return 0
	}

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.Contains(line, "Cpu(s):") {
			fields := strings.Fields(line)
			idleStr := strings.TrimSuffix(fields[7], "%id")
			idle, err := strconv.ParseFloat(idleStr, 64)
			if err != nil {
				logger.Error("Failed to parse CPU idle percentage:", err)
				return 0
			}
			return int(100.0 - idle) // CPU usage as an integer
		}
	}
	return 0
}

func getMemoryUsage() int {
	out, err := exec.Command("free", "-m").Output()
	if err != nil {
		logger.Error("Failed to read memory usage:", err)
		return 0
	}

	lines := strings.Split(string(out), "\n")
	if len(lines) < 2 {
		logger.Error("Unexpected output from free command")
		return 0
	}

	fields := strings.Fields(lines[1]) // Use the second line for memory usage
	totalMem, err := strconv.Atoi(fields[1])
	if err != nil {
		logger.Error("Failed to parse total memory:", err)
		return 0
	}
	usedMem, err := strconv.Atoi(fields[2])
	if err != nil {
		logger.Error("Failed to parse used memory:", err)
		return 0
	}

	return int((float64(usedMem) / float64(totalMem)) * 100.0)
}

func getLastLogLines(filePath string, n int) string {
	file, err := os.Open(filePath)
	if err != nil {
		logger.Error("Failed to open log file:", err)
		return "Failed to read logs"
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		logger.Error("Failed to read log file:", err)
		return "Failed to read logs"
	}

	if len(lines) > n {
		lines = lines[len(lines)-n:]
	}

	return strings.Join(lines, "\n")
}
