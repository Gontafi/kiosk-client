package main

import (
	"flag"
	"fmt"
	"kiosk-client/config"
	"kiosk-client/internal/kiosk"
	"kiosk-client/internal/monitor"
	"kiosk-client/internal/registration"
	"kiosk-client/pkg/logger"
	"log"
	"os"
)

func startProgram() {
	defer logger.Close()
	cfg := config.Load()
	uuid := registration.RegisterDevice(cfg)

	go monitor.StartHealthReportSender(cfg, &uuid)
	go kiosk.StartKioskController(cfg, &uuid)

	logger.Info("Application started")
	select {}
}

func statusCheck() {
	fmt.Println("Checking status of program...")

	pidFile := "/tmp/kiosk-client.pid"
	if _, err := os.Stat(pidFile); err == nil {
		fmt.Println("Program is running.")
	} else {
		fmt.Println("Program is not running.")
	}
}

func main() {
	startCmd := flag.Bool("start", false, "Start the program")
	statusCmd := flag.Bool("status", false, "Check the program status")
	flag.Parse()

	pidFile := "/tmp/kiosk-client.pid"

	switch {
	case *startCmd:
		fmt.Println("Starting program...")
		if err := writePID(pidFile); err != nil {
			log.Fatalf("Could not write PID file: %v", err)
		}
		startProgram()

	case *statusCmd:
		statusCheck()

	default:
		fmt.Println("Usage: kiosk-client -start | -status")
		flag.PrintDefaults()
	}
}

func writePID(pidFile string) error {
	pid := os.Getpid()
	file, err := os.Create(pidFile)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(fmt.Sprintf("%d", pid))
	return err
}
