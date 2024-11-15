package utils

import (
	"fmt"
	"os"

	"github.com/google/uuid"
)

const uuidFilePath = "device_uuid.txt"

func LoadOrCreateUUID() (string, bool) {
	if currentUUID := loadUUID(); currentUUID != "" {
		return currentUUID, true
	}
	newUUID := uuid.New().String()
	saveUUID(newUUID)
	return newUUID, false
}

func loadUUID() string {
	data, err := os.ReadFile(uuidFilePath)
	if err != nil {
		return ""
	}
	return string(data)
}

func saveUUID(uuid string) {
	err := os.WriteFile(uuidFilePath, []byte(uuid), 0644)
	if err != nil {
		fmt.Printf("Failed to save UUID to file: %v\n", err)
	}
}
