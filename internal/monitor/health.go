package monitor

import (
	"bufio"
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

const logFilePath = "application.log"

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
	url := fmt.Sprintf("%s%s", cfg.ServerURL, cfg.HealthReportPath)
	_, _, err := utils.MakePUTRequest(url, report)
	if err != nil {
		logger.Error("Failed to send health report:", err)
		return
	}
}

func getBrowserStatus() string {
	out, err := exec.Command("pgrep", "-f", "chromium").Output()
	if err != nil || len(out) == 0 {
		logger.Warn("Chromium browser is not running")
		return "not running"
	}
	return "working"
}

func getCPUTemperature() float64 {
	out, err := exec.Command("cat", "/sys/class/thermal/thermal_zone0/temp").Output()
	if err != nil {
		logger.Error("Failed to read CPU temperature:", err)
		return 0.0
	}

	tempMilli, err := strconv.Atoi(strings.TrimSpace(string(out)))
	if err != nil {
		logger.Error("Failed to parse CPU temperature:", err)
		return 0.0
	}
	return float64(tempMilli) / 1000.0
}

func getCPULoad() float64 {
	out, err := exec.Command("sh", "-c", "grep 'cpu ' /proc/stat | awk '{usage=($2+$4)*100/($2+$4+$5)} END {print usage}'").Output()
	if err != nil {
		logger.Error("Failed to read CPU load:", err)
		return 0.0
	}
	
	usageStr := strings.TrimSpace(string(out))
	usage, err := strconv.ParseFloat(usageStr, 64)
	if err != nil {
		logger.Error("Failed to parse CPU usage:", err)
		return 0.0
	}

	return usage
}

func getMemoryUsage() float64 {
	out, err := exec.Command("free", "-m").Output()
	if err != nil {
		logger.Error("Failed to read memory usage:", err)
		return 0.0
	}

	lines := strings.Split(string(out), "\n")
	if len(lines) < 2 {
		logger.Error("Unexpected output from free command")
		return 0.0
	}

	fields := strings.Fields(lines[1]) // Use the second line for memory usage
	totalMem, err := strconv.Atoi(fields[1])
	if err != nil {
		logger.Error("Failed to parse total memory:", err)
		return 0.0
	}
	usedMem, err := strconv.Atoi(fields[2])
	if err != nil {
		logger.Error("Failed to parse used memory:", err)
		return 0.0
	}

	return (float64(usedMem) / float64(totalMem)) * 100.0
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
