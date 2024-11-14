package main

import (
	"flag"
	"fmt"
	"os"
	"log"
	"kiosk-client/config"
	"kiosk-client/internal/kiosk"
	"kiosk-client/internal/monitor"
	"kiosk-client/internal/registration"
	"kiosk-client/pkg/logger"
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

	lockFile := "./kiosk-client.lock"
	if _, err := os.Stat(lockFile); err == nil {
		fmt.Println("Program is running.")
	} else {
		fmt.Println("Program is not running.")
	}
}

func main() {
	startCmd := flag.Bool("start", false, "Start the program")
	statusCmd := flag.Bool("status", false, "Check the program status")
	flag.Parse()

	lockFile := "./kiosk-client.lock"

	switch {
	case *startCmd:
		fmt.Println("Starting program...")

		if err := createLockFile(lockFile); err != nil {
			log.Fatalf("Could not create lock file: %v", err)
		}
		startProgram()

	case *statusCmd:
		statusCheck()

	default:
		fmt.Println("Usage: kiosk-client -start | -status")
		flag.PrintDefaults()
	}
}

func createLockFile(lockFile string) error {
	if _, err := os.Stat(lockFile); err == nil {
		return fmt.Errorf("another instance is already running")
	}

	file, err := os.Create(lockFile)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString("kiosk-client lock file\n")
	return err
}
