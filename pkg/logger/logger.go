package logger

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

var (
	logFile   *os.File
	logMutex  sync.Mutex
	logPrefix = "application.log"
)

func init() {
	var err error
	logFile, err = os.OpenFile(logPrefix, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	log.SetOutput(logFile)
	log.SetFlags(0)
}

func logMessage(level string, args ...interface{}) {
	currentTime := time.Now().In(time.Local)
	timestamp := currentTime.Format("2006-01-02 15:04:05 -0700")
	message := fmt.Sprintf("[%s] [%s] %v\n", timestamp, level, fmt.Sprint(args...))

	logMutex.Lock()
	defer logMutex.Unlock()

	if logFile != nil {
		_, err := logFile.WriteString(message)
		if err != nil {
			fmt.Printf("ERROR: Failed to write to log file: %v\n", err)
		}
	}

	fmt.Print(message)
}

func Info(args ...interface{}) {
	logMessage("INFO", args...)
}

func Error(args ...interface{}) {
	logMessage("ERROR", args...)
}

func Warn(args ...interface{}) {
	logMessage("WARN", args...)
}

func Close() {
	logMutex.Lock()
	defer logMutex.Unlock()

	if logFile != nil {
		err := logFile.Close()
		if err != nil {
			fmt.Printf("ERROR: Failed to close log file: %v\n", err)
		}
		logFile = nil
	}
}
