package monitor

import (
	"bufio"
	"encoding/json"
	"fmt"
	"kiosk-client/config"
	"kiosk-client/pkg/logger"
	"kiosk-client/pkg/models"
	_ "kiosk-client/pkg/utils"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const logFilePath = "application.log" // Path to the log file

func StartHealthReportSender(cfg *config.Config, uuid *string) {
	for range time.Tick(cfg.HealthInterval) {
		report := collectHealthData(uuid)
		sendHealthReport(cfg, report)
	}
}

func collectHealthData(uuid *string) *models.HealthRequest {
	temperature := getCPUTemperature()
	cpuLoad := getCPULoad()
	memoryUsage := getMemoryUsage()
	browserStatus := getBrowserStatus()
	logs := getLastLogLines(logFilePath, 50) // Read the last 50 lines from the log file

	return &models.HealthRequest{
		DeviceID:      *uuid,
		Temperature:   temperature,
		CPULoad:       cpuLoad,
		MemoryUsage:   memoryUsage,
		BrowserStatus: browserStatus,
		Logs:          logs,
	}
}

func sendHealthReport(cfg *config.Config, report *models.HealthRequest) {
	reportJSON, err := json.MarshalIndent(report, "","    ")
	if err != nil {
		logger.Error("Failed to marshal health report:", err)
		return
	}

	url := fmt.Sprintf("%s%s", cfg.ServerURL, cfg.HealthReportPath)
	
	fmt.Println(string(reportJSON), url)
	//_, _, err = utils.MakePOSTRequest(url, reportJSON)
	//if err != nil {
	//	logger.Error("Failed to send health report:", err)
	//	return
	//}
}

// getBrowserStatus checks if Chromium is running
func getBrowserStatus() string {
	out, err := exec.Command("pgrep", "-f", "chromium-browser").Output()
	if err != nil || len(out) == 0 {
		logger.Warn("Chromium browser is not running")
		return "not running"
	}
	return "working"
}

func getCPUTemperature() string {
	out, err := exec.Command("cat", "/sys/class/thermal/thermal_zone0/temp").Output()
	if err != nil {
		logger.Error("Failed to read CPU temperature:", err)
		return "unknown"
	}

	tempMilli, err := strconv.Atoi(strings.TrimSpace(string(out)))
	if err != nil {
		logger.Error("Failed to parse CPU temperature:", err)
		return "unknown"
	}
	return fmt.Sprintf("%.1f°C", float64(tempMilli)/1000.0)
}

func getCPULoad() string {
	out, err := exec.Command("top", "-bn1").Output()
	if err != nil {
		logger.Error("Failed to read CPU load:", err)
		return "unknown"
	}

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.Contains(line, "Cpu(s):") {
			fields := strings.Fields(line)
			idleStr := strings.TrimSuffix(fields[7], "%id")
			idle, err := strconv.ParseFloat(idleStr, 64)
			if err != nil {
				logger.Error("Failed to parse CPU idle percentage:", err)
				return "unknown"
			}
			return fmt.Sprintf("%.1f%%", 100.0-idle) // CPU usage = 100 - idle
		}
	}
	return "unknown"
}

func getMemoryUsage() string {
	out, err := exec.Command("free", "-m").Output()
	if err != nil {
		logger.Error("Failed to read memory usage:", err)
		return "unknown"
	}

	lines := strings.Split(string(out), "\n")
	if len(lines) < 2 {
		logger.Error("Unexpected output from free command")
		return "unknown"
	}

	fields := strings.Fields(lines[1]) // Use the second line for memory usage
	totalMem, err := strconv.Atoi(fields[1])
	if err != nil {
		logger.Error("Failed to parse total memory:", err)
		return "unknown"
	}
	usedMem, err := strconv.Atoi(fields[2])
	if err != nil {
		logger.Error("Failed to parse used memory:", err)
		return "unknown"
	}

	usage := (float64(usedMem) / float64(totalMem)) * 100.0
	return fmt.Sprintf("%.1f%%", usage)
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
